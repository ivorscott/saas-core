package admin

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/devpies/saas-core/internal/admin/config"
	"github.com/devpies/saas-core/internal/admin/handler"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/devpies/saas-core/pkg/web/mid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// Routes composes routes, middleware and handlers.
func Routes(
	log *zap.Logger,
	shutdown chan os.Signal,
	assets fs.FS,
	config config.Config,
	authHandler *handler.AuthHandler,
	webPageHandler *handler.WebPageHandler,
	registrationHandler *handler.RegistrationHandler,
	tenantsHandler *handler.TenantHandler,
) http.Handler {
	mux := chi.NewRouter()
	mux.Use(loadSession)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost", "https://devpie.io"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(assets))))

	middleware := []web.Middleware{
		mid.Logger(log),
		mid.Errors(log),
		withAuth(log, config.Cognito.Region, config.Cognito.UserPoolID),
		mid.Panics(log),
	}

	app := web.NewApp(mux, shutdown, log, middleware...)

	// Unauthenticated webpages.
	app.Handle(http.MethodGet, "/", withNoSession()(authHandler.LoginPage))
	app.Handle(http.MethodGet, "/force-new-password", withPasswordChallengeSession()(authHandler.ForceNewPasswordPage))
	app.Handle(http.MethodPost, "/secure-new-password", withNoSession()(authHandler.SetupNewUserWithSecurePassword))
	app.Handle(http.MethodPost, "/authenticate", withNoSession()(authHandler.AuthenticateCredentials))

	// Authenticated webpages.
	app.Handle(http.MethodGet, "/admin", withSession()(webPageHandler.DashboardPage))
	app.Handle(http.MethodGet, "/admin/tenants", withSession()(webPageHandler.TenantsPage))
	app.Handle(http.MethodGet, "/admin/create-tenant", withSession()(webPageHandler.CreateTenantPage))
	app.Handle(http.MethodGet, "/admin/logout", withSession()(authHandler.Logout))
	app.Handle(http.MethodGet, "/*", withSession()(webPageHandler.E404Page))

	// API endpoints
	app.Handle(http.MethodGet, "/admin/api/verify", authHandler.VerifyTokenNoop)
	app.Handle(http.MethodGet, "/admin/api/tenants", tenantsHandler.ListTenants)
	app.Handle(http.MethodPost, "/admin/api/send-registration", registrationHandler.ProcessRegistration)
	app.Handle(http.MethodPost, "/admin/api/resend-otp", registrationHandler.ResendOTP)

	return mux
}
