package service

import (
	"context"
	"fmt"
	"time"

	"velero-api-server/internal/model"

	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type VeleroService struct {
	namespace string
}

func NewVeleroService(namespace string) *VeleroService {
	return &VeleroService{namespace: namespace}
}

// ========== Backup ==========

func (s *VeleroService) CreateBackup(ctx context.Context, c client.Client, req model.CreateBackupRequest) (*model.BackupResponse, error) {
	backup := &velerov1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: s.namespace,
		},
		Spec: velerov1.BackupSpec{
			IncludedNamespaces: req.IncludedNamespaces,
			ExcludedNamespaces: req.ExcludedNamespaces,
			IncludedResources:  req.IncludedResources,
			ExcludedResources:  req.ExcludedResources,
			StorageLocation:    req.StorageLocation,
			SnapshotVolumes:    req.SnapshotVolumes,
		},
	}

	if req.TTL != "" {
		duration, err := time.ParseDuration(req.TTL)
		if err != nil {
			return nil, fmt.Errorf("invalid TTL format: %v", err)
		}
		backup.Spec.TTL = metav1.Duration{Duration: duration}
	}

	if len(req.LabelSelector) > 0 {
		backup.Spec.LabelSelector = &metav1.LabelSelector{
			MatchLabels: req.LabelSelector,
		}
	}

	if err := c.Create(ctx, backup); err != nil {
		return nil, fmt.Errorf("failed to create backup: %v", err)
	}

	return backupToResponse(backup), nil
}

func (s *VeleroService) ListBackups(ctx context.Context, c client.Client) ([]model.BackupResponse, error) {
	list := &velerov1.BackupList{}
	if err := c.List(ctx, list, client.InNamespace(s.namespace)); err != nil {
		return nil, fmt.Errorf("failed to list backups: %v", err)
	}

	result := make([]model.BackupResponse, 0, len(list.Items))
	for i := range list.Items {
		result = append(result, *backupToResponse(&list.Items[i]))
	}
	return result, nil
}

func (s *VeleroService) GetBackup(ctx context.Context, c client.Client, name string) (*model.BackupResponse, error) {
	backup := &velerov1.Backup{}
	key := client.ObjectKey{Namespace: s.namespace, Name: name}
	if err := c.Get(ctx, key, backup); err != nil {
		return nil, fmt.Errorf("failed to get backup %s: %v", name, err)
	}
	return backupToResponse(backup), nil
}

func (s *VeleroService) DeleteBackup(ctx context.Context, c client.Client, name string) error {
	deleteReq := &velerov1.DeleteBackupRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-delete-" + time.Now().Format("20060102150405"),
			Namespace: s.namespace,
			Labels: map[string]string{
				velerov1.BackupNameLabel: name,
			},
		},
		Spec: velerov1.DeleteBackupRequestSpec{
			BackupName: name,
		},
	}
	if err := c.Create(ctx, deleteReq); err != nil {
		return fmt.Errorf("failed to delete backup %s: %v", name, err)
	}
	return nil
}

// ========== Restore ==========

func (s *VeleroService) CreateRestore(ctx context.Context, c client.Client, req model.CreateRestoreRequest) (*model.RestoreResponse, error) {
	restore := &velerov1.Restore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: s.namespace,
		},
		Spec: velerov1.RestoreSpec{
			BackupName:         req.BackupName,
			IncludedNamespaces: req.IncludedNamespaces,
			ExcludedNamespaces: req.ExcludedNamespaces,
			IncludedResources:  req.IncludedResources,
			ExcludedResources:  req.ExcludedResources,
			RestorePVs:         req.RestorePVs,
		},
	}

	if len(req.LabelSelector) > 0 {
		restore.Spec.LabelSelector = &metav1.LabelSelector{
			MatchLabels: req.LabelSelector,
		}
	}

	if err := c.Create(ctx, restore); err != nil {
		return nil, fmt.Errorf("failed to create restore: %v", err)
	}
	return restoreToResponse(restore), nil
}

func (s *VeleroService) ListRestores(ctx context.Context, c client.Client, backupName string) ([]model.RestoreResponse, error) {
	list := &velerov1.RestoreList{}
	opts := []client.ListOption{client.InNamespace(s.namespace)}

	if backupName != "" {
		opts = append(opts, client.MatchingLabels{
			velerov1.BackupNameLabel: backupName,
		})
	}

	if err := c.List(ctx, list, opts...); err != nil {
		return nil, fmt.Errorf("failed to list restores: %v", err)
	}

	result := make([]model.RestoreResponse, 0, len(list.Items))
	for i := range list.Items {
		result = append(result, *restoreToResponse(&list.Items[i]))
	}
	return result, nil
}

func (s *VeleroService) GetRestore(ctx context.Context, c client.Client, name string) (*model.RestoreResponse, error) {
	restore := &velerov1.Restore{}
	key := client.ObjectKey{Namespace: s.namespace, Name: name}
	if err := c.Get(ctx, key, restore); err != nil {
		return nil, fmt.Errorf("failed to get restore %s: %v", name, err)
	}
	return restoreToResponse(restore), nil
}

// ========== Schedule ==========

func (s *VeleroService) CreateSchedule(ctx context.Context, c client.Client, req model.CreateScheduleRequest) (*model.ScheduleResponse, error) {
	schedule := &velerov1.Schedule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: s.namespace,
		},
		Spec: velerov1.ScheduleSpec{
			Schedule: req.Schedule,
			Template: velerov1.BackupSpec{
				IncludedNamespaces: req.IncludedNamespaces,
				ExcludedNamespaces: req.ExcludedNamespaces,
				IncludedResources:  req.IncludedResources,
				ExcludedResources:  req.ExcludedResources,
				StorageLocation:    req.StorageLocation,
			},
		},
	}

	if req.TTL != "" {
		duration, err := time.ParseDuration(req.TTL)
		if err != nil {
			return nil, fmt.Errorf("invalid TTL format: %v", err)
		}
		schedule.Spec.Template.TTL = metav1.Duration{Duration: duration}
	}

	if len(req.LabelSelector) > 0 {
		schedule.Spec.Template.LabelSelector = &metav1.LabelSelector{
			MatchLabels: req.LabelSelector,
		}
	}

	if err := c.Create(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %v", err)
	}
	return scheduleToResponse(schedule), nil
}

func (s *VeleroService) ListSchedules(ctx context.Context, c client.Client) ([]model.ScheduleResponse, error) {
	list := &velerov1.ScheduleList{}
	if err := c.List(ctx, list, client.InNamespace(s.namespace)); err != nil {
		return nil, fmt.Errorf("failed to list schedules: %v", err)
	}

	result := make([]model.ScheduleResponse, 0, len(list.Items))
	for i := range list.Items {
		result = append(result, *scheduleToResponse(&list.Items[i]))
	}
	return result, nil
}

func (s *VeleroService) GetSchedule(ctx context.Context, c client.Client, name string) (*model.ScheduleResponse, error) {
	schedule := &velerov1.Schedule{}
	key := client.ObjectKey{Namespace: s.namespace, Name: name}
	if err := c.Get(ctx, key, schedule); err != nil {
		return nil, fmt.Errorf("failed to get schedule %s: %v", name, err)
	}
	return scheduleToResponse(schedule), nil
}

func (s *VeleroService) DeleteSchedule(ctx context.Context, c client.Client, name string) error {
	schedule := &velerov1.Schedule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: s.namespace,
		},
	}
	if err := c.Delete(ctx, schedule); err != nil {
		return fmt.Errorf("failed to delete schedule %s: %v", name, err)
	}
	return nil
}

// ========== BackupStorageLocation ==========

func (s *VeleroService) ListBSLs(ctx context.Context, c client.Client) ([]model.BSLResponse, error) {
	list := &velerov1.BackupStorageLocationList{}
	if err := c.List(ctx, list, client.InNamespace(s.namespace)); err != nil {
		return nil, fmt.Errorf("failed to list backup storage locations: %v", err)
	}

	result := make([]model.BSLResponse, 0, len(list.Items))
	for i := range list.Items {
		result = append(result, *bslToResponse(&list.Items[i]))
	}
	return result, nil
}

func (s *VeleroService) GetBSL(ctx context.Context, c client.Client, name string) (*model.BSLResponse, error) {
	bsl := &velerov1.BackupStorageLocation{}
	key := client.ObjectKey{Namespace: s.namespace, Name: name}
	if err := c.Get(ctx, key, bsl); err != nil {
		return nil, fmt.Errorf("failed to get backup storage location %s: %v", name, err)
	}
	return bslToResponse(bsl), nil
}

// ========== Converters ==========

func backupToResponse(b *velerov1.Backup) *model.BackupResponse {
	ttl := ""
	if b.Spec.TTL.Duration != 0 {
		ttl = b.Spec.TTL.Duration.String()
	}

	resp := &model.BackupResponse{
		Name:               b.Name,
		Namespace:          b.Namespace,
		Phase:              string(b.Status.Phase),
		IncludedNamespaces: b.Spec.IncludedNamespaces,
		ExcludedNamespaces: b.Spec.ExcludedNamespaces,
		StorageLocation:    b.Spec.StorageLocation,
		TTL:                ttl,
		StartTimestamp:     b.Status.StartTimestamp,
		CompletionTime:     b.Status.CompletionTimestamp,
		Errors:             b.Status.Errors,
		Warnings:           b.Status.Warnings,
		Labels:             b.Labels,
		CreatedAt:          b.CreationTimestamp,
	}
	return resp
}

func restoreToResponse(r *velerov1.Restore) *model.RestoreResponse {
	return &model.RestoreResponse{
		Name:           r.Name,
		Namespace:      r.Namespace,
		BackupName:     r.Spec.BackupName,
		Phase:          string(r.Status.Phase),
		StartTimestamp: r.Status.StartTimestamp,
		CompletionTime: r.Status.CompletionTimestamp,
		Errors:         r.Status.Errors,
		Warnings:       r.Status.Warnings,
		CreatedAt:      r.CreationTimestamp,
	}
}

func scheduleToResponse(s *velerov1.Schedule) *model.ScheduleResponse {
	ttl := ""
	if s.Spec.Template.TTL.Duration != 0 {
		ttl = s.Spec.Template.TTL.Duration.String()
	}

	return &model.ScheduleResponse{
		Name:            s.Name,
		Namespace:       s.Namespace,
		Schedule:        s.Spec.Schedule,
		Phase:           string(s.Status.Phase),
		LastBackup:      s.Status.LastBackup,
		StorageLocation: s.Spec.Template.StorageLocation,
		TTL:             ttl,
		CreatedAt:       s.CreationTimestamp,
	}
}

func bslToResponse(b *velerov1.BackupStorageLocation) *model.BSLResponse {
	resp := &model.BSLResponse{
		Name:       b.Name,
		Namespace:  b.Namespace,
		Provider:   b.Spec.Provider,
		Phase:      string(b.Status.Phase),
		LastSynced: b.Status.LastSyncedTime,
		CreatedAt:  b.CreationTimestamp,
	}
	if b.Spec.ObjectStorage != nil {
		resp.Bucket = b.Spec.ObjectStorage.Bucket
		resp.Prefix = b.Spec.ObjectStorage.Prefix
	}
	if b.Spec.Default {
		resp.IsDefault = true
	}
	return resp
}

// suppress unused import
var _ = labels.Everything
