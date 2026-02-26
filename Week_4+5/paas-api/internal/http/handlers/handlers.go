package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	apiauth "github.com/dreule28/Week_4/paas-api/internal/http/auth"
	"github.com/dreule28/Week_4/paas-api/internal/model"
	"github.com/dreule28/Week_4/paas-api/internal/service"
	"github.com/labstack/echo/v4"
)

var validID = regexp.MustCompile(`^[a-z0-9]([a-z0-9\-]{0,61}[a-z0-9])?$`)

type InstancesHandler struct {
	svc  service.InstanceAPI
	logs service.LogsAPI
}

func NewInstanceHandler(svc service.InstanceAPI) *InstancesHandler {
	h := &InstancesHandler{svc: svc}
	if logsSvc, ok := svc.(service.LogsAPI); ok {
		h.logs = logsSvc
	}
	return h
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
		if errors.Is(err, service.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

func (h *InstancesHandler) ListLogs(c echo.Context) error {
	if h.logs == nil {
		return c.JSON(http.StatusNotImplemented, map[string]string{"error": "log retrieval not supported"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	logType := strings.ToLower(strings.TrimSpace(c.QueryParam("type")))
	switch logType {
	case "", model.LogTypeAudit, model.LogTypeService:
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": `type must be "audit" or "service"`})
	}

	limit := 100
	if raw := strings.TrimSpace(c.QueryParam("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "limit must be a positive integer"})
		}
		limit = n
	}

	items, err := h.logs.ListInstanceLogs(c.Request().Context(), model.LogQuery{
		InstanceID: id,
		Type:       logType,
		Limit:      limit,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, items)
}

func (h *InstancesHandler) ListAllLogs(c echo.Context) error {
	if h.logs == nil {
		return c.JSON(http.StatusNotImplemented, map[string]string{"error": "log retrieval not supported"})
	}

	logType := strings.ToLower(strings.TrimSpace(c.QueryParam("type")))
	switch logType {
	case "", model.LogTypeAudit, model.LogTypeService:
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": `type must be "audit" or "service"`})
	}

	limit := 100
	if raw := strings.TrimSpace(c.QueryParam("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "limit must be a positive integer"})
		}
		limit = n
	}

	items, err := h.logs.ListInstanceLogs(c.Request().Context(), model.LogQuery{
		InstanceID: strings.TrimSpace(c.QueryParam("instanceId")),
		Type:       logType,
		Limit:      limit,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, items)
}

func (h *InstancesHandler) Create(c echo.Context) error {
	var req model.CreateInstanceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
	}

	if req.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	if !validID.MatchString(req.ID) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id must be a valid RFC 1123 DNS label"})
	}
	if req.Instances <= 0 {
		req.Instances = 1
	}
	if req.StorageGi <= 0 {
		req.StorageGi = 10
	}

	out, err := h.svc.CreateDatabase(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.recordAudit(c, req.ID, "instance.create.requested", "User requested instance creation via API/UI")
	return c.JSON(http.StatusAccepted, out)
}

func (h *InstancesHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	if err := h.svc.DeleteDatabase(c.Request().Context(), id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.recordAudit(c, id, "instance.delete.requested", "User requested instance deletion via API/UI")
	return c.NoContent(http.StatusNoContent)
}

func (h *InstancesHandler) recordAudit(c echo.Context, instanceID, action, message string) {
	if h.logs == nil {
		return
	}

	user := "unknown"
	if claims := apiauth.GetClaims(c); claims != nil && claims.Subject != "" {
		user = claims.Subject
	}

	_ = h.logs.RecordAuditLog(c.Request().Context(), instanceID, user, action, message)
}
