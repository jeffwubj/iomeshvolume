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

const (
	// HelperContainerMountPath
	HelperContainerMountPath = "/dev/iomeshblock"
	// HelperContainerName
	HelperContainerName = "pause"
	// HelperContainerImage
	HelperContainerImage = "k8s.gcr.io/pause:3.6"
	// HelperPodTopologyKey
	HelperPodTopologyKey = "kubernetes.io/hostname"
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

func (r *IOMeshVolumeReconciler) podExists(iomv *iomeshv1.IOMeshVolume) bool {
	p := &corev1.Pod{}
	if err := r.Get(context.TODO(), types.NamespacedName{Namespace: iomv.Namespace, Name: iomv.Name}, p); err != nil {
		if k8serrors.IsNotFound(err) {
			return false
		}
	}
	return true
}

func (r *IOMeshVolumeReconciler) cleanup(log logr.Logger, iomv *iomeshv1.IOMeshVolume) error {
	// TODO try best to either clean pvc OR pod

	pvc := &corev1.PersistentVolumeClaim{}
	if err := r.Get(context.TODO(), types.NamespacedName{Namespace: iomv.Namespace, Name: iomv.Name}, pvc); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}
	err := r.Delete(context.TODO(), pvc)
	if err != nil {
		log.Info("failed to clean pvc", pvc.Name)
		return err
	}
	log.Info("pvc", pvc.Name, "is cleaned")

	p := &corev1.Pod{}
	if err := r.Get(context.TODO(), types.NamespacedName{Namespace: iomv.Namespace, Name: iomv.Name}, p); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
	}
	err = r.Delete(context.TODO(), p)
	if err != nil {
		log.Info("failed to clean pod", p.Name)
		return err
	}
	log.Info("pod", p.Name, "is cleaned")

	return nil
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

func (r *IOMeshVolumeReconciler) newHelperPod(iomv *iomeshv1.IOMeshVolume, pvc *corev1.PersistentVolumeClaim) (*corev1.Pod, error) {
	req := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      iomv.Name,
			Namespace: iomv.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  HelperContainerName,
					Image: HelperContainerImage,
					VolumeDevices: []corev1.VolumeDevice{
						{
							Name:       pvc.Name,
							DevicePath: HelperContainerMountPath,
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: pvc.Name,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvc.Name,
							ReadOnly:  false,
						},
					},
				},
			},
			NodeSelector: map[string]string{
				// TODO check Node exists and is schedulable
				HelperPodTopologyKey: iomv.Spec.Node,
			},
		},
	}

	req.OwnerReferences = []metav1.OwnerReference{
		*metav1.NewControllerRef(iomv, schema.GroupVersionKind{
			Group:   iomeshv1.GroupVersion.Group,
			Version: iomeshv1.GroupVersion.Version,
			Kind:    "IOMeshVolume",
		}),
	}
	return req, nil
}
