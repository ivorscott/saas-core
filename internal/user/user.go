package user

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/devpies/saas-core/internal/user/config"
	"github.com/devpies/saas-core/internal/user/service"
	"github.com/devpies/saas-core/pkg/log"
	"github.com/devpies/saas-core/pkg/msg"

	"github.com/ardanlabs/conf"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Run contains the app setup.
func Run() error {
	var (
		cfg     config.Config
		logger  *zap.Logger
		logPath = "log/out.log"
		err     error
	)

	if err = conf.Parse(os.Args[1:], "USER", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var usage string
			usage, err = conf.Usage("USER", &cfg)
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

	if cfg.Web.Production {
		logger, err = log.NewProductionLogger(logPath)
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		logger.Error("error creating logger", zap.Error(err))
		return err
	}
	defer logger.Sync()

	// Initialize AWS clients.
	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("error loading aws config", zap.Error(err))
		return err
	}
	cognitoClient := cip.NewFromConfig(awsCfg)
	userService := service.NewUserService(logger, cognitoClient)

	// Initialize channels for graceful shutdown.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Initialize NATS JetStream.
	js := msg.NewStreamContext(logger, shutdown, cfg.Nats.Address, cfg.Nats.Port)
	opts := []nats.SubOpt{nats.DeliverAll(), nats.ManualAck()}

	logger.Info(fmt.Sprintf("Starting user service on %s:%s", cfg.Web.Address, cfg.Web.Port))

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("listener panic: %v", r)
				logger.Error(fmt.Sprintf("%s", debug.Stack()), zap.Error(err))
			}
		}()

		js.Listen(
			string(msg.TypeTenantRegistered),
			msg.SubjectRegistered,
			"user_consumer",
			userService.CreateTenantUserFromMessage,
			opts...,
		)
	}()

	select {
	case sig := <-shutdown:
		logger.Info(fmt.Sprintf("Start shutdown due to %s signal", sig))

		switch {
		case sig == syscall.SIGSTOP:
			logger.Error("error on integrity issue caused shutdown")
		default:
		}
	}

	return nil
}
