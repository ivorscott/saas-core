package render

import (
	"embed"
	"fmt"
	"github.com/devpies/core/admin/pkg/config"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"strings"
)

type Render struct {
	logger     *zap.Logger
	cfg        config.Config
	cache      map[string]*template.Template
	templateFS embed.FS
}

type templateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	Data            map[string]interface{}
	IsAuthenticated int
	UserID          int
	API             string
}

// functions for templates.
var functions = template.FuncMap{}

func New(logger *zap.Logger, config config.Config, templateCache map[string]*template.Template, templateFS embed.FS) *Render {
	return &Render{
		logger:     logger,
		cfg:        config,
		cache:      templateCache,
		templateFS: templateFS,
	}
}

func (re *Render) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.API = re.cfg.Web.APIAddress
	td.IsAuthenticated = 0
	td.UserID = 0

	//if app.Session.Exists(r.Context(), "userID") {
	//	td.IsAuthenticated = 1
	//	td.UserID = app.Session.GetInt(r.Context(), "userID")
	//}

	return td
}

// renderTemplate renders a template for the application.
// During development, renderTemplate will never use the template cache.
func (re *Render) renderTemplate(
	w http.ResponseWriter,
	r *http.Request,
	page string,
	td *templateData,
	partials ...string) error {
	var t *template.Template
	var err error

	// The template to render.
	tmpl := fmt.Sprintf("templates/%s.page.gohtml", page)

	if val, ok := re.cache[tmpl]; ok {
		t = val
	} else {
		t, err = re.parseTemplate(partials, page, tmpl)
		if err != nil {
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}

	td = re.addDefaultData(td, r)
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
			partials[i] = fmt.Sprintf("templates/%s.partial.gohtml", x)
		}

		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).
			ParseFS(re.templateFS, "templates/base.layout.gohtml", strings.Join(partials, ","), tmpl)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).
			ParseFS(re.templateFS, "templates/base.layout.gohtml", tmpl)
	}

	re.cache[tmpl] = t

	return t, err
}
