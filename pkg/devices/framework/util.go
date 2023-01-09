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

package framework

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"

	"k8s.io/klog"
	"k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func RegisterDevicePlugin(req *v1beta1.RegisterRequest) error {
	klog.Infof("Connect to kubelet at <%s>.", v1beta1.KubeletSocket)

	timeout := 3 * time.Second
	conn, err := grpc.Dial(v1beta1.KubeletSocket, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		return err
	}

	client := v1beta1.NewRegistrationClient(conn)
	if _, err = client.Register(context.Background(), req); err != nil {
		return err
	}

	klog.Infof("Register to kubelet successfully.")

	return nil
}
