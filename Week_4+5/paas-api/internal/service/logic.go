package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/dreule28/Week_4/paas-api/internal/config"
	"github.com/dreule28/Week_4/paas-api/internal/kube"
	"github.com/dreule28/Week_4/paas-api/internal/model"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InstanceService struct {
	cfg	config.Config
	k	*kube.Client
}

func NewInstanceService(cfg config.Config, k *kube.Client) *InstanceService {
	return &InstanceService{cfg: cfg, k: k}
}

func getStatus(c *cnpgv1.Cluster) string {
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

func (s *InstanceService) ListDatabases(ctx context.Context) ([]model.Instance, error) {
	var clusters cnpgv1.ClusterList
	if err := s.k.K8sClient.List(ctx, &clusters, client.InNamespace(s.cfg.Namespace)); err != nil {
		return nil, err
	}
	out := make([]model.Instance, 0, len(clusters.Items))
	for _, c := range clusters.Items {
		out = append(out, model.Instance{
			ID:     c.Name,
			Status: getStatus(&c),
		})
	}
	return out, nil
}

func (s *InstanceService) GetDatabase(ctx context.Context, id string) (model.InstanceDetails, error) {
	//get CR
	var cluster cnpgv1.Cluster
	if err := s.k.K8sClient.Get(ctx, types.NamespacedName{
		Namespace: s.cfg.Namespace,
		Name:      id,
	}, &cluster); err != nil {
		if apierrors.IsNotFound(err) {
			return model.InstanceDetails{}, fmt.Errorf("instance %q: %w", id, ErrNotFound)
		}
		return model.InstanceDetails{}, err
	}
	status := getStatus(&cluster)
	out := model.InstanceDetails{
		ID:     id,
		Status: status,
	}

	//if CR not ready, return without connection info
	if status != "ready" {
		return out, nil
	}

	secretName := fmt.Sprintf("%s-app", id)

	var sec corev1.Secret
	if err := s.k.K8sClient.Get(ctx, types.NamespacedName{
		Namespace: s.cfg.Namespace,
		Name:      secretName,
	}, &sec); err != nil {
		return out, nil
	}
	password := string(sec.Data["password"])

	host := fmt.Sprintf("%s-rw.%s.svc.cluster.local", id, s.cfg.Namespace)
	port := 5432
	endpoint := fmt.Sprintf("postgres://app@%s:%d/app", host, port)

	out.Connection = &model.ConnectionInfo{
		Host:     host,
		Port:     port,
		Database: "app",
		User:     "app",
		Password: password,
		Endpoint: endpoint,
	}
	return out, nil
}

func (s *InstanceService) CreateDatabase(ctx context.Context, req model.CreateInstanceRequest) (model.Instance, error) {
	cluster := &cnpgv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.ID,
			Namespace: s.cfg.Namespace,
			Labels: map[string]string{
				"paas.stackit.dev/managed": "true",
			},
		},
		Spec: cnpgv1.ClusterSpec{
			Instances: req.Instances,
			StorageConfiguration: cnpgv1.StorageConfiguration{
				Size: fmt.Sprintf("%dGi", req.StorageGi),
			},
		},
	}

	if err := s.k.K8sClient.Create(ctx, cluster); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return model.Instance{}, fmt.Errorf("instance %q: %w", req.ID, ErrAlreadyExists)
		}
		return model.Instance{}, err
	}
	return model.Instance{
		ID:     req.ID,
		Status: "creating",
	}, nil
}

func (s *InstanceService) DeleteDatabase(ctx context.Context, id string) error {
	cluster := &cnpgv1.Cluster{}
	cluster.Name = id
	cluster.Namespace = s.cfg.Namespace

	if err := s.k.K8sClient.Delete(ctx, cluster); err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("instance %q: %w", id, ErrNotFound)
		}
		return err
	}
	return nil
}

