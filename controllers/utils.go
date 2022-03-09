package controllers

import (
	"context"

	"github.com/go-logr/logr"
	iomeshv1 "iomesh.com/cdi-iomesh/api/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func (r *IOMeshVolumeReconciler) pvcExists(iomv *iomeshv1.IOMeshVolume) bool {
	pvc := &corev1.PersistentVolumeClaim{}
	if err := r.Get(context.TODO(), types.NamespacedName{Namespace: iomv.Namespace, Name: iomv.Name}, pvc); err != nil {
		if k8serrors.IsNotFound(err) {
			return false
		}
	}
	return true
}

func (r *IOMeshVolumeReconciler) cleanup(log logr.Logger, iomv *iomeshv1.IOMeshVolume) error {
	pvc := &corev1.PersistentVolumeClaim{}
	if err := r.Get(context.TODO(), types.NamespacedName{Namespace: iomv.Namespace, Name: iomv.Name}, pvc); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
	}
	log.Info("Find pvc", pvc.Name)
	return r.Delete(context.TODO(), pvc)
}

func (r *IOMeshVolumeReconciler) newPersistentVolumeClaim(iomv *iomeshv1.IOMeshVolume) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: iomv.Namespace,
			Name:      iomv.Name,
		},
		Spec: *iomv.Spec.PVC,
	}

	pvc.OwnerReferences = []metav1.OwnerReference{
		*metav1.NewControllerRef(iomv, schema.GroupVersionKind{
			Group:   iomeshv1.GroupVersion.Group,
			Version: iomeshv1.GroupVersion.Version,
			Kind:    "IOMeshVolume",
		}),
	}
	return pvc, nil
}
