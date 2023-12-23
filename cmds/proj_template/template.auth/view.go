package auth

import (
	"context"
	"fmt"
	"github.com/volatiletech/authboss/v3"
	"net/http"
	"{{PROJECT_NAME}}/infra"
)

type Responder struct {
	authboss.HTTPResponder
	vt *infra.ViewTemplate
}

func (resp *Responder) Respond(w http.ResponseWriter, r *http.Request, code int, templateName string, data authboss.HTMLData) error {
	path := fmt.Sprintf("auth/templates/%s.html", templateName)
	resp.vt.Path = path
	return infra.RenderTemplate(w, resp.vt, data)
}

type Renderer struct{ authboss.Renderer }

func (r *Renderer) Load(names ...string) error {
	return nil
}

func (r *Renderer) Render(ctx context.Context, page string, data authboss.HTMLData) (output []byte, contentType string, err error) {
	return nil, "text/html", nil
}
