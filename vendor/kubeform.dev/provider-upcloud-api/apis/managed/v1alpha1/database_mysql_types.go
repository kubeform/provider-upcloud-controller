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

type DatabaseMysql struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DatabaseMysqlSpec   `json:"spec,omitempty"`
	Status            DatabaseMysqlStatus `json:"status,omitempty"`
}

type DatabaseMysqlSpecComponents struct {
	// Type of the component
	// +optional
	Component *string `json:"component,omitempty" tf:"component"`
	// Hostname of the component
	// +optional
	Host *string `json:"host,omitempty" tf:"host"`
	// Port number of the component
	// +optional
	Port *int64 `json:"port,omitempty" tf:"port"`
	// Component network route type
	// +optional
	Route *string `json:"route,omitempty" tf:"route"`
	// Usage of the component
	// +optional
	Usage *string `json:"usage,omitempty" tf:"usage"`
}

type DatabaseMysqlSpecNodeStates struct {
	// Name plus a node iteration
	// +optional
	Name *string `json:"name,omitempty" tf:"name"`
	// Role of the node
	// +optional
	Role *string `json:"role,omitempty" tf:"role"`
	// State of the node
	// +optional
	State *string `json:"state,omitempty" tf:"state"`
}

type DatabaseMysqlSpecPropertiesMigration struct {
	// Database name for bootstrapping the initial connection
	// +optional
	Dbname *string `json:"dbname,omitempty" tf:"dbname"`
	// Hostname or IP address of the server where to migrate data from
	// +optional
	Host *string `json:"host,omitempty" tf:"host"`
	// Comma-separated list of databases, which should be ignored during migration (supported by MySQL only at the moment)
	// +optional
	IgnoreDbs *string `json:"ignoreDbs,omitempty" tf:"ignore_dbs"`
	// Password for authentication with the server where to migrate data from
	// +optional
	Password *string `json:"-" sensitive:"true" tf:"password"`
	// Port number of the server where to migrate data from
	// +optional
	Port *int64 `json:"port,omitempty" tf:"port"`
	// The server where to migrate data from is secured with SSL
	// +optional
	Ssl *bool `json:"ssl,omitempty" tf:"ssl"`
	// User name for authentication with the server where to migrate data from
	// +optional
	Username *string `json:"username,omitempty" tf:"username"`
}

type DatabaseMysqlSpecProperties struct {
	// Custom password for admin user. Defaults to random string. This must be set only when a new service is being created.
	// +optional
	AdminPassword *string `json:"-" sensitive:"true" tf:"admin_password"`
	// Custom username for admin user. This must be set only when a new service is being created.
	// +optional
	AdminUsername *string `json:"adminUsername,omitempty" tf:"admin_username"`
	// Automatic utility network IP Filter
	// +optional
	AutomaticUtilityNetworkIPFilter *bool `json:"automaticUtilityNetworkIPFilter,omitempty" tf:"automatic_utility_network_ip_filter"`
	// The hour of day (in UTC) when backup for the service is started. New backup is only started if previous backup has already completed.
	// +optional
	BackupHour *int64 `json:"backupHour,omitempty" tf:"backup_hour"`
	// The minute of an hour when backup for the service is started. New backup is only started if previous backup has already completed.
	// +optional
	BackupMinute *int64 `json:"backupMinute,omitempty" tf:"backup_minute"`
	// The minimum amount of time in seconds to keep binlog entries before deletion. This may be extended for services that require binlog entries for longer than the default for example if using the MySQL Debezium Kafka connector.
	// +optional
	BinlogRetentionPeriod *int64 `json:"binlogRetentionPeriod,omitempty" tf:"binlog_retention_period"`
	// connect_timeout
	// +optional
	ConnectTimeout *int64 `json:"connectTimeout,omitempty" tf:"connect_timeout"`
	// default_time_zone
	// +optional
	DefaultTimeZone *string `json:"defaultTimeZone,omitempty" tf:"default_time_zone"`
	// group_concat_max_len
	// +optional
	GroupConcatMaxLen *int64 `json:"groupConcatMaxLen,omitempty" tf:"group_concat_max_len"`
	// information_schema_stats_expiry
	// +optional
	InformationSchemaStatsExpiry *int64 `json:"informationSchemaStatsExpiry,omitempty" tf:"information_schema_stats_expiry"`
	// innodb_ft_min_token_size
	// +optional
	InnodbFtMinTokenSize *int64 `json:"innodbFtMinTokenSize,omitempty" tf:"innodb_ft_min_token_size"`
	// innodb_ft_server_stopword_table
	// +optional
	InnodbFtServerStopwordTable *string `json:"innodbFtServerStopwordTable,omitempty" tf:"innodb_ft_server_stopword_table"`
	// innodb_lock_wait_timeout
	// +optional
	InnodbLockWaitTimeout *int64 `json:"innodbLockWaitTimeout,omitempty" tf:"innodb_lock_wait_timeout"`
	// innodb_log_buffer_size
	// +optional
	InnodbLogBufferSize *int64 `json:"innodbLogBufferSize,omitempty" tf:"innodb_log_buffer_size"`
	// innodb_online_alter_log_max_size
	// +optional
	InnodbOnlineAlterLogMaxSize *int64 `json:"innodbOnlineAlterLogMaxSize,omitempty" tf:"innodb_online_alter_log_max_size"`
	// innodb_print_all_deadlocks
	// +optional
	InnodbPrintAllDeadlocks *bool `json:"innodbPrintAllDeadlocks,omitempty" tf:"innodb_print_all_deadlocks"`
	// innodb_rollback_on_timeout
	// +optional
	InnodbRollbackOnTimeout *bool `json:"innodbRollbackOnTimeout,omitempty" tf:"innodb_rollback_on_timeout"`
	// interactive_timeout
	// +optional
	InteractiveTimeout *int64 `json:"interactiveTimeout,omitempty" tf:"interactive_timeout"`
	// internal_tmp_mem_storage_engine
	// +optional
	InternalTmpMemStorageEngine *string `json:"internalTmpMemStorageEngine,omitempty" tf:"internal_tmp_mem_storage_engine"`
	// IP filter
	// +optional
	// +kubebuilder:validation:MaxItems=1024
	IpFilter []string `json:"ipFilter,omitempty" tf:"ip_filter"`
	// long_query_time
	// +optional
	LongQueryTime *float64 `json:"longQueryTime,omitempty" tf:"long_query_time"`
	// max_allowed_packet
	// +optional
	MaxAllowedPacket *int64 `json:"maxAllowedPacket,omitempty" tf:"max_allowed_packet"`
	// max_heap_table_size
	// +optional
	MaxHeapTableSize *int64 `json:"maxHeapTableSize,omitempty" tf:"max_heap_table_size"`
	// Migrate data from existing server
	// +optional
	Migration *DatabaseMysqlSpecPropertiesMigration `json:"migration,omitempty" tf:"migration"`
	// net_read_timeout
	// +optional
	NetReadTimeout *int64 `json:"netReadTimeout,omitempty" tf:"net_read_timeout"`
	// net_write_timeout
	// +optional
	NetWriteTimeout *int64 `json:"netWriteTimeout,omitempty" tf:"net_write_timeout"`
	// Public Access
	// +optional
	PublicAccess *bool `json:"publicAccess,omitempty" tf:"public_access"`
	// slow_query_log
	// +optional
	SlowQueryLog *bool `json:"slowQueryLog,omitempty" tf:"slow_query_log"`
	// sort_buffer_size
	// +optional
	SortBufferSize *int64 `json:"sortBufferSize,omitempty" tf:"sort_buffer_size"`
	// sql_mode
	// +optional
	SqlMode *string `json:"sqlMode,omitempty" tf:"sql_mode"`
	// sql_require_primary_key
	// +optional
	SqlRequirePrimaryKey *bool `json:"sqlRequirePrimaryKey,omitempty" tf:"sql_require_primary_key"`
	// tmp_table_size
	// +optional
	TmpTableSize *int64 `json:"tmpTableSize,omitempty" tf:"tmp_table_size"`
	// MySQL major version
	// +optional
	Version *string `json:"version,omitempty" tf:"version"`
	// wait_timeout
	// +optional
	WaitTimeout *int64 `json:"waitTimeout,omitempty" tf:"wait_timeout"`
}

type DatabaseMysqlSpec struct {
	State *DatabaseMysqlSpecResource `json:"state,omitempty" tf:"-"`

	Resource DatabaseMysqlSpecResource `json:"resource" tf:"resource"`

	UpdatePolicy base.UpdatePolicy `json:"updatePolicy,omitempty" tf:"-"`

	TerminationPolicy base.TerminationPolicy `json:"terminationPolicy,omitempty" tf:"-"`

	ProviderRef core.LocalObjectReference `json:"providerRef" tf:"-"`

	SecretRef *core.LocalObjectReference `json:"secretRef,omitempty" tf:"-"`

	BackendRef *core.LocalObjectReference `json:"backendRef,omitempty" tf:"-"`
}

type DatabaseMysqlSpecResource struct {
	ID string `json:"id,omitempty" tf:"id,omitempty"`

	// Service component information
	// +optional
	Components []DatabaseMysqlSpecComponents `json:"components,omitempty" tf:"components"`
	// Maintenance window day of week. Lower case weekday name (monday, tuesday, ...)
	// +optional
	MaintenanceWindowDow *string `json:"maintenanceWindowDow,omitempty" tf:"maintenance_window_dow"`
	// Maintenance window UTC time in hh:mm:ss format
	// +optional
	MaintenanceWindowTime *string `json:"maintenanceWindowTime,omitempty" tf:"maintenance_window_time"`
	// Name of the service. The name is used as a prefix for the logical hostname. Must be unique within an account
	Name *string `json:"name" tf:"name"`
	// Information about nodes providing the managed service
	// +optional
	NodeStates []DatabaseMysqlSpecNodeStates `json:"nodeStates,omitempty" tf:"node_states"`
	// Service plan to use. This determines how much resources the instance will have
	Plan *string `json:"plan" tf:"plan"`
	// The administrative power state of the service
	// +optional
	Powered *bool `json:"powered,omitempty" tf:"powered"`
	// Primary database name
	// +optional
	PrimaryDatabase *string `json:"primaryDatabase,omitempty" tf:"primary_database"`
	// Database Engine properties for MySQL
	// +optional
	Properties *DatabaseMysqlSpecProperties `json:"properties,omitempty" tf:"properties"`
	// Hostname to the service instance
	// +optional
	ServiceHost *string `json:"serviceHost,omitempty" tf:"service_host"`
	// Primary username's password to the service instance
	// +optional
	ServicePassword *string `json:"-" sensitive:"true" tf:"service_password"`
	// Port to the service instance
	// +optional
	ServicePort *string `json:"servicePort,omitempty" tf:"service_port"`
	// URI to the service instance
	// +optional
	ServiceURI *string `json:"-" sensitive:"true" tf:"service_uri"`
	// Primary username to the service instance
	// +optional
	ServiceUsername *string `json:"serviceUsername,omitempty" tf:"service_username"`
	// State of the service
	// +optional
	State *string `json:"state,omitempty" tf:"state"`
	// Title of a managed database instance
	// +optional
	Title *string `json:"title,omitempty" tf:"title"`
	// Type of the service
	// +optional
	Type *string `json:"type,omitempty" tf:"type"`
	// Zone where the instance resides
	Zone *string `json:"zone" tf:"zone"`
}

type DatabaseMysqlStatus struct {
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

// DatabaseMysqlList is a list of DatabaseMysqls
type DatabaseMysqlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of DatabaseMysql CRD objects
	Items []DatabaseMysql `json:"items,omitempty"`
}
