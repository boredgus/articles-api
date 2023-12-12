package internal

import (
	"net/http"

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
		return app.Login.Register(NewContext(c))
	})
	e.POST("/authorize", func(c echo.Context) error {
		return app.Login.Authorize(NewContext(c))
	})
	e.GET("/articles", func(c echo.Context) error {
		return app.Article.GetForUser(NewContext((c)))
	})
	e.GET("/articles/:article_id", func(c echo.Context) error {
		return app.Article.Get(NewContext(c))
	})

	protectedArticles := e.Group("/articles", jwtAuthMiddleware())
	protectedArticles.POST("", func(c echo.Context) error {
		return app.Article.Create(NewContext(c))
	})
	protectedArticles.PUT("/:article_id", func(c echo.Context) error {
		return app.Article.Update(NewContext(c))
	})
	protectedArticles.DELETE("/:article_id", func(c echo.Context) error {
		return app.Article.Delete(NewContext(c))
	})

	return e
}

func GetRouter(cntrs AppController) *echo.Echo {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogError:  true,
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

	registerRoutes(e, cntrs)

	return e
}
