package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TaskPriority defines the importance level of a Task.
// +kubebuilder:validation:Enum=Low;Medium;High;Critical
type TaskPriority string

const (
	PriorityLow      TaskPriority = "Low"
	PriorityMedium   TaskPriority = "Medium"
	PriorityHigh     TaskPriority = "High"
	PriorityCritical TaskPriority = "Critical"
)

// TaskPhase defines the current lifecycle state of a Task.
// +kubebuilder:validation:Enum=Pending;InProgress;Completed
type TaskPhase string

const (
	PhasePending    TaskPhase = "Pending"
	PhaseInProgress TaskPhase = "InProgress"
	PhaseCompleted  TaskPhase = "Completed"
)

// ProjectRef contains the reference details to the parent Project resource.
type ProjectRef struct {
	// Name of the referenced Project resource.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Namespace of the referenced Project resource. Defaults to current namespace if omitted.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// TaskSpec defines the desired state of a Task.
type TaskSpec struct {
	// Title of the task.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Title string `json:"title"`

	// Detailed description of the task requirements.
	// +optional
	Description string `json:"description,omitempty"`

	// Priority level of the task.
	// +kubebuilder:validation:Enum=Low;Medium;High;Critical
	// +kubebuilder:default=Medium
	// +optional
	Priority TaskPriority `json:"priority,omitempty"`

	// DesiredPhase requested by user/client.
	// +kubebuilder:validation:Enum=Pending;InProgress;Completed
	// +kubebuilder:default=Pending
	// +optional
	DesiredPhase TaskPhase `json:"desiredPhase,omitempty"`

	// ProjectRef links this task to a parent Project CR.
	// +kubebuilder:validation:Required
	ProjectRef ProjectRef `json:"projectRef"`

	// DueDate in ISO-8601 format (YYYY-MM-DD).
	// +optional
	DueDate string `json:"dueDate,omitempty"`
}

// TaskStatus defines the observed state of a Task reconciled by the Controller.
type TaskStatus struct {
	// Phase observed and reconciled by the controller.
	// +optional
	Phase TaskPhase `json:"phase,omitempty"`

	// ObservedGeneration tracks the generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// SyncTimestamp records the last time state was synchronized.
	// +optional
	SyncTimestamp *metav1.Time `json:"syncTimestamp,omitempty"`

	// Conditions represents the latest available observations of the task state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=tk
// +kubebuilder:printcolumn:name="PRIORITY",type="string",JSONPath=".spec.priority"
// +kubebuilder:printcolumn:name="DESIRED",type="string",JSONPath=".spec.desiredPhase"
// +kubebuilder:printcolumn:name="PHASE",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="PROJECT",type="string",JSONPath=".spec.projectRef.name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// Task is the Schema for the tasks API.
type Task struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskSpec   `json:"spec,omitempty"`
	Status TaskStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TaskList contains a list of Task resources.
type TaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Task `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Task{}, &TaskList{})
}
