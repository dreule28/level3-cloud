package http_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/dreule28/Week_4/paas-api/internal/config"
	myhttp "github.com/dreule28/Week_4/paas-api/internal/http"
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
var _ service.InstanceAPI = fakeSvc{}

// TestHealthz_OK verifies that the /healthz endpoint returns 200
// and the expected body. This is the simplest smoke test to confirm
// the API is alive and the router is wired up.
func TestHealthz_OK(t *testing.T) {
	e := echo.New()

	svc := fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn:  func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	}

	cfg := config.Config{}
	myhttp.RegisterRoutes(e, svc, cfg)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "You're amazin") {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

// TestGetInstance_Route_OK verifies that GET /instances/:id reaches
// the correct handler and returns the instance JSON with a 200 status.
func TestGetInstance_Route_OK(t *testing.T) {
	e := echo.New()

	svc := fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn: func(ctx context.Context, id string) (model.InstanceDetails, error) {
			return model.InstanceDetails{ID: id, Status: "ready"}, nil
		},
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	}

	cfg := config.Config{}
	myhttp.RegisterRoutes(e, svc, cfg)

	req := httptest.NewRequest(http.MethodGet, "/instances/pg-demo", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"id":"pg-demo"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

// TestPostInstances_Route_Accepted verifies that POST /instances with
// a valid JSON body reaches the service and returns 202 Accepted.
// It also checks that the request payload (id) arrives correctly.
func TestPostInstances_Route_Accepted(t *testing.T) {
	e := echo.New()

	svc := fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn:  func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(ctx context.Context, req model.CreateInstanceRequest) (model.Instance, error) {
			// verify request reached service
			if req.ID != "pg-new" {
				t.Fatalf("expected id pg-new, got %q", req.ID)
			}
			return model.Instance{ID: req.ID, Status: "creating"}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	}

	cfg := config.Config{}
	myhttp.RegisterRoutes(e, svc, cfg)

	body := `{"id":"pg-new","instances":1,"storageGi":10}`
	req := httptest.NewRequest(http.MethodPost, "/instances", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"id":"pg-new"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

// TestDeleteInstances_Route_NoContent verifies that DELETE /instances/:id
// calls the delete service with the correct ID and returns 204 No Content.
// Uses a "called" flag to ensure the service was actually invoked.
func TestDeleteInstances_Route_NoContent(t *testing.T) {
	e := echo.New()

	called := false
	svc := fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn:  func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(ctx context.Context, id string) error {
			called = true
			if id != "pg-del" {
				t.Fatalf("expected id pg-del, got %q", id)
			}
			return nil
		},
	}

	cfg := config.Config{}
	myhttp.RegisterRoutes(e, svc, cfg)

	req := httptest.NewRequest(http.MethodDelete, "/instances/pg-del", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected delete service to be called")
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", rec.Code, rec.Body.String())
	}
}

// ---------------------------------------------------------------------------
// PaaS Edge Case Tests
// ---------------------------------------------------------------------------

// TestListInstances_Route_EmptyArray verifies that when there are zero
// instances, the API returns an empty JSON array "[]" and NOT "null".
// This is important because frontend consumers like JavaScript will
// crash on null.length but handle [].length fine.
func TestListInstances_Route_EmptyArray(t *testing.T) {
	e := echo.New()

	svc := fakeSvc{
		// Return an empty slice (not nil) to ensure JSON encodes as []
		listFn: func(context.Context) ([]model.Instance, error) {
			return []model.Instance{}, nil
		},
		getFn: func(context.Context, string) (model.InstanceDetails, error) { return model.InstanceDetails{}, nil },
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	}

	cfg := config.Config{}
	myhttp.RegisterRoutes(e, svc, cfg)

	req := httptest.NewRequest(http.MethodGet, "/instances", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// The body must be "[]" (empty JSON array), never "null"
	body := strings.TrimSpace(rec.Body.String())
	if body != "[]" {
		t.Fatalf("expected empty JSON array '[]', got %q", body)
	}
}

// TestGetInstance_Route_NotFound verifies that requesting a non-existent
// instance returns 404 Not Found instead of a generic 500. This is critical
// for PaaS consumers so they can distinguish "does not exist" from
// "something broke on the server."
func TestGetInstance_Route_NotFound(t *testing.T) {
	e := echo.New()

	svc := fakeSvc{
		listFn: func(context.Context) ([]model.Instance, error) { return nil, nil },
		getFn: func(ctx context.Context, id string) (model.InstanceDetails, error) {
			return model.InstanceDetails{}, fmt.Errorf("instance %q not found", id)
		},
		createFn: func(context.Context, model.CreateInstanceRequest) (model.Instance, error) {
			return model.Instance{}, nil
		},
		deleteFn: func(context.Context, string) error { return nil },
	}

	cfg := config.Config{}
	myhttp.RegisterRoutes(e, svc, cfg)

	req := httptest.NewRequest(http.MethodGet, "/instances/does-not-exist", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestGetInstance_Route_ConnectionInfo verifies that when a database
// instance is "ready", the response includes connection details
// (host, port, endpoint, credentials). This is the core value of a
// PaaS API — users need this info to connect to their database.
func TestGetInstance_Route_ConnectionInfo(t *testing.T) {
	e := echo.New()

	svc := fakeSvc{
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
	}

	cfg := config.Config{}
	myhttp.RegisterRoutes(e, svc, cfg)

	req := httptest.NewRequest(http.MethodGet, "/instances/pg-demo", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()

	// Verify all essential connection fields are present in the JSON response
	for _, field := range []string{`"host"`, `"port"`, `"endpoint"`, `"user"`, `"password"`, `"database"`} {
		if !strings.Contains(body, field) {
			t.Fatalf("expected connection field %s in response body, got: %s", field, body)
		}
	}
}
