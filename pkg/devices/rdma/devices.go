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

// #cgo LDFLAGS: -libverbs
// #include <stdlib.h>
// #include <infiniband/verbs.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Device is a Go struct for ibv_device in libibverbs.
type Device struct {
	Name       string
	DevName    string
	DevPath    string
	IbvDevPath string
}

func foreachDevice(callback func(device *Device)) error {
	var c_num C.int
	var c_devptr *C.struct_ibv_device

	// Get device list by libibverbs
	c_devList := C.ibv_get_device_list(&c_num)
	if c_devList == nil {
		return fmt.Errorf("failed to get IB devices list")
	}
	defer C.ibv_free_device_list(c_devList)

	// Store device info into Go structure
	cur := uintptr(unsafe.Pointer(c_devList))
	for i := 0; i <= int(c_num); i++ {
		dev := *(**C.struct_ibv_device)(unsafe.Pointer(cur))
		if dev == nil {
			continue
		}

		d := &Device{
			Name:       C.GoString(&(*dev).name[0]),
			DevName:    C.GoString(&(*dev).dev_name[0]),
			DevPath:    C.GoString(&(*dev).dev_path[0]),
			IbvDevPath: C.GoString(&(*dev).ibdev_path[0]),
		}

		callback(d)
		cur += unsafe.Sizeof(c_devptr)
	}

	return nil
}

func getDevices() (map[string]*Device, error) {
	devMap := map[string]*Device{}
	if err := foreachDevice(func(d *Device) {
		devMap[d.Name] = d
	}); err != nil {
		return nil, err
	}

	return devMap, nil
}

func getDeviceIDs() ([]string, error) {
	var ids []string
	if err := foreachDevice(func(d *Device) {
		ids = append(ids, d.Name)
	}); err != nil {
		return nil, err
	}

	return ids, nil
}
