package service

import (
	"context"

	"github.com/dreule28/Week_4/paas-api/internal/model"
)

type InstanceAPI interface {
	ListDatabases(ctx context.Context) ([]model.Instance, error)
	GetDatabase(ctx context.Context, id string) (model.InstanceDetails, error)
	CreateDatabase(ctx context.Context, req model.CreateInstanceRequest) (model.Instance, error)
	DeleteDatabase(ctx context.Context, id string) error
}

