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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus-community/windows_exporter/internal/headers/virtdisk"
	"github.com/prometheus-community/windows_exporter/internal/pdh"
	"github.com/prometheus-community/windows_exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

// Hyper-V Virtual Storage Device metrics
type collectorVirtualStorageDevice struct {
	perfDataCollectorVirtualStorageDevice *pdh.Collector
	perfDataObjectVirtualStorageDevice    []perfDataCounterValuesVirtualStorageDevice

	virtualStorageDeviceErrorCount               *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Error Count
	virtualStorageDeviceQueueLength              *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Queue Length
	virtualStorageDeviceReadBytes                *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Read Bytes/sec
	virtualStorageDeviceReadOperations           *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Read Operations/Sec
	virtualStorageDeviceWriteBytes               *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Write Bytes/sec
	virtualStorageDeviceWriteOperations          *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Write Operations/Sec
	virtualStorageDeviceLatency                  *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Latency
	virtualStorageDeviceThroughput               *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Throughput
	virtualStorageDeviceNormalizedThroughput     *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Normalized Throughput
	virtualStorageDeviceLowerQueueLength         *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Lower Queue Length
	virtualStorageDeviceLowerLatency             *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\Lower Latency
	virtualStorageDeviceIOQuotaReplenishmentRate *prometheus.Desc // \Hyper-V Virtual Storage Device(*)\IO Quota Replenishment Rate
	virtualStorageDeviceVirtualSizeBytes         *prometheus.Desc // Virtual size of the VHD/VHDX file
	virtualStorageDevicePhysicalSizeBytes        *prometheus.Desc // Physical size of the VHD/VHDX file on disk
}

type perfDataCounterValuesVirtualStorageDevice struct {
	Name string

	VirtualStorageDeviceErrorCount               float64 `perfdata:"Error Count"`
	VirtualStorageDeviceQueueLength              float64 `perfdata:"Queue Length"`
	VirtualStorageDeviceReadBytes                float64 `perfdata:"Read Bytes/sec"`
	VirtualStorageDeviceReadOperations           float64 `perfdata:"Read Count"`
	VirtualStorageDeviceWriteBytes               float64 `perfdata:"Write Bytes/sec"`
	VirtualStorageDeviceWriteOperations          float64 `perfdata:"Write Count"`
	VirtualStorageDeviceLatency                  float64 `perfdata:"Latency"`
	VirtualStorageDeviceThroughput               float64 `perfdata:"Throughput"`
	VirtualStorageDeviceNormalizedThroughput     float64 `perfdata:"Normalized Throughput"`
	VirtualStorageDeviceLowerQueueLength         float64 `perfdata:"Lower Queue Length"`
	VirtualStorageDeviceLowerLatency             float64 `perfdata:"Lower Latency"`
	VirtualStorageDeviceIOQuotaReplenishmentRate float64 `perfdata:"IO Quota Replenishment Rate"`
}

func (c *Collector) buildVirtualStorageDevice() error {
	var err error

	c.perfDataCollectorVirtualStorageDevice, err = pdh.NewCollector[perfDataCounterValuesVirtualStorageDevice](c.logger, pdh.CounterTypeRaw, "Hyper-V Virtual Storage Device", pdh.InstancesAll)
	if err != nil {
		return fmt.Errorf("failed to create Hyper-V Virtual Storage Device collector: %w", err)
	}

	c.virtualStorageDeviceErrorCount = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_error_count_total"),
		"Represents the total number of errors that have occurred on this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceQueueLength = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_queue_length"),
		"Represents the average queue length on this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceReadBytes = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_bytes_read"),
		"Represents the total number of bytes that have been read on this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceReadOperations = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_operations_read_total"),
		"Represents the total number of read operations that have occurred on this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceWriteBytes = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_bytes_written"),
		"Represents the total number of bytes that have been written on this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceWriteOperations = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_operations_written_total"),
		"Represents the total number of write operations that have occurred on this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceLatency = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_latency_seconds"),
		"Represents the average IO transfer latency for this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceThroughput = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_throughput_total"),
		"Represents the total number of 8KB IO transfers completed by this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceNormalizedThroughput = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_normalized_throughput"),
		"Represents the average number of IO transfers completed by this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceLowerQueueLength = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_lower_queue_length"),
		"Represents the average queue length on the underlying storage subsystem for this device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceLowerLatency = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_lower_latency_seconds"),
		"Represents the average IO transfer latency on the underlying storage subsystem for this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceIOQuotaReplenishmentRate = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "io_quota_replenishment_rate"),
		"Represents the IO quota replenishment rate for this virtual device.",
		[]string{"device"},
		nil,
	)
	c.virtualStorageDeviceVirtualSizeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_size_bytes"),
		"Virtual size of the VHD/VHDX file in bytes.",
		[]string{"device", "path"},
		nil,
	)
	c.virtualStorageDevicePhysicalSizeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "virtual_storage_device_physical_size_bytes"),
		"Physical size of the VHD/VHDX file on disk in bytes.",
		[]string{"device", "path"},
		nil,
	)

	return nil
}

func (c *Collector) collectVirtualStorageDevice(ch chan<- prometheus.Metric) error {
	err := c.perfDataCollectorVirtualStorageDevice.Collect(&c.perfDataObjectVirtualStorageDevice)
	if err != nil {
		return fmt.Errorf("failed to collect Hyper-V Virtual Storage Device metrics: %w", err)
	}

	for _, data := range c.perfDataObjectVirtualStorageDevice {
		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceErrorCount,
			prometheus.CounterValue,
			data.VirtualStorageDeviceErrorCount,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceQueueLength,
			prometheus.GaugeValue,
			data.VirtualStorageDeviceQueueLength,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceReadBytes,
			prometheus.CounterValue,
			data.VirtualStorageDeviceReadBytes,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceReadOperations,
			prometheus.CounterValue,
			data.VirtualStorageDeviceReadOperations,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceWriteBytes,
			prometheus.CounterValue,
			data.VirtualStorageDeviceWriteBytes,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceWriteOperations,
			prometheus.CounterValue,
			data.VirtualStorageDeviceWriteOperations,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceLatency,
			prometheus.GaugeValue,
			data.VirtualStorageDeviceLatency,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceThroughput,
			prometheus.CounterValue,
			data.VirtualStorageDeviceThroughput,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceNormalizedThroughput,
			prometheus.GaugeValue,
			data.VirtualStorageDeviceNormalizedThroughput,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceLowerQueueLength,
			prometheus.GaugeValue,
			data.VirtualStorageDeviceLowerQueueLength,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceLowerLatency,
			prometheus.GaugeValue,
			data.VirtualStorageDeviceLowerLatency,
			data.Name,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceIOQuotaReplenishmentRate,
			prometheus.GaugeValue,
			data.VirtualStorageDeviceIOQuotaReplenishmentRate,
			data.Name,
		)

		// Attempt to get disk size information
		// The Name field contains the encoded path to the VHD/VHDX file
		diskPath := c.resolveVirtualDiskPath(data.Name)

		// Always emit metrics with -1 to indicate "unknown" if we can't get the size
		virtualSize := float64(-1)
		physicalSize := float64(-1)
		resolvedPath := "unknown"

		if diskPath != "" {
			resolvedPath = diskPath
			vSize, pSize, err := virtdisk.GetVirtualDiskSize(diskPath)
			if err == nil {
				virtualSize = float64(vSize)
				physicalSize = float64(pSize)
			} else {
				// Log the error for debugging but continue processing other devices
				c.logger.Debug("Failed to get virtual disk size",
					"device", data.Name,
					"path", diskPath,
					"error", err,
				)
			}
		} else {
			// Log when we can't resolve the path for debugging
			c.logger.Debug("Unable to resolve virtual disk path",
				"device", data.Name,
			)
		}

		// Always emit the size metrics (with -1 if unknown)
		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDeviceVirtualSizeBytes,
			prometheus.GaugeValue,
			virtualSize,
			data.Name,
			resolvedPath,
		)

		ch <- prometheus.MustNewConstMetric(
			c.virtualStorageDevicePhysicalSizeBytes,
			prometheus.GaugeValue,
			physicalSize,
			data.Name,
			resolvedPath,
		)
	}

	return nil
}

// resolveVirtualDiskPath attempts to resolve the full path to a VHD/VHDX file
// based on the performance counter instance name.
//
// The instance name often contains the encoded path itself, where:
//   - Backslashes (\) are replaced with hyphens (-)
//   - Drive letter colon (:) becomes (:-)
//   - Prefix \\?\ becomes --?-
//
// Example: "--?-C:-ClusterStorage-Volume-VM-disk.vhdx"
// Becomes: "C:\ClusterStorage\Volume\VM\disk.vhdx"
//
// To customize VHD search paths for fallback, set the HYPERV_VHD_PATHS environment variable
// with semicolon-separated paths. Example:
//
//	HYPERV_VHD_PATHS=D:\VMs;E:\ClusterStorage\Volume1
//
// Enable debug logging to troubleshoot path resolution issues.
func (c *Collector) resolveVirtualDiskPath(instanceName string) string {
	// First, try to decode the path from the instance name itself
	// Performance counter instance names encode the full path
	decodedPath := decodeVirtualDiskPath(instanceName)
	if decodedPath != "" {
		c.logger.Debug("Decoded virtual disk path",
			"device", instanceName,
			"decodedPath", decodedPath,
		)
		// Verify the decoded path exists
		if _, err := os.Stat(decodedPath); err == nil {
			return decodedPath
		} else {
			c.logger.Debug("Decoded path does not exist",
				"decodedPath", decodedPath,
				"error", err,
			)
		}
	}

	// Fallback to searching common locations
	// Common Hyper-V virtual disk storage locations
	// Can be customized via HYPERV_VHD_PATHS environment variable (semicolon-separated)
	commonPaths := []string{
		`C:\ClusterStorage`,
		`C:\ProgramData\Microsoft\Windows\Hyper-V`,
		`C:\ProgramData\Microsoft\Windows\Hyper-V\Virtual Hard Disks`,
		`C:\Users\Public\Documents\Hyper-V\Virtual Hard Disks`,
		`D:\Hyper-V`,
		`D:\Hyper-V\Virtual Hard Disks`,
		`E:\Hyper-V`,
		`E:\Hyper-V\Virtual Hard Disks`,
	}

	// Allow custom paths from environment variable
	if customPaths := os.Getenv("HYPERV_VHD_PATHS"); customPaths != "" {
		customPathsList := strings.Split(customPaths, ";")
		// Prepend custom paths so they're checked first
		commonPaths = append(customPathsList, commonPaths...)
	}

	// Try to extract a meaningful filename from the instance name
	// Instance names might be in format like "VMName_DiskName" or just "DiskName"
	possibleNames := []string{
		instanceName + ".vhdx",
		instanceName + ".vhd",
		instanceName + ".vhdset",
	}

	// Also try splitting on underscore and using the last part
	parts := strings.Split(instanceName, "_")
	if len(parts) > 1 {
		lastPart := parts[len(parts)-1]
		possibleNames = append(possibleNames,
			lastPart+".vhdx",
			lastPart+".vhd",
			lastPart+".vhdset",
		)
	}

	// Try the full instance name as-is if it looks like a filename
	if strings.Contains(instanceName, ".vhd") {
		possibleNames = append([]string{instanceName}, possibleNames...)
	}

	// Search in common paths
	for _, basePath := range commonPaths {
		for _, name := range possibleNames {
			// Try direct path
			fullPath := filepath.Join(basePath, name)
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath
			}

			// Try searching in subdirectories (up to 2 levels deep for VM folders)
			pattern := filepath.Join(basePath, "*", name)
			matches, err := filepath.Glob(pattern)
			if err == nil && len(matches) > 0 {
				return matches[0]
			}

			// Try 2 levels deep
			pattern = filepath.Join(basePath, "*", "*", name)
			matches, err = filepath.Glob(pattern)
			if err == nil && len(matches) > 0 {
				return matches[0]
			}
		}
	}

	return ""
}

// decodeVirtualDiskPath decodes a Hyper-V performance counter instance name
// into a Windows file path.
//
// The encoding format used by Hyper-V performance counters:
//   - Backslashes (\) are replaced with hyphens (-)
//   - Drive letter colon (C:\) becomes (C:-)
//   - UNC prefix (\\?\) becomes (--?-)
//
// Since directory names can contain hyphens, we try multiple interpretations
// and return the first one that exists on disk.
//
// Examples:
//   - "--?-C:-ClusterStorage-Volume-VM-disk.vhdx" -> "C:\ClusterStorage\Volume\VM\disk.vhdx"
//   - "C:-ClusterStorage-HAMC010011-FTT2-CSV01-VM-disk.vhdx" -> "C:\ClusterStorage\HAMC010011-FTT2-CSV01\VM\disk.vhdx"
func decodeVirtualDiskPath(instanceName string) string {
	if instanceName == "" {
		return ""
	}

	// Check if this looks like an encoded path
	// Encoded paths typically contain ":-" (drive letter) or start with "--?-"
	if !strings.Contains(instanceName, ":-") && !strings.HasPrefix(instanceName, "--?-") {
		return ""
	}

	path := instanceName

	// Remove UNC prefix if present: --?- becomes \\?\
	// For simplicity, we'll just remove it as Windows can work without it
	path = strings.TrimPrefix(path, "--?-")

	// Handle drive letter: "C:-" becomes "C:\"
	// Find pattern like "X:-" where X is a letter
	if len(path) < 3 || path[1] != ':' || path[2] != '-' {
		return ""
	}

	driveLetter := string(path[0])
	remainingPath := path[3:] // Everything after "C:-"

	// The simple approach: replace all hyphens with backslashes
	// This works when directory names don't contain hyphens
	simplePath := driveLetter + `:\` + strings.ReplaceAll(remainingPath, "-", `\`)

	// Try the simple path first
	if fileExists(simplePath) {
		return simplePath
	}

	// If simple approach didn't work, directory names likely contain hyphens
	// Split by hyphens and try progressively merging parts
	parts := strings.Split(remainingPath, "-")

	// Try to build path by treating consecutive parts as one directory name
	// when separated hyphens don't create a valid path
	return tryPathCombinations(driveLetter+`:\`, parts)
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// tryPathCombinations tries different combinations of joining parts
// to find a valid file path, accounting for hyphens in directory names
func tryPathCombinations(basePath string, parts []string) string {
	if len(parts) == 0 {
		return ""
	}

	// Start building path from the beginning
	var currentPath string

	for i := 0; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}

		// Try adding this part as a new path component
		testPath := filepath.Join(basePath, parts[i])

		// If this is the last part and it has a file extension, check if it exists
		if i == len(parts)-1 && (strings.Contains(parts[i], ".vhd") || strings.Contains(parts[i], ".vmgs")) {
			if currentPath != "" {
				finalPath := filepath.Join(currentPath, parts[i])
				if fileExists(finalPath) {
					return finalPath
				}
			}
			// Also try without the accumulated path
			if fileExists(testPath) {
				return testPath
			}
		}

		// Check if this directory exists
		if fileInfo, err := os.Stat(testPath); err == nil && fileInfo.IsDir() {
			currentPath = testPath
			basePath = currentPath
		} else if i < len(parts)-1 {
			// Try merging with next part (this part might have a hyphen in its name)
			merged := parts[i] + "-" + parts[i+1]
			testPath = filepath.Join(basePath, merged)
			if fileInfo, err := os.Stat(testPath); err == nil && fileInfo.IsDir() {
				currentPath = testPath
				basePath = currentPath
				i++ // Skip the next part since we merged it
			}
		}
	}

	return ""
}
