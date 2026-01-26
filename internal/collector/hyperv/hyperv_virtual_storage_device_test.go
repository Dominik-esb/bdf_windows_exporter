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

package hyperv

import (
	"testing"
)

func TestDecodeVirtualDiskPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ClusterStorage with UNC prefix",
			input:    "--?-C:-ClusterStorage-BCNC010001-FTT1-CSV01-BCNA010003-Virtual Machines-B8DA40B3-5BD9-40AB-9C51-66DBA1F5EE35.vmgs",
			expected: `C:\ClusterStorage\BCNC010001-FTT1-CSV01\BCNA010003\Virtual Machines\B8DA40B3-5BD9-40AB-9C51-66DBA1F5EE35.vmgs`,
		},
		{
			name:     "ClusterStorage without UNC prefix",
			input:    "C:-ClusterStorage-HAMC010011-FTT2-CSV01-VM01-disk.vhdx",
			expected: `C:\ClusterStorage\HAMC010011-FTT2-CSV01\VM01\disk.vhdx`,
		},
		{
			name:     "Simple path",
			input:    "C:-Users-Public-Documents-Hyper-V-Virtual Hard Disks-test.vhdx",
			expected: `C:\Users\Public\Documents\Hyper-V\Virtual Hard Disks\test.vhdx`,
		},
		{
			name:     "D drive",
			input:    "D:-VMs-Production-VM-disk1.vhd",
			expected: `D:\VMs\Production\VM\disk1.vhd`,
		},
		{
			name:     "Non-encoded path returns empty",
			input:    "VMName_DiskName",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}










}	}		})			}				t.Errorf("decodeVirtualDiskPath(%q) = %q, want %q", tt.input, result, tt.expected)			if result != tt.expected {			result := decodeVirtualDiskPath(tt.input)		t.Run(tt.name, func(t *testing.T) {	for _, tt := range tests {
