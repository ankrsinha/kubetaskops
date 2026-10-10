package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProjectHealth indicates the overall status/progress of a project.
// +kubebuilder:validation:Enum=OnTrack;AtRisk;Overdue
type ProjectHealth string

const (
	HealthOnTrack ProjectHealth = "OnTrack"
	HealthAtRisk  ProjectHealth = "AtRisk"
	HealthOverdue ProjectHealth = "Overdue"
)

// ProjectSpec defines the desired state of a Project.
type ProjectSpec struct {
	// Name of the project.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Detailed description of the project goals.
	// +optional
	Description string `json:"description,omitempty"`

	// Lead assigned to manage the project.
	// +optional
	Lead string `json:"lead,omitempty"`

	// TargetDate in ISO-8601 format (YYYY-MM-DD).
	// +optional
	TargetDate string `json:"targetDate,omitempty"`
}

// ProjectStatus defines the observed state of a Project reconciled by the Controller.
type ProjectStatus struct {
	// TaskCount represents the total number of Tasks linked to this project.
	// +optional
	// +kubebuilder:default=0
	TaskCount int32 `json:"taskCount"`

	// CompletedCount represents the number of completed Tasks linked to this project.
	// +optional
	// +kubebuilder:default=0
	CompletedCount int32 `json:"completedCount"`

	// Health status calculated based on task deadlines and progress.
	// +optional
	// +kubebuilder:validation:Enum=OnTrack;AtRisk;Overdue
	// +kubebuilder:default=OnTrack
	Health ProjectHealth `json:"health,omitempty"`

	// ObservedGeneration tracks the generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions represents the latest available observations of the project state.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=prj
// +kubebuilder:printcolumn:name="LEAD",type="string",JSONPath=".spec.lead"
// +kubebuilder:printcolumn:name="HEALTH",type="string",JSONPath=".status.health"
// +kubebuilder:printcolumn:name="TASKS",type="integer",JSONPath=".status.taskCount"
// +kubebuilder:printcolumn:name="COMPLETED",type="integer",JSONPath=".status.completedCount"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// Project is the Schema for the projects API.
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectSpec   `json:"spec,omitempty"`
	Status ProjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProjectList contains a list of Project resources.
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Project{}, &ProjectList{})
}
