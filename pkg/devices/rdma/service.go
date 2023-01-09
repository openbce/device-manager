/*
Copyright 2022 The openBCE Authors.
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

package rdma

import (
	"context"
	"net"
	"syscall"

	"google.golang.org/grpc"

	"k8s.io/klog"
	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

	fw "openbce.io/device-manager/pkg/devices/framework"
)

type rdmaSvc struct {
}

func (r *rdmaSvc) GetDevicePluginOptions(ctx context.Context, empty *v1beta1.Empty) (*v1beta1.DevicePluginOptions, error) {
	return &v1beta1.DevicePluginOptions{
		PreStartRequired: true,
	}, nil
}

func (r *rdmaSvc) ListAndWatch(empty *v1beta1.Empty, server v1beta1.DevicePlugin_ListAndWatchServer) error {
	var stop <-chan struct{}

	devs, err := getDevices()
	if err != nil {
		return err
	}

	resp := &v1beta1.ListAndWatchResponse{}
	for k, _ := range devs {
		dev := &v1beta1.Device{
			ID:     k,
			Health: v1beta1.Healthy,
		}
		resp.Devices = append(resp.Devices, dev)

		klog.Infof("dev info: %v", dev)
	}

	if err = server.Send(resp); err != nil {
		return err
	}

	// TODO(k82cn): update device healthy accordingly.
	for {
		select {
		case <-stop:
			return nil
		}
	}
}

func (r *rdmaSvc) GetPreferredAllocation(ctx context.Context, request *v1beta1.PreferredAllocationRequest) (*v1beta1.PreferredAllocationResponse, error) {
	resp := &v1beta1.PreferredAllocationResponse{}

	ids, err := getDeviceIDs()
	if err != nil {
		return resp, err
	}

	for i := 0; i < len(request.ContainerRequests); i++ {
		resp.ContainerResponses = append(resp.ContainerResponses, &v1beta1.ContainerPreferredAllocationResponse{
			DeviceIDs: ids,
		})
	}

	return resp, nil
}

func (r *rdmaSvc) Allocate(ctx context.Context, request *v1beta1.AllocateRequest) (*v1beta1.AllocateResponse, error) {
	resp := &v1beta1.AllocateResponse{}

	devMap, err := getDevices()
	if err != nil {
		return resp, err
	}

	for _, r := range request.ContainerRequests {
		cr := &v1beta1.ContainerAllocateResponse{}
		for _, id := range r.DevicesIDs {
			d, found := devMap[id]
			if !found {
				klog.Errorf("Failed to find device <%s> in RDMA device plugin.", id)
				continue
			}
			ibPath := IbPath + d.DevName
			cr.Devices = append(cr.Devices, &v1beta1.DeviceSpec{
				HostPath:      ibPath,
				ContainerPath: ibPath,
				Permissions:   "rw",
			})
			cr.Devices = append(cr.Devices, &v1beta1.DeviceSpec{
				HostPath:      d.DevPath,
				ContainerPath: d.DevPath,
				Permissions:   "rw",
			})
			cr.Devices = append(cr.Devices, &v1beta1.DeviceSpec{
				HostPath:      d.IbvDevPath,
				ContainerPath: d.IbvDevPath,
				Permissions:   "rw",
			})

			klog.Infof("Allocate %v ", d)
		}

		resp.ContainerResponses = append(resp.ContainerResponses, cr)
	}

	return resp, nil
}

func (r *rdmaSvc) PreStartContainer(ctx context.Context, request *v1beta1.PreStartContainerRequest) (*v1beta1.PreStartContainerResponse, error) {
	return &v1beta1.PreStartContainerResponse{}, nil
}

func (r *rdmaSvc) Register(opt *fw.Options) error {
	// Register to kubelet
	rr := &v1beta1.RegisterRequest{
		Version:      v1beta1.Version,
		Endpoint:     Entrypoint,
		ResourceName: opt.Vendor + ResourceName,
		Options: &v1beta1.DevicePluginOptions{
			PreStartRequired: true,
		},
	}

	if err := fw.RegisterDevicePlugin(rr); err != nil {
		return err
	}

	return nil
}

func (r *rdmaSvc) Run() error {
	// Start RDMA device plugin service
	s := grpc.NewServer()
	v1beta1.RegisterDevicePluginServer(s, r)

	// Remove sock firstly
	syscall.Unlink(Socket)

	klog.Infof("Start RDMA device plugin at <%s>.", Socket)
	ibsvc, err := net.Listen("unix", Socket)
	if err != nil {
		return err
	}

	if err = s.Serve(ibsvc); err != nil {
		return err
	}

	return nil
}
