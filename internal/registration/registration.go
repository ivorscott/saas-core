package registration

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/devpies/saas-core/internal/registration/config"
	"github.com/devpies/saas-core/internal/registration/db"
	"github.com/devpies/saas-core/internal/registration/handler"
	"github.com/devpies/saas-core/internal/registration/repository"
	"github.com/devpies/saas-core/internal/registration/service"
	"github.com/devpies/saas-core/pkg/log"
	"github.com/devpies/saas-core/pkg/msg"

	"github.com/ardanlabs/conf"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

// Run contains the app setup.
func Run() error {
	var (
		cfg      config.Config
		logger   *zap.Logger
		dbClient *dynamodb.Client
		logPath  = "log/out.log"
		err      error
	)

	if err = conf.Parse(os.Args[1:], "REGISTRATION", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var usage string
			usage, err = conf.Usage("REGISTRATION", &cfg)
			if err != nil {
				logger.Error("error generating config usage", zap.Error(err))
				return err
			}
			fmt.Println(usage)
			return nil
		}
		logger.Error("error parsing config", zap.Error(err))
		return err
	}

	ctx := context.Background()

	if cfg.Web.Production {
		logger, err = log.NewProductionLogger(logPath)
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	defer logger.Sync()

	dbClient = db.NewDynamoDBClient(ctx, cfg)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	// Initialize AWS clients.
	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("error loading aws config", zap.Error(err))
		return err
	}
	cognitoClient := cip.NewFromConfig(awsCfg)

	jetStream := msg.NewStreamContext(logger, shutdown, cfg.Nats.Address, cfg.Nats.Port)

	_ = jetStream.Create(msg.StreamTenants)

	// Initialize 3-layered architecture.
	authInfoRepo := repository.NewAuthInfoRepository(logger, dbClient, cfg.Dynamodb.AuthTable)

	idpService := service.NewIDPService(logger, cfg, cognitoClient, authInfoRepo, jetStream)
	registrationService := service.NewRegistrationService(logger, cfg.Cognito.Region, idpService, jetStream)
	registrationHandler := handler.NewRegistrationHandler(logger, registrationService)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Web.Port),
		WriteTimeout: cfg.Web.WriteTimeout,
		ReadTimeout:  cfg.Web.ReadTimeout,
		Handler:      Routes(logger, shutdown, registrationHandler, cfg),
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting registration service on %s:%s", cfg.Web.Address, cfg.Web.Port))
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err = <-serverErrors:
		logger.Error("error on startup", zap.Error(err))
		return err
	case sig := <-shutdown:
		logger.Info(fmt.Sprintf("Start shutdown due to %s signal", sig))

		// Give on going tasks a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err = srv.Shutdown(ctx)
		if err != nil {
			err = srv.Close()
		}

		switch {
		case sig == syscall.SIGSTOP:
			logger.Error("error on integrity issue caused shutdown", zap.Error(err))
			return err
		case err != nil:
			logger.Error("error on gracefully shutdown", zap.Error(err))
			return err
		}
	}

	return err
}
