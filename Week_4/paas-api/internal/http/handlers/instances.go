package handlers

import (
	"net/http"


	"github.com/dreule28/Week_4/paas-api/internal/service"
	"github.com/labstack/echo/v4"
)

type InstancesHandler struct {
	svc *service.InstanceService
}

func (h *InstancesHandler) List(c echo.Context) error {
	items, err := h.svc.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

func NewInstanceHandler(svc *service.InstanceService) *InstancesHandler {
	return &InstancesHandler{svc: svc}
}
