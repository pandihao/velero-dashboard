package model

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// ========== Backup ==========

type CreateBackupRequest struct {
	Name               string            `json:"name" binding:"required"`
	IncludedNamespaces []string          `json:"includedNamespaces,omitempty"`
	ExcludedNamespaces []string          `json:"excludedNamespaces,omitempty"`
	IncludedResources  []string          `json:"includedResources,omitempty"`
	ExcludedResources  []string          `json:"excludedResources,omitempty"`
	LabelSelector      map[string]string `json:"labelSelector,omitempty"`
	StorageLocation    string            `json:"storageLocation,omitempty"`
	TTL                string            `json:"ttl,omitempty"` // e.g. "720h"
	SnapshotVolumes    *bool             `json:"snapshotVolumes,omitempty"`
}

type BackupResponse struct {
	Name               string            `json:"name"`
	Namespace          string            `json:"namespace"`
	Phase              string            `json:"phase"`
	IncludedNamespaces []string          `json:"includedNamespaces,omitempty"`
	ExcludedNamespaces []string          `json:"excludedNamespaces,omitempty"`
	StorageLocation    string            `json:"storageLocation"`
	TTL                string            `json:"ttl"`
	StartTimestamp     *metav1.Time      `json:"startTimestamp,omitempty"`
	CompletionTime     *metav1.Time      `json:"completionTimestamp,omitempty"`
	Errors             int               `json:"errors"`
	Warnings           int               `json:"warnings"`
	Labels             map[string]string `json:"labels,omitempty"`
	CreatedAt          metav1.Time       `json:"createdAt"`
}

// ========== Restore ==========

type CreateRestoreRequest struct {
	Name               string            `json:"name" binding:"required"`
	BackupName         string            `json:"backupName" binding:"required"`
	IncludedNamespaces []string          `json:"includedNamespaces,omitempty"`
	ExcludedNamespaces []string          `json:"excludedNamespaces,omitempty"`
	IncludedResources  []string          `json:"includedResources,omitempty"`
	ExcludedResources  []string          `json:"excludedResources,omitempty"`
	LabelSelector      map[string]string `json:"labelSelector,omitempty"`
	RestorePVs         *bool             `json:"restorePVs,omitempty"`
}

type RestoreResponse struct {
	Name           string       `json:"name"`
	Namespace      string       `json:"namespace"`
	BackupName     string       `json:"backupName"`
	Phase          string       `json:"phase"`
	StartTimestamp *metav1.Time `json:"startTimestamp,omitempty"`
	CompletionTime *metav1.Time `json:"completionTimestamp,omitempty"`
	Errors         int          `json:"errors"`
	Warnings       int          `json:"warnings"`
	CreatedAt      metav1.Time  `json:"createdAt"`
}

// ========== Schedule ==========

type CreateScheduleRequest struct {
	Name               string            `json:"name" binding:"required"`
	Schedule           string            `json:"schedule" binding:"required"` // cron expression
	IncludedNamespaces []string          `json:"includedNamespaces,omitempty"`
	ExcludedNamespaces []string          `json:"excludedNamespaces,omitempty"`
	IncludedResources  []string          `json:"includedResources,omitempty"`
	ExcludedResources  []string          `json:"excludedResources,omitempty"`
	LabelSelector      map[string]string `json:"labelSelector,omitempty"`
	StorageLocation    string            `json:"storageLocation,omitempty"`
	TTL                string            `json:"ttl,omitempty"`
}

type ScheduleResponse struct {
	Name           string       `json:"name"`
	Namespace      string       `json:"namespace"`
	Schedule       string       `json:"schedule"`
	Phase          string       `json:"phase"`
	LastBackup     *metav1.Time `json:"lastBackup,omitempty"`
	StorageLocation string      `json:"storageLocation"`
	TTL            string       `json:"ttl"`
	CreatedAt      metav1.Time  `json:"createdAt"`
}

// ========== BackupStorageLocation ==========

type BSLResponse struct {
	Name        string      `json:"name"`
	Namespace   string      `json:"namespace"`
	Provider    string      `json:"provider"`
	Bucket      string      `json:"bucket"`
	Prefix      string      `json:"prefix,omitempty"`
	Phase       string      `json:"phase"`
	LastSynced  *metav1.Time `json:"lastSyncedTime,omitempty"`
	IsDefault   bool        `json:"isDefault"`
	CreatedAt   metav1.Time `json:"createdAt"`
}

// ========== Common ==========

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
