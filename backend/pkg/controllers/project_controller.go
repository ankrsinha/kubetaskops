package controllers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1alpha1 "github.com/ankrsinha/kubetaskops/backend/pkg/apis/v1alpha1"
)

// ProjectReconciler reconciles a Project object
type ProjectReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile aggregates child task counts and updates Project health
func (r *ProjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("project", req.NamespacedName)

	project := &v1alpha1.Project{}
	if err := r.Get(ctx, req.NamespacedName, project); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("Project resource deleted from cluster")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Project resource")
		return ctrl.Result{}, err
	}
	previousStatus := project.DeepCopy().Status

	// 1. List all Task resources. Task project references may cross namespaces.
	taskList := &v1alpha1.TaskList{}
	if err := r.List(ctx, taskList); err != nil {
		logger.Error(err, "Failed to list child Tasks for Project")
		return ctrl.Result{}, err
	}

	// 2. Aggregate statistics for tasks belonging to this project
	var totalTasks int32
	var completedTasks int32

	for _, task := range taskList.Items {
		projectNamespace := task.Spec.ProjectRef.Namespace
		if projectNamespace == "" {
			projectNamespace = task.Namespace
		}

		if task.Spec.ProjectRef.Name == project.Name && projectNamespace == project.Namespace {
			totalTasks++
			if task.Status.Phase == v1alpha1.PhaseCompleted {
				completedTasks++
			}
		}
	}

	// 3. Compute Project Health
	health := v1alpha1.HealthOnTrack
	if project.Spec.TargetDate != "" {
		targetTime, err := time.Parse("2006-01-02", project.Spec.TargetDate)
		if err == nil {
			now := time.Now()
			if now.After(targetTime) && (totalTasks == 0 || completedTasks < totalTasks) {
				health = v1alpha1.HealthOverdue
			} else if targetTime.Sub(now) < 7*24*time.Hour && totalTasks > 0 && float64(completedTasks)/float64(totalTasks) < 0.5 {
				health = v1alpha1.HealthAtRisk
			}
		}
	}

	// 4. Update Project Status
	project.Status.TaskCount = totalTasks
	project.Status.CompletedCount = completedTasks
	project.Status.Health = health
	project.Status.ObservedGeneration = project.Generation

	meta.SetStatusCondition(&project.Status.Conditions, metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		Reason:             "ProjectAggregated",
		Message:            fmt.Sprintf("Aggregated %d total tasks (%d completed)", totalTasks, completedTasks),
		ObservedGeneration: project.Generation,
	})

	if !reflect.DeepEqual(previousStatus, project.Status) {
		if err := r.Status().Update(ctx, project); err != nil {
			logger.Error(err, "Failed to update Project status")
			return ctrl.Result{}, err
		}
	}

	logger.Info("Successfully aggregated Project statistics",
		"taskCount", totalTasks,
		"completedCount", completedTasks,
		"health", health,
	)
	return ctrl.Result{}, nil
}

// findProjectsForTask maps an incoming Task event to its referenced parent Project
func (r *ProjectReconciler) findProjectsForTask(ctx context.Context, obj client.Object) []reconcile.Request {
	task, ok := obj.(*v1alpha1.Task)
	if !ok || task.Spec.ProjectRef.Name == "" {
		return nil
	}

	ns := task.Spec.ProjectRef.Namespace
	if ns == "" {
		ns = task.Namespace
	}

	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Name:      task.Spec.ProjectRef.Name,
				Namespace: ns,
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager and watches child Tasks
func (r *ProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Project{}).
		Watches(&v1alpha1.Task{}, handler.EnqueueRequestsFromMapFunc(r.findProjectsForTask)).
		Complete(r)
}
