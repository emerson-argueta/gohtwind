package auth

import (
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/auth"
	"github.com/volatiletech/authboss/v3/defaults"
	"github.com/volatiletech/authboss/v3/logout"
	"net/http"
	"os"
	"{{PROJECT_NAME}}/infra"
)

// Duration of the session cookie in hours
const SessionDuration = 24

var (
	AuthConfig        *authboss.Authboss
	Auth              *auth.Auth
	Logout            *logout.Logout
	ProtectMiddleware func(http.Handler) http.Handler
	sessionMiddleware func(http.Handler) http.Handler
	//Register *register.Register
)

func init() {
	AuthConfig = authboss.New()
	AuthConfig.Config.Core.ViewRenderer = &Renderer{}
	AuthConfig.Config.Core.Responder = &Responder{
		vt: &infra.ViewTemplate{
			BasePath: "templates",
		},
	}
	//AuthConfig.Config.Storage.Server = myDatabaseImplementation
	AuthConfig.Config.Storage.SessionState = &SessionState{}
	//AuthConfig.Config.Storage.CookieState = myCookieImplementation
	//defaults.SetCore(&AuthConfig.Config, false, false)
	logger := defaults.NewLogger(os.Stdout)
	AuthConfig.Core.Router = defaults.NewRouter()
	AuthConfig.Core.ErrorHandler = defaults.NewErrorHandler(logger)
	AuthConfig.Core.Redirector = defaults.NewRedirector(AuthConfig.Core.ViewRenderer, authboss.FormValueRedirect)
	AuthConfig.Core.BodyReader = defaults.NewHTTPBodyReader(false, false)
	AuthConfig.Core.Mailer = defaults.NewLogMailer(os.Stdout)
	AuthConfig.Core.Logger = logger
	if err := AuthConfig.Init(); err != nil {
		panic(err)
	}
	ProtectMiddleware = authboss.Middleware2(AuthConfig, authboss.RequireNone, authboss.RespondRedirect)
	sessionMiddleware = AuthConfig.LoadClientStateMiddleware
	Auth = &auth.Auth{AuthConfig}
	Logout = &logout.Logout{AuthConfig}
	//Register = &register.Register{AuthConfig}
}
