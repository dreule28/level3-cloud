package http

import (
	"github.com/labstack/echo/v4"

	"github.com/dreule28/Week_4/paas-api/internal/http/handlers"
	"github.com/dreule28/Week_4/paas-api/internal/service"
)

func RegisterRoutes(e *echo.Echo, svc service.InstanceAPI) {
	h := handlers.NewInstanceHandler(svc)

	e.GET("/healthz", func(c echo.Context) error { return c.String(200, "You're amazin") })

	g := e.Group("/instances")
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.DELETE("/:id", h.Delete)
}