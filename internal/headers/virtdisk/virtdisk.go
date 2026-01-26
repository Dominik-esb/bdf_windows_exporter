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
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

// GetVirtualDiskSize returns the virtual and physical size of a VHD/VHDX file
func GetVirtualDiskSize(path string) (virtualSize uint64, physicalSize uint64, err error) {
	// Convert path to UTF16
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert path to UTF16: %w", err)
	}

	// Determine storage type based on file extension
	storageType := VIRTUAL_STORAGE_TYPE{
		DeviceID: VIRTUAL_STORAGE_TYPE_DEVICE_UNKNOWN,
		VendorID: VIRTUAL_STORAGE_TYPE_VENDOR_MICROSOFT,
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".vhd":
		storageType.DeviceID = VIRTUAL_STORAGE_TYPE_DEVICE_VHD
	case ".vhdx":
		storageType.DeviceID = VIRTUAL_STORAGE_TYPE_DEVICE_VHDX
	case ".vhdset":
		storageType.DeviceID = VIRTUAL_STORAGE_TYPE_DEVICE_VHDSET
	case ".iso":
		storageType.DeviceID = VIRTUAL_STORAGE_TYPE_DEVICE_ISO
	}

	// Open the virtual disk with flags that allow opening even when in use
	// Use READ access mask which includes ATTACH_RO, DETACH, and GET_INFO
	var handle windows.Handle
	err = OpenVirtualDisk(
		&storageType,
		pathPtr,
		VIRTUAL_DISK_ACCESS_READ,
		OPEN_VIRTUAL_DISK_FLAG_NO_PARENTS|OPEN_VIRTUAL_DISK_FLAG_CACHED_IO, // Allow opening VHDs in use with cached I/O
		nil,
		&handle,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open virtual disk: %w", err)
	}
	defer windows.CloseHandle(handle)

	// Get size information
	diskInfo := GET_VIRTUAL_DISK_INFO{
		Version: GET_VIRTUAL_DISK_INFO_SIZE,
	}

	diskInfoSize := uint32(unsafe.Sizeof(diskInfo))
	var sizeUsed uint32

	err = GetVirtualDiskInformation(
		handle,
		&diskInfoSize,
		&diskInfo,
		&sizeUsed,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get virtual disk information: %w", err)
	}

	return diskInfo.Size.VirtualSize, diskInfo.Size.PhysicalSize, nil
}
