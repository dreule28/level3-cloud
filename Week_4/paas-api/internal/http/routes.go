package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/dreule28/Week_4/paas-api/internal/http/handlers"
	"github.com/dreule28/Week_4/paas-api/internal/service"
)

func RegisterRoutes(e *echo.Echo, svc service.InstanceAPI) {
	h := handlers.NewInstanceHandler(svc)

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "You're amazin")
	})

	// CPU-burn endpoint for HPA demo
	// GET /work?ms=50
	e.GET("/work", func(c echo.Context) error {
		ms, err := strconv.Atoi(c.QueryParam("ms"))
		if err != nil || ms <= 0 {
			ms = 50
		}
		if ms > 2000 { // safety cap so you don't melt your cluster
			ms = 2000
		}

		deadline := time.Now().Add(time.Duration(ms) * time.Millisecond)

		x := 0
		for time.Now().Before(deadline) {
			x++
		}

		return c.JSON(http.StatusOK, map[string]any{
			"ok":    true,
			"iters": x,
			"ms":    ms,
		})
	})

	g := e.Group("/instances")
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.DELETE("/:id", h.Delete)
}
