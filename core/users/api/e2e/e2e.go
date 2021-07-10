package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-core/users/schema"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type config struct {
	Web struct {
		Port                 string        `conf:"default::4000"`
		Debug                string        `conf:"default:localhost:6060"`
		Production           bool          `conf:"default:false"`
		ReadTimeout          time.Duration `conf:"default:5s"`
		WriteTimeout         time.Duration `conf:"default:5s"`
		ShutdownTimeout      time.Duration `conf:"default:5s"`
		CorsOrigins          string        `conf:"default:https://localhost:3000"`
		AuthDomain           string        `conf:"default:none"`
		AuthAudience         string        `conf:"default:none"`
		AuthTestClientID     string        `conf:"default:none"`
		AuthTestClientSecret string        `conf:"default:none"`
		AuthM2MClient        string        `conf:"default:none"`
		AuthM2MClientSecret  string        `conf:"default:none"`
		AuthM2MSecret        string        `conf:"default:none"`
		AuthMAPIAudience     string        `conf:"default:none"`
		SendgridAPIKey       string        `conf:"default:none,noprint"`
	}
	DB struct {
		User       string `conf:"default:postgres"`
		Password   string `conf:"default:postgres"`
		Host       string `conf:"default:localhost"`
		Name       string `conf:"default:postgres"`
		DisableTLS bool   `conf:"default:true"`
	}
	Nats struct {
		URL       string `conf:"default:nats://"`
		ClientID  string `conf:"default:client-id"`
		ClusterID string `conf:"default:cluster-id"`
	}
}

func setupTests(t *testing.T) (config, *database.Repository, func(), *log.Logger) {
	infolog := log.New(os.Stderr, "E2E TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	var cfg config

	if err := conf.Parse(os.Args[1:], "API", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			var help string
			help, err = conf.Usage("API", &cfg)
			if err != nil {
				infolog.Fatal(err, "generating usage")
			}
			fmt.Println(help)
		}
		infolog.Fatal(err, "error: parsing config")
	}

	repo, rClose := newTestRepository(t, cfg)

	return cfg, repo, rClose, infolog
}

func newTestRepository(t *testing.T, cfg config) (*database.Repository, func()) {
	ctx := context.Background()

	postgresPort := nat.Port("5432/tcp")

	postgres, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "postgres",
			ExposedPorts: []string{postgresPort.Port()},
			Env: map[string]string{
				"POSTGRES_PASSWORD": cfg.DB.Password,
				"POSTGRES_USER":     cfg.DB.User,
			},
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort(postgresPort),
			),
		},
		Started: true, // auto-start the container
	})
	if err != nil {
		t.Fatal("start:", err)
	}

	// MappedPort gets the externally mapped port for the container
	hostPort, err := postgres.MappedPort(ctx, postgresPort)
	if err != nil {
		t.Fatal("map:", err)
	}

	repo, rClose, err := database.NewRepository(database.Config{
		User:       cfg.DB.User,
		Host:       cfg.DB.Host + ":" + hostPort.Port(),
		Name:       cfg.DB.Name,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		t.Fatal(err, "connecting to db")
	}

	t.Logf("Postgres container started, running at:  %s\n", repo.URL.String())

	if err := schema.Migrate(repo.URL.String()); err != nil {
		t.Fatal(err)
	}

	if err := schema.Seed(repo.SqlxStorer, "users"); err != nil {
		t.Fatal(err)
	}
	if err := schema.Seed(repo.SqlxStorer, "teams"); err != nil {
		t.Fatal(err)
	}
	if err := schema.Seed(repo.SqlxStorer, "memberships"); err != nil {
		t.Fatal(err)
	}
	if err := schema.Seed(repo.SqlxStorer, "projects"); err != nil {
		t.Fatal(err)
	}
	return repo, rClose
}

// Test stores auth dependencies to support e2e tests
type Test struct {
	t                 *testing.T
	Auth0Domain       string
	Auth0Audience     string
	Auth0ClientID     string
	Auth0ClientSecret string
}

func (t *Test) token(username, password string) string {
	const TokenNotFound = "could not retrieve access token"

	urlStr := fmt.Sprintf("https://%s/oauth/token", t.Auth0Domain)
	jsonStr := fmt.Sprintf(`{ "username":"%s", "password":"%s", "audience":"%s", "client_id":"%s", "client_secret": "%s", "grant_type": "password" }`, username, password, t.Auth0Audience, t.Auth0ClientID, t.Auth0ClientSecret)
	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(jsonStr))
	if err != nil {
		t.t.Fatalf("%s: %v", TokenNotFound, err)
	}

	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.t.Fatalf("%s: %v", TokenNotFound, err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var token struct {
		Token string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &token); err != nil {
		t.t.Fatalf("%s: %v", TokenNotFound, err)
	}
	return token.Token
}
