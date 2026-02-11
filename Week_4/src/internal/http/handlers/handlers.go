package handlers

import (
	"net/http"
	"strings"

	"github.com/dreule28/Week_4/paas-api/internal/model"
	"github.com/dreule28/Week_4/paas-api/internal/service"
	"github.com/labstack/echo/v4"
)

type InstancesHandler struct {
	svc service.InstanceAPI
}

func NewInstanceHandler(svc service.InstanceAPI) *InstancesHandler {
	return &InstancesHandler{svc: svc}
}

func (h *InstancesHandler) List(c echo.Context) error {
	items, err := h.svc.ListDatabases(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

func (h *InstancesHandler) Get(c echo.Context) error {
	id := c.Param("id")

	out, err := h.svc.GetDatabase(c.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
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

	out, err := h.svc.CreateDatabase(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusAccepted, out)
}

func (h *InstancesHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	if err := h.svc.DeleteDatabase(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
