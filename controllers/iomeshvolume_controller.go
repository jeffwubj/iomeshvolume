/*
Copyright 2022.

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

package controllers

import (
	"context"

	iomeshv1 "iomesh.com/cdi-iomesh/api/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// ErrResourceExists provides a const to indicate a resource exists error
	ErrResourceExists = "ErrResourceExists"
	// MessageResourceExists provides a const to form a resource exists error message
	MessageResourceExists = "Resource %q already exists and is not managed by IOMeshVolume"
)

// IOMeshVolumeReconciler reconciles a IOMeshVolume object
type IOMeshVolumeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=poc.iomesh.com,resources=iomeshvolumes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=poc.iomesh.com,resources=iomeshvolumes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=poc.iomesh.com,resources=iomeshvolumes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the IOMeshVolume object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *IOMeshVolumeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// TODO get and show node & mount point

	l := log.FromContext(ctx)

	iomv := &iomeshv1.IOMeshVolume{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, iomv)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// New IOMeshVolume
			return ctrl.Result{}, nil
		}
		l.Error(err, "Failed to get IOMeshVolume", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if iomv.DeletionTimestamp != nil {
		l.Info("IOMeshVolume marked for deletion, cleaning up")
		err := r.cleanup(l, iomv)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	pvcExists := r.pvcExists(iomv)
	pvc := &corev1.PersistentVolumeClaim{}
	// TODO handle the case pvc exists but no managed by this iomeshvolume
	if !pvcExists {
		pvc, err = r.newPersistentVolumeClaim(iomv)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err := r.Create(context.TODO(), pvc); err != nil {
			return reconcile.Result{}, err
		}
	}

	// TODO update Pod with new node location
	// TODO handle the case pod exists but no managed by this iomeshvolume

	podExists := r.podExists(iomv)
	if !podExists {
		p, err := r.newHelperPod(iomv, pvc)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err := r.Create(context.TODO(), p); err != nil {
			return reconcile.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IOMeshVolumeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&iomeshv1.IOMeshVolume{}).
		Complete(r)
}
