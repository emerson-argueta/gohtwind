package auth

import (
	"context"
	"github.com/volatiletech/authboss/v3"
	"net/http"
	"town/infra"
)

type Responder struct {
	vt *infra.ViewTemplate
}

func (res *Responder) Respond(w http.ResponseWriter, r *http.Request, code int, templateName string, data authboss.HTMLData) error {
	res.vt = &infra.ViewTemplate{
		BasePath: "templates",
		Path:     "auth/templates/" + templateName + ".html",
	}
	fv, err := infra.NewView(res.vt)
	if err != nil {
		return err
	}
	fv.RenderTemplate(w, data)
	return nil
}

type Renderer struct{}

func (r *Renderer) Render(ctx context.Context, name string, data authboss.HTMLData) (output []byte, contentType string, err error) {
	return nil, "text/html", nil
}

func (r *Renderer) Load(names ...string) error {
	return nil
}
