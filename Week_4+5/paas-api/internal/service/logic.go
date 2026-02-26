package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/dreule28/Week_4/paas-api/internal/config"
	"github.com/dreule28/Week_4/paas-api/internal/kube"
	"github.com/dreule28/Week_4/paas-api/internal/model"
	logstore "github.com/dreule28/Week_4/paas-api/internal/store"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InstanceService struct {
	cfg       config.Config
	k         *kube.Client
	logs      logstore.LogStore
	statusMu  sync.Mutex
	statusMap map[string]string
}

func NewInstanceService(cfg config.Config, k *kube.Client) *InstanceService {
	ls, err := logstore.NewFileLogStore(cfg.LogsPath)
	if err != nil {
		log.Printf("log store disabled: %v", err)
	}
	return &InstanceService{
		cfg:       cfg,
		k:         k,
		logs:      ls,
		statusMap: map[string]string{},
	}
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
		status := getStatus(&c)
		out = append(out, model.Instance{
			ID:     c.Name,
			Status: status,
		})
		s.recordStatusChange(c.Name, status)
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
	s.recordStatusChange(id, status)
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
	s.recordServiceLog(
		req.ID,
		"instance.create.accepted",
		fmt.Sprintf("Service accepted create request and submitted cluster resource (replicas=%d, storageGi=%d)", req.Instances, req.StorageGi),
	)
	s.recordStatusChange(req.ID, "creating")
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
	s.recordServiceLog(id, "instance.delete.accepted", "Service accepted delete request and submitted cluster deletion")
	s.statusMu.Lock()
	delete(s.statusMap, id)
	s.statusMu.Unlock()
	return nil
}

func (s *InstanceService) ListInstanceLogs(_ context.Context, q model.LogQuery) ([]model.LogEntry, error) {
	if s.logs == nil {
		return []model.LogEntry{}, nil
	}
	return s.logs.List(q)
}

func (s *InstanceService) RecordAuditLog(_ context.Context, instanceID, user, action, message string) error {
	if user == "" {
		user = "unknown"
	}
	if message == "" {
		message = action
	}
	_, err := s.appendLog(model.LogEntry{
		InstanceID: instanceID,
		Type:       model.LogTypeAudit,
		Action:     action,
		Message:    message,
		User:       user,
	})
	return err
}

func (s *InstanceService) recordStatusChange(instanceID, status string) {
	if instanceID == "" || status == "" {
		return
	}

	s.statusMu.Lock()
	prev, ok := s.statusMap[instanceID]
	if ok && prev == status {
		s.statusMu.Unlock()
		return
	}

	msg := fmt.Sprintf("Platform observed status=%s", status)
	if ok {
		msg = fmt.Sprintf("Platform observed status transition %s -> %s", prev, status)
	}
	if err := s.recordServiceLog(instanceID, "status.changed", msg); err == nil {
		s.statusMap[instanceID] = status
	}
	s.statusMu.Unlock()
}

func (s *InstanceService) recordServiceLog(instanceID, action, message string) error {
	_, err := s.appendLog(model.LogEntry{
		InstanceID: instanceID,
		Type:       model.LogTypeService,
		Action:     action,
		Message:    message,
	})
	return err
}

func (s *InstanceService) appendLog(entry model.LogEntry) (model.LogEntry, error) {
	if s.logs == nil {
		return model.LogEntry{}, nil
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	return s.logs.Append(entry)
}
