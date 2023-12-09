package auth

import (
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/auth"
	"github.com/volatiletech/authboss/v3/defaults"
	"github.com/volatiletech/authboss/v3/logout"
	"os"
)

var AuthConfig *authboss.Authboss
var Auth *auth.Auth
var Logout *logout.Logout

//var Register *register.Register

func init() {
	AuthConfig = authboss.New()
	AuthConfig.Config.Core.ViewRenderer = &Renderer{}
	AuthConfig.Config.Core.Responder = &Responder{}

	//ab.Config.Storage.Server = myDatabaseImplementation
	//ab.Config.Storage.SessionState = mySessionImplementation
	//ab.Config.Storage.CookieState = myCookieImplementation

	// This instantiates and uses every default implementation
	// in the Config.Core area that exist in the defaults package.
	// Just a convenient helper if you don't want to do anything fancy.
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

	Auth = &auth.Auth{AuthConfig}
	Logout = &logout.Logout{AuthConfig}
	//Register = &register.Register{AuthConfig}
}
