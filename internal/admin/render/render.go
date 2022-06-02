package render

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/devpies/core/internal/adminclient/config"
)

type Render struct {
	logger     *zap.Logger
	cfg        config.Config
	cache      templateCache
	templateFS fs.FS
}

type templateCache map[string]*template.Template

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	Data            map[string]interface{}
	IsAuthenticated int
	UserID          int
	API             string
}

// functions for templates.
var functions = template.FuncMap{}

func New(logger *zap.Logger, config config.Config, templateFS fs.FS) *Render {
	cache := make(templateCache)
	return &Render{
		logger:     logger,
		cfg:        config,
		cache:      cache,
		templateFS: templateFS,
	}
}

func (re *Render) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.API = re.cfg.Web.Backend
	td.IsAuthenticated = 0
	td.UserID = 0

	//if app.Session.Exists(r.Context(), "userID") {
	//	td.IsAuthenticated = 1
	//	td.UserID = app.Session.GetInt(r.Context(), "userID")
	//}

	return td
}

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
			return err
		}
	}

	if td == nil {
		td = &TemplateData{}
	}

	td = re.AddDefaultData(td, r)
	err = t.Execute(w, td)
	if err != nil {
		re.logger.Error("", zap.Error(err))
		return err
	}

	return err
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
