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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "github.com/ankrsinha/kubetaskops/backend/pkg/apis/v1alpha1"
)

const (
	TaskFinalizer = "kubetaskops.io/finalizer"
)

// TaskReconciler reconciles a Task object
type TaskReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile handles Task lifecycle reconciliation
func (r *TaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("task", req.NamespacedName)

	task := &v1alpha1.Task{}
	if err := r.Get(ctx, req.NamespacedName, task); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("Task resource deleted from cluster")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Task resource")
		return ctrl.Result{}, err
	}
	previousStatus := task.DeepCopy().Status

	// 1. Handle Deletion & Finalizer
	if task.DeletionTimestamp.IsZero() {
		// Resource is active, ensure finalizer is attached
		if !controllerutil.ContainsFinalizer(task, TaskFinalizer) {
			controllerutil.AddFinalizer(task, TaskFinalizer)
			if err := r.Update(ctx, task); err != nil {
				logger.Error(err, "Failed to add finalizer to Task")
				return ctrl.Result{}, err
			}
			logger.Info("Added finalizer to Task", "finalizer", TaskFinalizer)
			return ctrl.Result{}, nil
		}
	} else {
		// Resource is undergoing deletion
		if controllerutil.ContainsFinalizer(task, TaskFinalizer) {
			logger.Info("Executing finalizer cleanup for Task")
			// Cleanup logic (e.g. projection sync in Phase 4) executed here
			controllerutil.RemoveFinalizer(task, TaskFinalizer)
			if err := r.Update(ctx, task); err != nil {
				logger.Error(err, "Failed to remove finalizer from Task")
				return ctrl.Result{}, err
			}
			logger.Info("Removed finalizer from Task, deletion completed")
		}
		return ctrl.Result{}, nil
	}

	// 2. Validate Project Reference
	projNamespace := task.Spec.ProjectRef.Namespace
	if projNamespace == "" {
		projNamespace = task.Namespace
	}

	project := &v1alpha1.Project{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      task.Spec.ProjectRef.Name,
		Namespace: projNamespace,
	}, project)

	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("Referenced Project not found", "project", task.Spec.ProjectRef.Name)
			meta.SetStatusCondition(&task.Status.Conditions, metav1.Condition{
				Type:               "ProjectBound",
				Status:             metav1.ConditionFalse,
				Reason:             "ProjectNotFound",
				Message:            fmt.Sprintf("Referenced Project '%s' does not exist in namespace '%s'", task.Spec.ProjectRef.Name, projNamespace),
				ObservedGeneration: task.Generation,
			})
			meta.SetStatusCondition(&task.Status.Conditions, metav1.Condition{
				Type:               "Ready",
				Status:             metav1.ConditionFalse,
				Reason:             "ProjectNotFound",
				Message:            "Task cannot be Ready until referenced Project exists",
				ObservedGeneration: task.Generation,
			})
			task.Status.ObservedGeneration = task.Generation

			if updateErr := r.updateStatusIfChanged(ctx, task, previousStatus); updateErr != nil {
				logger.Error(updateErr, "Failed to update Task status for missing Project")
				return ctrl.Result{}, updateErr
			}
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		}
		logger.Error(err, "Failed to lookup referenced Project")
		return ctrl.Result{}, err
	}

	meta.SetStatusCondition(&task.Status.Conditions, metav1.Condition{
		Type:               "ProjectBound",
		Status:             metav1.ConditionTrue,
		Reason:             "ProjectBound",
		Message:            fmt.Sprintf("Successfully bound to Project '%s'", project.Name),
		ObservedGeneration: task.Generation,
	})

	// 3. Reconcile Desired Phase to Status Phase
	desiredPhase := task.Spec.DesiredPhase
	if desiredPhase == "" {
		desiredPhase = v1alpha1.PhasePending
	}
	task.Status.Phase = desiredPhase

	meta.SetStatusCondition(&task.Status.Conditions, metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		Reason:             "TaskReconciled",
		Message:            fmt.Sprintf("Task successfully reconciled to phase '%s'", desiredPhase),
		ObservedGeneration: task.Generation,
	})

	// 4. Update Generation & Timestamp Metadata
	task.Status.ObservedGeneration = task.Generation

	if err := r.updateStatusIfChanged(ctx, task, previousStatus); err != nil {
		logger.Error(err, "Failed to update Task status")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully reconciled Task", "phase", task.Status.Phase)
	return ctrl.Result{}, nil
}

// updateStatusIfChanged persists status only when reconciliation changed it.
// This prevents controller-generated status events from causing an update loop.
func (r *TaskReconciler) updateStatusIfChanged(ctx context.Context, task *v1alpha1.Task, previousStatus v1alpha1.TaskStatus) error {
	if reflect.DeepEqual(previousStatus, task.Status) {
		return nil
	}

	now := metav1.Now()
	task.Status.SyncTimestamp = &now
	return r.Status().Update(ctx, task)
}

// SetupWithManager sets up the controller with the Manager.
func (r *TaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Task{}).
		Complete(r)
}
