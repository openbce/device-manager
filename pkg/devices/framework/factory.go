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
	"fmt"
	"os"
	"sync"

	"k8s.io/klog"

	"openbce.io/device-manager/pkg/devices/common"
)

var devices = map[string]Service{}

var mutex sync.Mutex

func Register(name string, svc Service) error {
	mutex.Lock()
	defer mutex.Unlock()

	if _, found := devices[name]; found {
		return fmt.Errorf("duplicated device plugin <%s>", name)
	}

	devices[name] = svc

	return nil
}

func setupEnv() error {
	if err := os.MkdirAll(common.DevicePluginPath, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func Run(opt *Options) error {
	mutex.Lock()
	defer mutex.Unlock()

	if err := setupEnv(); err != nil {
		return err
	}

	for k, v := range devices {
		if err := v.Register(opt); err != nil {
			klog.Errorf("Failed to register <%s> device plugin: %v.", k, err)
			continue
		}

		if err := v.Run(); err != nil {
			klog.Errorf("Failed to start <%s> device plugin: %v.", k, err)
		} else {
			klog.Infof("The device plugin <%s> was started.")
		}
	}

	return nil
}
