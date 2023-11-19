package infrastructure

import (
	"net/http"
	"user-management/internal/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func registerRoutes(e *echo.Echo, app controllers.AppController) *echo.Echo {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ message string }{message: "alive"})
	})
	e.POST("/register", func(c echo.Context) error {
		return app.User.Register(NewContext(c))
	})
	e.GET("/authorize", func(c echo.Context) error {
		return app.User.Authorize(NewContext(c))
	})

	return e
}

func GetRouter(cntrs controllers.AppController) *echo.Echo {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} ${method} ${uri} ${status} error:${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logrus.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))

	registerRoutes(e, cntrs)

	return e
}
