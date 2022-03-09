# IOMesh Raw Disk Volume
A toy project, alonside providing volume for containers by iomesh, this project tries to provide raw disks on the host directly.

```
make mainifests
make install
make run
kubectl apply -f config/samples/iomesh_v1_iomeshvolume.yaml
```
It will create/mount a block device in the given Kubernetes node automatically by setting up PVC, Pod, etc.

It will use a iomesh storageclass iomesh-csi-driver to create Kubernetes PVC. The block device could be used outside containerzied application. Please configure node name or storageclass with proper information.

iomesh.com has a good sample to spin up a iomesh cluster.
