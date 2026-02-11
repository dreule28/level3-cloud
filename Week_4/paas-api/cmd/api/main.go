package main

import (
	"log"

	"github.com/dreule28/Week_4/paas-api/internal/config"
	"github.com/dreule28/Week_4/paas-api/internal/http"
	"github.com/dreule28/Week_4/paas-api/internal/kube"
	"github.com/dreule28/Week_4/paas-api/internal/service"
	// "sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/controller-runtime/pkg/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

)

func main() {
	cfg := config.MustLoad()

	k8sclient, err := kube.NewClient(cfg)
	if err != nil {
		log.Fatalf("kube client: %v", err)
	}

	svc := service.NewInstanceService(cfg, k8sclient)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	http.RegisterRoutes(e, svc)

	log.Fatal(e.Start(":8080"))
}