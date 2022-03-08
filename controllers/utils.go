package controllers

import (
	iomeshv1 "iomesh.com/cdi-iomesh/api/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

func RenderDataVolumeSpec(iomeshVolumeSpec iomeshv1.IOMeshVolumeSpec) (*cdiv1.DataVolumeSpec, error) {
	spec := &cdiv1.DataVolumeSpec{}
	return spec, nil
}
