package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/dreule28/Week_4/paas-api/internal/config"
	"github.com/dreule28/Week_4/paas-api/internal/kube"
	"github.com/dreule28/Week_4/paas-api/internal/model"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)


type InstanceService struct {
	cfg	config.Config
	k	*kube.Client
}

func deriveStatus(c *cnpgv1.Cluster) string {
	phase := strings.ToLower(strings.TrimSpace(c.Status.Phase))
	switch {
	case phase == "":
			return "creating"
	case strings.Contains(phase, "healthy") || strings.Contains(phase, "ready"):
		return "ready"
	case strings.Contains(phase, "error") || strings.Contains(phase, "fail"):
		return "error"
	default:
		return "creating"
	}
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

func (s *InstanceService) Get(ctx context.Context, id string) (model.InstanceDetails, error) {
	//get CR
	var cluster cnpgv1.Cluster
	if err := s.k.K8sClient.Get(ctx, types.NamespacedName{
		Namespace:	s.cfg.Namespace,
		Name:		id ,
	}, &cluster); err != nil {
		if apierrors.IsNotFound(err) {
			return model.InstanceDetails{}, fmt.Errorf("instance %q not found", id)
		}
		return model.InstanceDetails{}, err
	}
	status := deriveStatus(&cluster)
	out := model.InstanceDetails{
		ID:		id,
		Status:	status,
	}

	//if CR not ready, return without connection info
	if status != "ready" {
		return out, nil
	}

	host := fmt.Sprintf("%s-rw.%s.svc.cluster.local", id, s.cfg.Namespace)
	port := 5432
	endpoint := fmt.Sprintf("postgres://app@%s:%d/app", host, port)

	out.Connection = &model.ConnectionInfo{
		Host:		host,
		Port:		port,
		Database:	"app",
		User:		"app",
		Password:	"womp womp",
		Endpoint:	endpoint,
	}
	return out, nil
}

func (s *InstanceService) Create(ctx context.Context, req model.CreateInstanceRequest) (model.Instance, error) {
	cluster := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "postgresql.cnpg.io/v1",
			"kind":       "Cluster",
			"metadata": map[string]any{
				"name":      req.ID,
				"namespace": s.cfg.Namespace,
				"labels": map[string]any{
					"paas.stackit.dev/managed": "true",
				},
			},
			"spec": map[string]any{
				"instances": req.Instances,
				"storage": map[string]any{
					"size": fmt.Sprintf("%dGi", req.StorageGi),
				},
			},
		},
	}

	if err := s.k.K8sClient.Create(ctx, cluster); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return model.Instance{}, fmt.Errorf("instance %q already exists", req.ID)
		}
		return	model.Instance{}, err
	}
	return model.Instance{
		ID:		req.ID,
		Status:	"creating",
	}, nil
}

func (s *InstanceService) Delete(ctx context.Context, id string) error {
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion("postgresql.cnpg.io/v1")
	obj.SetKind("Cluster")
	obj.SetName(id)
	obj.SetNamespace(s.cfg.Namespace)

	if err := s.k.K8sClient.Delete(ctx, obj); err != nil {
		if apierrors.IsNotFound(err) {
			return	fmt.Errorf("instance %q not found", id)
		}
		return err
	}
	return nil
}

func NewInstanceService(cfg config.Config, k *kube.Client) *InstanceService {
	return &InstanceService{cfg: cfg, k: k}
}