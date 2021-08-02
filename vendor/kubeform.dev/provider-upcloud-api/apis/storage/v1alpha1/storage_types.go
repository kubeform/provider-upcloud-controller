/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by Kubeform. DO NOT EDIT.

package v1alpha1

import (
	base "kubeform.dev/apimachinery/api/v1alpha1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`

type Storage struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              StorageSpec   `json:"spec,omitempty"`
	Status            StorageStatus `json:"status,omitempty"`
}

type StorageSpecBackupRule struct {
	// The weekday when the backup is created
	Interval *string `json:"interval" tf:"interval"`
	// The number of days before a backup is automatically deleted
	Retention *int64 `json:"retention" tf:"retention"`
	// The time of day when the backup is created
	Time *string `json:"time" tf:"time"`
}

type StorageSpecClone struct {
	// The unique identifier of the storage/template to clone
	ID *string `json:"ID" tf:"id"`
}

type StorageSpecImport struct {
	// sha256 sum of the imported data
	// +optional
	Sha256sum *string `json:"sha256sum,omitempty" tf:"sha256sum"`
	// The mode of the import task. One of `http_import` or `direct_upload`.
	Source *string `json:"source" tf:"source"`
	// For `direct_upload`; an optional hash of the file to upload.
	// +optional
	SourceHash *string `json:"sourceHash,omitempty" tf:"source_hash"`
	// The location of the file to import. For `http_import` an accessible URL for `direct_upload` a local file.
	SourceLocation *string `json:"sourceLocation" tf:"source_location"`
	// Number of bytes imported
	// +optional
	WrittenBytes *int64 `json:"writtenBytes,omitempty" tf:"written_bytes"`
}

type StorageSpec struct {
	State *StorageSpecResource `json:"state,omitempty" tf:"-"`

	Resource StorageSpecResource `json:"resource" tf:"resource"`

	UpdatePolicy base.UpdatePolicy `json:"updatePolicy,omitempty" tf:"-"`

	TerminationPolicy base.TerminationPolicy `json:"terminationPolicy,omitempty" tf:"-"`

	ProviderRef core.LocalObjectReference `json:"providerRef" tf:"-"`
}

type StorageSpecResource struct {
	ID string `json:"id,omitempty" tf:"id,omitempty"`

	// The criteria to backup the storage
	// +optional
	BackupRule *StorageSpecBackupRule `json:"backupRule,omitempty" tf:"backup_rule"`
	// Block defining another storage/template to clone to storage
	// +optional
	Clone *StorageSpecClone `json:"clone,omitempty" tf:"clone"`
	// Block defining external data to import to storage
	// +optional
	Import *StorageSpecImport `json:"import,omitempty" tf:"import"`
	// The size of the storage in gigabytes
	Size *int64 `json:"size" tf:"size"`
	// The storage tier to use
	// +optional
	Tier *string `json:"tier,omitempty" tf:"tier"`
	// A short, informative description
	Title *string `json:"title" tf:"title"`
	// The zone in which the storage will be created
	Zone *string `json:"zone" tf:"zone"`
}

type StorageStatus struct {
	// Resource generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// +optional
	Phase status.Status `json:"phase,omitempty"`
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// StorageList is a list of Storages
type StorageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Storage CRD objects
	Items []Storage `json:"items,omitempty"`
}