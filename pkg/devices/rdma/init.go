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

import fw "openbce.io/device-manager/pkg/devices/framework"

func init() {
	if err := fw.Register("rdma", New()); err != nil {
		panic("failed to register <rdma> device plugin.")
	}
}

func New() fw.Service {
	return &rdmaSvc{}
}
