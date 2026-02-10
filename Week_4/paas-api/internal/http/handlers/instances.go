package handlers

import (
	"net/http"

	"github.com/dreule28/Week_4/paas-api/internal/model"
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

func (h *InstancesHandler) Get(c echo.Context) error {
	id := c.Param("id")

	out, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

func (h *InstancesHandler) Create(c echo.Context) error {
	var req model.CreateInstanceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
	}

	if req.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	if req.Instances <= 0 {
		req.Instances = 1
	}
	if req.StorageGi <= 0 {
		req.StorageGi = 10
	}

	out, err := h.svc.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusAccepted, out)
}

func NewInstanceHandler(svc *service.InstanceService) *InstancesHandler {
	return &InstancesHandler{svc: svc}
}
