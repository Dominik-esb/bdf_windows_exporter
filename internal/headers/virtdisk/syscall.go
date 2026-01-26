// SPDX-License-Identifier: Apache-2.0
//
// Copyright The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build windows

package virtdisk

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	virtdisk = windows.NewLazySystemDLL("virtdisk.dll")

	procOpenVirtualDisk           = virtdisk.NewProc("OpenVirtualDisk")
	procGetVirtualDiskInformation = virtdisk.NewProc("GetVirtualDiskInformation")
)

// OpenVirtualDisk opens a virtual hard disk (VHD or VHDX) or CD or DVD image file (ISO) for use.
func OpenVirtualDisk(
	virtualStorageType *VIRTUAL_STORAGE_TYPE,
	path *uint16,
	virtualDiskAccessMask VIRTUAL_DISK_ACCESS_MASK,
	flags OPEN_VIRTUAL_DISK_FLAG,
	parameters *OPEN_VIRTUAL_DISK_PARAMETERS,
	handle *windows.Handle,
) error {
	r1, _, err := procOpenVirtualDisk.Call(
		uintptr(unsafe.Pointer(virtualStorageType)),
		uintptr(unsafe.Pointer(path)),
		uintptr(virtualDiskAccessMask),
		uintptr(flags),
		uintptr(unsafe.Pointer(parameters)),
		uintptr(unsafe.Pointer(handle)),
	)
	if r1 != 0 {
		return err
	}
	return nil
}

// GetVirtualDiskInformation retrieves information about a VHD.
func GetVirtualDiskInformation(
	virtualDiskHandle windows.Handle,
	virtualDiskInfoSize *uint32,
	virtualDiskInfo *GET_VIRTUAL_DISK_INFO,
	sizeUsed *uint32,
) error {
	r1, _, err := procGetVirtualDiskInformation.Call(
		uintptr(virtualDiskHandle),
		uintptr(unsafe.Pointer(virtualDiskInfoSize)),
		uintptr(unsafe.Pointer(virtualDiskInfo)),
		uintptr(unsafe.Pointer(sizeUsed)),
	)
	if r1 != 0 {
		return err
	}
	return nil
}
