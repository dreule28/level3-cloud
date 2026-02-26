package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/dreule28/Week_4/paas-api/internal/config"
	"github.com/dreule28/Week_4/paas-api/internal/http/auth"
	"github.com/dreule28/Week_4/paas-api/internal/http/handlers"
	"github.com/dreule28/Week_4/paas-api/internal/service"
)

func RegisterRoutes(e *echo.Echo, svc service.InstanceAPI, cfg config.Config) {
	h := handlers.NewInstanceHandler(svc)

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "You're amazin")
	})

	authCfg := auth.Config{
		AuthUser:    cfg.AuthUser,
		AuthPass:    cfg.AuthPass,
		JWTSecret:   []byte(cfg.JWTSecret),
		JWTIssuer:   cfg.JWTIssuer,
		JWTAudience: cfg.JWTAudience,
		JWTTL:       cfg.JWTTL,
	}

	// Auth endpoints
	auth.RegisterAuthRoutes(e, authCfg)

	// Protected routes
	protected := e.Group("")
	protected.Use(auth.RequireJWT(authCfg))
	protected.GET("/logs", h.ListAllLogs)

	// CPU-burn endpoint for HPA demo (protected)
	// GET /work?ms=50
	protected.GET("/work", func(c echo.Context) error {
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

	instances := protected.Group("/instances")
	instances.GET("", h.List)
	instances.GET("/:id/logs", h.ListLogs)
	instances.GET("/:id", h.Get)

	admin := instances.Group("")
	admin.Use(auth.RequireRole("admin"))
	admin.POST("", h.Create)
	admin.DELETE("/:id", h.Delete)
}
