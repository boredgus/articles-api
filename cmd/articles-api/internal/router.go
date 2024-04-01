package internal

import (
	"net/http"
	"strings"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func registerRoutes(e *echo.Echo, app AppController) *echo.Echo {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Message string `json:"message"`
		}{Message: "alive"})
	})
	e.POST("/register", func(c echo.Context) error {
		return app.User.Register(NewContext(c))
	})
	e.POST("/authorize", func(c echo.Context) error {
		return app.User.Authorize(NewContext(c))
	})
	e.GET("/articles", func(c echo.Context) error {
		return app.Article.GetForUser(NewContext((c)))
	})
	e.GET("/articles/:article_id", func(c echo.Context) error {
		return app.Article.Get(NewContext(c))
	})

	protected := e.Group("", jwtAuthMiddleware())
	protected.POST("/articles", func(c echo.Context) error {
		return app.Article.Create(NewContext(c))
	})
	protected.PUT("/articles/:article_id", func(c echo.Context) error {
		return app.Article.Update(NewContext(c))
	})
	protected.PUT("/articles/:article_id/reaction", func(c echo.Context) error {
		return app.Article.UpdateReactionForArticle(NewContext(c))
	})
	protected.DELETE("/articles/:article_id", func(c echo.Context) error {
		return app.Article.Delete(NewContext(c))
	})
	protected.DELETE("/users/:user_id", func(c echo.Context) error {
		return app.User.Delete(NewContext(c))
	})
	protected.PATCH("/users/:user_id/role", func(c echo.Context) error {
		return app.User.UpdateRole(NewContext(c))
	})

	return e
}

const MetricsPath = "/metrics"

func GetRouter(cntrs AppController) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogError:  true,
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), MetricsPath)
		},
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logrus.WithFields(logrus.Fields{
				"method": values.Method,
				"URI":    values.URI,
				"status": values.Status,
				"error":  values.Error,
			}).Info("request")

			return nil
		},
	}))

	e.Use(middleware.CORS())

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Subsystem: "articles_service",
	}))
	e.GET(MetricsPath, echoprometheus.NewHandler())

	registerRoutes(e, cntrs)

	return e
}
