package service

import (
	"context"

	"github.com/dreule28/Week_4/paas-api/internal/config"
	"github.com/dreule28/Week_4/paas-api/internal/kube"
	"github.com/dreule28/Week_4/paas-api/internal/model"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)


type InstanceService struct {
	cfg	config.Config
	k	*kube.Client
}

func deriveStatus(c *cnpgv1.Cluster) string {
	if c.Status.Phase != "" {
		return string(c.Status.Phase)
	}

	if c.Status.ReadyInstances == 0 {
		return "Creating"
	}

	return "Unknown"
}

func (s *InstanceService) List(ctx context.Context) ([]model.Instance, error) {
	var clusters cnpgv1.ClusterList
	if err := s.k.K8sClient.List(ctx, &clusters, client.InNamespace(s.cfg.Namespace),);
	 err != nil {
		return nil, err
	}
	out := make([]model.Instance, 0, len(clusters.Items))
	for _, c := range clusters.Items {
		out = append(out, model.Instance{
			ID:		c.Name,
			Status:	deriveStatus(&c),
		})
	}
	return out, nil
}

func NewInstanceService(cfg config.Config, k *kube.Client) *InstanceService {
	return &InstanceService{cfg: cfg, k: k}
}