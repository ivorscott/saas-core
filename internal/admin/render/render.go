// Package render manages the web app page rendering logic.
package render

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/devpies/saas-core/internal/admin/config"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/alexedwards/scs/v2"

	"go.uber.org/zap"
)

// Render contains various methods for rendering .gohtml templates.
type Render struct {
	logger     *zap.Logger
	cfg        config.Config
	cache      templateCache
	templateFS fs.FS
	session    *scs.SessionManager
}

type templateCache map[string]*template.Template

// TemplateData models data options in templates.
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	Data            map[string]interface{}
	IsAuthenticated int
	UserID          string
	Email           string
	API             string
}

var functions = template.FuncMap{
	"formatCurrency": formatCurrency,
}

var currencyMap = map[string]string{
	"eur": "€",
	"usd": "$",
}

func formatCurrency(n int, currency string) string {
	f := float32(n) / float32(100)
	return fmt.Sprintf("%s %.2f", currencyMap[currency], f)
}

// New returns a new Render with template rendering logic.
func New(logger *zap.Logger, config config.Config, templateFS fs.FS, session *scs.SessionManager) *Render {
	cache := make(templateCache)
	return &Render{
		logger:     logger,
		cfg:        config,
		cache:      cache,
		templateFS: templateFS,
		session:    session,
	}
}

// AddDefaultData provides .gohtml templates with default data.
func (re *Render) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.API = fmt.Sprintf("%s:%s", re.cfg.Web.Address, re.cfg.Web.Port)
	td.IsAuthenticated = 0
	td.UserID = ""
	td.Email = ""

	if re.session.Exists(r.Context(), "UserID") {
		td.IsAuthenticated = 1
		td.UserID = re.session.GetString(r.Context(), "UserID")
		td.Email = re.session.GetString(r.Context(), "Email")
	}

	return td
}

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Template renders a template for the application.
// During development, renderTemplate will never use the template cache.
func (re *Render) Template(
	w http.ResponseWriter,
	r *http.Request,
	page string,
	td *TemplateData,
	partials ...string) error {
	var t *template.Template
	var err error

	// The template to render.
	tmpl := fmt.Sprintf("%s.page.gohtml", page)

	if val, ok := re.cache[tmpl]; ok {
		t = val
	} else {
		t, err = re.parseTemplate(partials, page, tmpl)
		if err != nil {
			re.logger.Error("", zap.Error(err))
			return web.NewShutdownError(err.Error())
		}
	}

	if td == nil {
		td = &TemplateData{}
	}

	td = re.AddDefaultData(td, r)
	err = t.Execute(w, td)
	if err != nil {
		re.logger.Error("", zap.Error(err))
		return web.NewShutdownError(err.Error())
	}

	web.SetContextStatusCode(r.Context(), http.StatusOK)

	return nil
}

// parseTemplate parses the desired page template with or with partials.
func (re *Render) parseTemplate(partials []string, page, tmpl string) (*template.Template, error) {
	var t *template.Template
	var err error

	// Retrieve partial templates if they exist.
	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("%s.partial.gohtml", x)
		}

		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).
			ParseFS(re.templateFS, "base.layout.gohtml", strings.Join(partials, ","), tmpl)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).
			ParseFS(re.templateFS, "base.layout.gohtml", tmpl)
	}

	re.cache[tmpl] = t

	return t, err
}
