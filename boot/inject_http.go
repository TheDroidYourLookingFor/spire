package boot

import (
	"github.com/Akkadius/spire/internal/assets"
	"github.com/Akkadius/spire/internal/backup"
	"github.com/Akkadius/spire/internal/clientfiles"
	"github.com/Akkadius/spire/internal/deploy"
	"github.com/Akkadius/spire/internal/eqemuanalytics"
	"github.com/Akkadius/spire/internal/eqemuchangelog"
	"github.com/Akkadius/spire/internal/eqemuserverapi"
	apphttp "github.com/Akkadius/spire/internal/http"
	"github.com/Akkadius/spire/internal/http/controllers"
	appmiddleware "github.com/Akkadius/spire/internal/http/middleware"
	"github.com/Akkadius/spire/internal/http/routes"
	"github.com/Akkadius/spire/internal/http/staticmaps"
	"github.com/Akkadius/spire/internal/models"
	"github.com/Akkadius/spire/internal/occulus"
	"github.com/Akkadius/spire/internal/permissions"
	"github.com/Akkadius/spire/internal/questapi"
	"github.com/Akkadius/spire/internal/serverconfig"
	"github.com/Akkadius/spire/internal/spireuser"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

var httpSet = wire.NewSet(
	apphttp.NewServer,
	appmiddleware.NewUserContextMiddleware,
	appmiddleware.NewRequestLogMiddleware,
	appmiddleware.NewReadOnlyMiddleware,
	appmiddleware.NewPermissionsMiddleware,
	appmiddleware.NewLocalUserAuthMiddleware,
	controllers.NewAnalyticsController,
	controllers.NewHelloWorldController,
	controllers.NewConnectionsController,
	controllers.NewMeController,
	controllers.NewAuthController,
	controllers.NewDocsController,
	questapi.NewQuestApiController,
	controllers.NewAppController,
	controllers.NewQueryController,
	eqemuanalytics.NewAnalyticsController,
	eqemuchangelog.NewEqemuChangelogController,
	clientfiles.NewClientFilesController,
	assets.NewAssetsController,
	permissions.NewPermissionsController,
	spireuser.NewUsersController,
	occulus.NewController,
	staticmaps.NewStaticMapController,
	deploy.NewDeployController,
	eqemuserverapi.NewController,
	serverconfig.NewController,
	backup.NewController,
	provideControllers,
	NewRouter,
)

type appControllerGroups struct {
	authControllers       []routes.Controller
	helloWorldControllers []routes.Controller
	v1controllers         []routes.Controller
	v1controllersNoAuth   []routes.Controller
	v1Analytics           []routes.Controller
}

func NewRouter(
	cg *appControllerGroups,
	crudc *crudControllers,
	userContextMiddleware *appmiddleware.UserContextMiddleware,
	readOnlyModeMiddleware *appmiddleware.ReadOnlyMiddleware,
	permissionsMiddleware *appmiddleware.PermissionsMiddleware,
	logMiddleware *appmiddleware.RequestLogMiddleware,
	localUserAuthMiddleware *appmiddleware.LocalUserAuthMiddleware,
	assets *assets.SpireAssets,
) *routes.Router {
	return routes.NewHttpRouter(

		// pre middleware
		[]echo.MiddlewareFunc{
			middleware.RemoveTrailingSlash(),
		},

		// post middleware
		[]echo.MiddlewareFunc{
			logMiddleware.Handle(),
			//middleware.Logger(), // json logger
			middleware.Recover(),
			middleware.CORSWithConfig(
				middleware.CORSConfig{
					AllowOrigins: []string{"*"},
					AllowMethods: []string{
						http.MethodGet,
						http.MethodHead,
						http.MethodPut,
						http.MethodPatch,
						http.MethodPost,
						http.MethodDelete,
						http.MethodOptions,
					},
				},
			),
			assets.ServeStatic(),
		},

		// controller groups
		[]*routes.ControllerGroup{
			routes.NewControllerGroup(
				"/auth/",
				cg.authControllers,
				[]echo.MiddlewareFunc{}...,
			),
			routes.NewControllerGroup(
				"/api/v1/",
				cg.v1controllers,
				userContextMiddleware.HandleHeader(),
				userContextMiddleware.HandleQuerystring(),
				localUserAuthMiddleware.Handle(),
				readOnlyModeMiddleware.Handle(),
				v1RateLimit(),
			),
			routes.NewControllerGroup(
				"/api/v1/",
				cg.v1controllersNoAuth,
				v1RateLimit(),
				middleware.GzipWithConfig(middleware.GzipConfig{Level: 1}),
			),
			routes.NewControllerGroup(
				"/api/v1/",
				cg.v1Analytics,
				middleware.GzipWithConfig(middleware.GzipConfig{Level: 1}),
			),
			routes.NewControllerGroup(
				"/api/v1/",
				crudc.routes,
				userContextMiddleware.HandleHeader(),
				readOnlyModeMiddleware.Handle(),
				localUserAuthMiddleware.Handle(),
				permissionsMiddleware.Handle(),
				v1RateLimit(),
				middleware.GzipWithConfig(middleware.GzipConfig{Level: 1}),
			),
		},
	)
}

// controllers provider
func provideControllers(
	hello *controllers.HelloWorldController,
	auth *controllers.AuthController,
	me *controllers.MeController,
	analytics *controllers.AnalyticsController,
	connections *controllers.ConnectionsController,
	docs *controllers.DocsController,
	quest *questapi.QuestApiController,
	app *controllers.AppController,
	query *controllers.QueryController,
	clientFilesController *clientfiles.ClientFilesController,
	staticMaps *staticmaps.StaticMapController,
	analyticsController *eqemuanalytics.AnalyticsController,
	changelogController *eqemuchangelog.EqemuChangelogController,
	deployController *deploy.DeployController,
	assetsController *assets.AssetsController,
	permissionsController *permissions.PermissionsController,
	usersController *spireuser.UsersController,
	adminController *occulus.Controller,
	eqemuserverapiController *eqemuserverapi.Controller,
	serverconfigController *serverconfig.Controller,
	backupController *backup.Controller,
) *appControllerGroups {
	return &appControllerGroups{
		authControllers: []routes.Controller{
			auth,
		},
		v1controllers: []routes.Controller{
			me,
			analytics,
			connections,
			hello,
			docs,
			query,
			clientFilesController,
			permissionsController,
			usersController,
			adminController,
			eqemuserverapiController,
			serverconfigController,
			backupController,
		},
		v1controllersNoAuth: []routes.Controller{
			quest,
			app,
			staticMaps,
			assetsController,
			changelogController,
			deployController,
		},
		v1Analytics: []routes.Controller{
			analyticsController,
		},
	}
}

func v1RateLimit() echo.MiddlewareFunc {
	return appmiddleware.RateLimiterWithConfig(
		appmiddleware.RateLimiterConfig{
			Skipper: func(c echo.Context) bool {

				// if there is a valid user - skip the middleware
				user, ok := c.Get("user").(models.User)
				if ok {
					if user.ID > 0 {
						return true
					}
				}

				return false
			},
			LimitConfig: appmiddleware.LimiterConfig{
				Max:      5000,
				Duration: time.Minute * 1,
				Strategy: "ip",
				Key:      "",
			},
			Prefix:                       "LIMIT",
			Client:                       nil,
			SkipRateLimiterInternalError: false,
			OnRateLimit:                  nil,
		},
	)
}
