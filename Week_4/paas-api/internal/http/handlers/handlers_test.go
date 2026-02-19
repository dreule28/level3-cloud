package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/dreule28/Week_4/paas-api/internal/http/handlers"
	"github.com/dreule28/Week_4/paas-api/internal/model"
	"github.com/dreule28/Week_4/paas-api/internal/service"
)

// fakeSvc is a mock implementation of service.InstanceAPI.
// Each field is a function that lets us control what the service returns
// per test, without needing a real Kubernetes cluster.
type fakeSvc struct {
	listFn   func(context.Context) ([]model.Instance, error)
	getFn    func(context.Context, string) (model.InstanceDetails, error)
	createFn func(context.Context, model.CreateInstanceRequest) (model.Instance, error)
	deleteFn func(context.Context, string) error
}

func (f fakeSvc) ListDatabases(ctx context.Context) ([]model.Instance, error) {
	return f.listFn(ctx)
}
func (f fakeSvc) GetDatabase(ctx context.Context, id string) (model.InstanceDetails, error) {
	return f.getFn(ctx, id)
}
func (f fakeSvc) CreateDatabase(ctx context.Context, req model.CreateInstanceRequest) (model.Instance, error) {
	return f.createFn(ctx, req)
}
func (f fakeSvc) DeleteDatabase(ctx context.Context, id string) error {
	return f.deleteFn(ctx, id)
}

// Compile-time check: fakeSvc must satisfy the InstanceAPI interface.
// If fakeSvc is missing a method, this line will cause a build error.
var _ service.InstanceAPI = fakeSvc{}

// ---------------------------------------------------------------------------
// Happy Path Tests — everything works as expected
// ---------------------------------------------------------------------------

// TestList_OK verifies that when the service returns instances,
// the handler responds with 200 and a JSON array containing them.
func TestList_OK(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		listFn: func(ctx context.Context) ([]model.Instance, error) {
			return []model.Instance{{ID: "pg-demo", Status: "ready"}}, nil
		},
		getFn: func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/instances", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.List(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var got []model.Instance
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got) != 1 || got[0].ID != "pg-demo" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

// TestGet_OK verifies that the Get handler extracts the :id path
// parameter, calls the service, and returns 200 with instance details.
func TestGet_OK(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn: func(ctx context.Context, id string) (model.InstanceDetails, error) {
			return model.InstanceDetails{ID: id, Status: "ready"}, nil
		},
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/instances/pg-demo", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/instances/:id")
	c.SetParamNames("id")
	c.SetParamValues("pg-demo")

	if err := h.Get(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

// ---------------------------------------------------------------------------
// Bad Input Tests — client sends invalid data
// ---------------------------------------------------------------------------

// TestCreate_BadJSON verifies that sending malformed JSON returns 400
// and the service is never called. The t.Fatal inside createFn acts
// as a guard — if the service IS called, the test fails immediately.
func TestCreate_BadJSON(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn:  func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			t.Fatal("should not be called")
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	})

	req := httptest.NewRequest(http.MethodPost, "/instances", bytes.NewBufferString("{bad json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Create(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestCreate_InvalidID_400 verifies that an ID containing uppercase
// letters or special characters is rejected with 400 Bad Request.
func TestCreate_InvalidID_400(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		listFn:   func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn:    func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			t.Fatal("should not be called")
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	})

	req := httptest.NewRequest(http.MethodPost, "/instances", bytes.NewBufferString(`{"id":"INVALID ID!"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Create(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestDelete_MissingID verifies that calling delete without providing
// an :id parameter returns 400. The service should never be reached.
func TestDelete_MissingID(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn:  func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { t.Fatal("should not be called"); return nil },
	})

	req := httptest.NewRequest(http.MethodDelete, "/instances/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// no id param set — simulates a request with no :id in the path

	if err := h.Delete(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

// ---------------------------------------------------------------------------
// Service Error Tests — the backend (Kubernetes) fails
// ---------------------------------------------------------------------------

// TestGet_ServiceError_500 verifies that when the service returns a
// generic error (e.g. Kubernetes API unreachable), the handler responds
// with 500 Internal Server Error.
func TestGet_ServiceError_500(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn: func(context.Context, string) (model.InstanceDetails, error) {
			return model.InstanceDetails{}, fmt.Errorf("boom")
		},
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/instances/pg-demo", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/instances/:id")
	c.SetParamNames("id")
	c.SetParamValues("pg-demo")

	if err := h.Get(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestList_ServiceError_500 verifies that when ListDatabases fails,
// the handler returns 500 instead of crashing or returning partial data.
func TestList_ServiceError_500(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) {
			return nil, fmt.Errorf("boom")
		},
		getFn: func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/instances", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.List(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

// TestDelete_ServiceError_500 verifies that when DeleteDatabase fails,
// the handler returns 500. This catches cases like Kubernetes rejecting
// the delete or the cluster being unreachable.
func TestDelete_ServiceError_500(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn:  func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return fmt.Errorf("boom") },
	})

	req := httptest.NewRequest(http.MethodDelete, "/instances/pg-demo", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/instances/:id")
	c.SetParamNames("id")
	c.SetParamValues("pg-demo")

	if err := h.Delete(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// PaaS Edge Case Tests
// ---------------------------------------------------------------------------

// TestList_Empty_ReturnsEmptyArray verifies that when there are zero
// instances, the handler returns "[]" (empty JSON array) and NOT "null".
// This matters because frontend consumers (JavaScript, etc.) will crash
// on null.length but handle [].length just fine.
func TestList_Empty_ReturnsEmptyArray(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		// Return an empty slice (not nil!) so JSON encodes as []
		listFn: func(ctx context.Context) ([]model.Instance, error) {
			return []model.Instance{}, nil
		},
		getFn: func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/instances", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.List(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Must be "[]", not "null"
	body := strings.TrimSpace(rec.Body.String())
	if body != "[]" {
		t.Fatalf("expected empty JSON array '[]', got %q", body)
	}
}

// TestGet_NotFound_404 verifies that when the service returns a
// "not found" error (e.g. the Kubernetes CR doesn't exist), the handler
// responds with 404 instead of a generic 500. This lets PaaS consumers
// distinguish "instance does not exist" from "server is broken."
func TestGet_NotFound_404(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn: func(ctx context.Context, id string) (model.InstanceDetails, error) {
			return model.InstanceDetails{}, fmt.Errorf("instance %q: %w", id, service.ErrNotFound)
		},
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/instances/does-not-exist", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/instances/:id")
	c.SetParamNames("id")
	c.SetParamValues("does-not-exist")

	if err := h.Get(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestGet_ConnectionInfo_Present verifies that when a database instance
// is "ready", the JSON response includes all connection details:
// host, port, database, user, password, and endpoint.
// This is the core value of a PaaS — users need this info to connect
// to their provisioned database.
func TestGet_ConnectionInfo_Present(t *testing.T) {
	e := echo.New()

	h := handlers.NewInstanceHandler(fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn: func(ctx context.Context, id string) (model.InstanceDetails, error) {
			return model.InstanceDetails{
				ID:     id,
				Status: "ready",
				Connection: &model.ConnectionInfo{
					Host:     "pg-demo-rw.default.svc.cluster.local",
					Port:     5432,
					Database: "app",
					User:     "app",
					Password: "secret",
					Endpoint: "postgres://app@pg-demo-rw.default.svc.cluster.local:5432/app",
				},
			}, nil
		},
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/instances/pg-demo", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/instances/:id")
	c.SetParamNames("id")
	c.SetParamValues("pg-demo")

	if err := h.Get(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Unmarshal and verify connection fields are present and populated
	var got model.InstanceDetails
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.Connection == nil {
		t.Fatal("expected connection info to be present, got nil")
	}
	if got.Connection.Host == "" {
		t.Fatal("expected connection host to be set")
	}
	if got.Connection.Port == 0 {
		t.Fatal("expected connection port to be set")
	}
	if got.Connection.Endpoint == "" {
		t.Fatal("expected connection endpoint to be set")
	}
	if got.Connection.User == "" {
		t.Fatal("expected connection user to be set")
	}
	if got.Connection.Password == "" {
		t.Fatal("expected connection password to be set")
	}
	if got.Connection.Database == "" {
		t.Fatal("expected connection database to be set")
	}
}
