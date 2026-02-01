package device

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/ArkaprabhaChakraborty/winx/handle"
	"github.com/ArkaprabhaChakraborty/winx/service"
)

// TestClfsDriver tests interactions with the Common Log File System driver (clfs.sys)
func TestClfsDriver(t *testing.T) {
	t.Run("DiscoverClfsDevice", func(t *testing.T) {
		t.Log("Searching for CLFS driver...")

		// First, search for symbolic links with "clfs" in the name
		t.Log("Searching DOS device symbolic links...")
		symlinks, err := FindSymbolicLinksByPattern("clfs")
		if err != nil {
			t.Logf("Error searching symbolic links: %v", err)
		} else if len(symlinks) > 0 {
			t.Logf("Found %d CLFS symbolic link(s):", len(symlinks))
			for device, targets := range symlinks {
				t.Logf("  %s ->", device)
				for _, target := range targets {
					t.Logf("    %s", target)
				}
			}
		} else {
			t.Log("No CLFS symbolic links found")
		}

		// Also search for devices with service name "CLFS"
		t.Log("Searching for CLFS devices by service name...")
		devices, err := FindDevicesByService("CLFS")
		if err != nil {
			t.Logf("Error searching for CLFS devices: %v", err)
		}

		if len(devices) == 0 {
			t.Log("No CLFS devices found through service name search")
		} else {
			t.Logf("Found %d CLFS device(s):", len(devices))
			for i, device := range devices {
				t.Logf("  Device %d:", i)
				if device.Description != "" {
					t.Logf("    Description: %s", device.Description)
				}
				if device.Service != "" {
					t.Logf("    Service: %s", device.Service)
				}
				if device.Class != "" {
					t.Logf("    Class: %s", device.Class)
				}
				if device.HardwareID != "" {
					t.Logf("    Hardware ID: %s", device.HardwareID)
				}
			}
		}
	})

	t.Run("OpenClfsDevice", func(t *testing.T) {
		var workingPath string
		var hDevice handle.HANDLE
		var err error

		// Get common CLFS device paths
		clfsDevicePaths := GetDriverDevicePaths("CLFS")
		// Add some known CLFS-specific paths
		clfsDevicePaths = append(clfsDevicePaths, `\\.\ClfsLog`, `\\Device\clfs`, `\\.\clfscntrl`)

		// Try each known path
		for _, devicePath := range clfsDevicePaths {
			t.Logf("Trying device path: %s", devicePath)

			hDevice, err = CreateFile(
				devicePath,
				GENERIC_READ|GENERIC_WRITE,
				FILE_SHARE_READ|FILE_SHARE_WRITE,
				nil,
				OPEN_EXISTING,
				0,
				0,
			)

			if err == nil && hDevice != 0 && uintptr(hDevice) != INVALID_HANDLE_VALUE {
				workingPath = devicePath
				t.Logf("Successfully opened CLFS device at %s with handle: 0x%x", devicePath, hDevice)
				break
			}

			t.Logf("  Failed: %v", err)
		}

		if workingPath == "" {
			t.Log("Could not open any CLFS device path (may not exist or requires admin)")
			t.Skip("Skipping CLFS tests - device not accessible")
			return
		}
		defer CloseHandle(hDevice)
	})

	t.Run("ClfsDeviceIoControl", func(t *testing.T) {
		var hDevice handle.HANDLE
		var err error

		// Get common CLFS device paths
		clfsDevicePaths := GetDriverDevicePaths("CLFS")
		clfsDevicePaths = append(clfsDevicePaths, `\\.\ClfsLog`, `\\Device\clfs`)

		// Find a working device path
		for _, devicePath := range clfsDevicePaths {
			hDevice, err = OpenDeviceReadWrite(devicePath)
			if err == nil && hDevice != 0 && uintptr(hDevice) != INVALID_HANDLE_VALUE {
				break
			}
		}

		if hDevice == 0 || uintptr(hDevice) == INVALID_HANDLE_VALUE {
			t.Skip("CLFS device not accessible")
			return
		}
		defer CloseHandle(hDevice)

		// CLFS IOCTL codes (FILE_DEVICE_FILE_SYSTEM = 0x00000009)
		const FILE_DEVICE_FILE_SYSTEM = 0x00000009

		ioctlCodes := []struct {
			name string
			code uint32
		}{
			// Common file system IOCTLs that CLFS might respond to
			{"FSCTL_GET_REPARSE_POINT", CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 42, METHOD_BUFFERED, FILE_ANY_ACCESS)},
			{"FSCTL_GET_OBJECT_ID", CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 39, METHOD_BUFFERED, FILE_ANY_ACCESS)},
			{"FSCTL_QUERY_ALLOCATED_RANGES", CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 51, METHOD_NEITHER, FILE_READ_ACCESS)},
		}

		for _, ioctl := range ioctlCodes {
			t.Run(ioctl.name, func(t *testing.T) {
				var bytesReturned uint32
				outBuffer := make([]byte, 1024)

				success, err := DeviceIoControl(
					hDevice,
					ioctl.code,
					nil,
					0,
					unsafe.Pointer(&outBuffer[0]),
					uint32(len(outBuffer)),
					&bytesReturned,
					nil,
				)

				if success {
					t.Logf("IOCTL %s (0x%08X) succeeded, returned %d bytes", ioctl.name, ioctl.code, bytesReturned)
					// Log first few bytes if any returned
					if bytesReturned > 0 {
						displayLen := bytesReturned
						if displayLen > 16 {
							displayLen = 16
						}
						t.Logf("  First bytes: % X", outBuffer[:displayLen])
					}
				} else {
					t.Logf("IOCTL %s (0x%08X) failed: %v (expected for most IOCTLs)", ioctl.name, ioctl.code, err)
				}
			})
		}
	})

	// New Phase 2 Tests: IOCTL Discovery and Driver Loading
	t.Run("LoadClfsDriverFromSystem32", func(t *testing.T) {
		// Import service package for driver loading
		t.Log("Testing CLFS driver loading from system32...")

		driverPath := `C:\Windows\System32\drivers\CLFS.SYS`
		serviceName := "CLFS_Test_Winx"

		t.Logf("Attempting to load CLFS driver from: %s", driverPath)
		hService, err := LoadDriver(driverPath, serviceName)
		if err != nil {
			t.Logf("Failed to load CLFS driver: %v", err)
			t.Skip("Requires administrator privileges to load drivers")
			return
		}

		t.Logf("Successfully loaded CLFS driver with service handle: 0x%X", hService)

		// Query driver status
		var status service.SERVICE_STATUS
		ok, err := service.QueryServiceStatus(hService, &status)
		if ok {
			t.Logf("Driver status: CurrentState=%d (1=STOPPED, 4=RUNNING)", status.CurrentState)
		}

		// Cleanup: Unload the driver
		t.Cleanup(func() {
			t.Log("Cleaning up: Unloading CLFS driver")
			err := UnloadDriver(hService)
			if err != nil {
				t.Logf("Warning: Failed to unload driver: %v", err)
			}
		})
	})

	t.Run("DiscoverClfsIOCTLs", func(t *testing.T) {
		// Try to open CLFS device
		var hDevice handle.HANDLE
		var foundDevice bool

		clfsDevicePaths := GetDriverDevicePaths("CLFS")
		clfsDevicePaths = append(clfsDevicePaths, `\\.\CLFS`, `\\.\ClfsLog`)

		for _, path := range clfsDevicePaths {
			h, err := OpenDevice(path, GENERIC_READ|GENERIC_WRITE)
			if err == nil {
				hDevice = h
				foundDevice = true
				t.Logf("Opened CLFS device: %s", path)
				break
			}
		}

		if !foundDevice {
			t.Skip("CLFS device not accessible")
			return
		}
		defer CloseHandle(hDevice)

		// Discover IOCTLs for FILE_DEVICE_FILE_SYSTEM (common for CLFS)
		const FILE_DEVICE_FILE_SYSTEM = 0x00000009
		t.Logf("Scanning for valid IOCTL codes (device type: FILE_DEVICE_FILE_SYSTEM)...")
		results := DiscoverIOCTLsByDeviceType(hDevice, FILE_DEVICE_FILE_SYSTEM)

		t.Logf("Found %d valid IOCTL codes", len(results))
		if len(results) > 0 {
			t.Log("Valid IOCTLs found:")
			for i, result := range results {
				if i >= 10 {
					t.Logf("  ... and %d more", len(results)-10)
					break
				}
				t.Logf("  [%d] %s (bytes: %d)", i+1, FormatIOCTL(result.Code), result.BytesReturned)
			}
		}
	})

	t.Run("ProbeSpecificIOCTLs", func(t *testing.T) {
		var hDevice handle.HANDLE
		clfsDevicePaths := GetDriverDevicePaths("CLFS")
		clfsDevicePaths = append(clfsDevicePaths, `\\.\CLFS`, `\\.\ClfsLog`)

		for _, path := range clfsDevicePaths {
			h, err := OpenDevice(path, GENERIC_READ|GENERIC_WRITE)
			if err == nil {
				hDevice = h
				break
			}
		}

		if hDevice == 0 {
			t.Skip("CLFS device not accessible")
			return
		}
		defer CloseHandle(hDevice)

		// Probe specific IOCTLs to see which are recognized
		const FILE_DEVICE_FILE_SYSTEM = 0x00000009
		testIOCTLs := []struct {
			name string
			code uint32
		}{
			{"Function 0 BUFFERED", CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 0, METHOD_BUFFERED, FILE_ANY_ACCESS)},
			{"Function 1 BUFFERED", CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 1, METHOD_BUFFERED, FILE_ANY_ACCESS)},
			{"Function 10 BUFFERED", CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 10, METHOD_BUFFERED, FILE_ANY_ACCESS)},
			{"Function 42 BUFFERED", CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 42, METHOD_BUFFERED, FILE_ANY_ACCESS)},
			{"Function 100 BUFFERED", CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 100, METHOD_BUFFERED, FILE_ANY_ACCESS)},
		}

		t.Log("Probing specific IOCTL codes...")
		for _, tc := range testIOCTLs {
			result := ProbeIOCTL(hDevice, tc.code)
			t.Logf("  %-25s: 0x%08X - valid=%v, bytes=%d, err=%v",
				tc.name, tc.code, result.Valid, result.BytesReturned, result.ErrorCode)
		}
	})

	t.Run("SendTestIOCTLs", func(t *testing.T) {
		var hDevice handle.HANDLE
		clfsDevicePaths := GetDriverDevicePaths("CLFS")
		clfsDevicePaths = append(clfsDevicePaths, `\\.\CLFS`, `\\.\ClfsLog`)

		for _, path := range clfsDevicePaths {
			h, err := OpenDevice(path, GENERIC_READ|GENERIC_WRITE)
			if err == nil {
				hDevice = h
				t.Logf("Opened CLFS device: %s", path)
				break
			}
		}

		if hDevice == 0 {
			t.Skip("CLFS device not accessible")
			return
		}
		defer CloseHandle(hDevice)

		// Send a test IOCTL with varying buffer sizes
		const FILE_DEVICE_FILE_SYSTEM = 0x00000009
		testCode := CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 0, METHOD_BUFFERED, FILE_ANY_ACCESS)

		bufferSizes := []int{16, 64, 256, 1024, 4096}
		t.Logf("Testing IOCTL 0x%08X with various buffer sizes...", testCode)

		for _, size := range bufferSizes {
			inBuf := make([]byte, size)
			outBuf := make([]byte, size)
			var bytesReturned uint32

			success, err := DeviceIoControl(
				hDevice,
				testCode,
				unsafe.Pointer(&inBuf[0]),
				uint32(len(inBuf)),
				unsafe.Pointer(&outBuf[0]),
				uint32(len(outBuf)),
				&bytesReturned,
				nil,
			)

			t.Logf("  Buffer size %5d: success=%v, bytesReturned=%d, err=%v",
				size, success, bytesReturned, err)
		}
	})

	// Comprehensive test: Use system CLFS or PhysicalDrive for IOCTL testing
	t.Run("LoadDriverAndTestIOCTLs", func(t *testing.T) {
		// CLFS is typically already loaded as a Windows system driver
		// Instead of trying to load it separately, let's use PhysicalDrive0 for IOCTL testing
		// which is always available and supports standard disk IOCTLs

		var hDevice handle.HANDLE
		var devicePath string

		// Try PhysicalDrive0 which is always available for IOCTL testing
		testDevices := []string{
			`\\.\PhysicalDrive0`,
			`\\.\C:`,
			`\\.\GLOBALROOT\Device\HarddiskVolume1`,
		}

		t.Log("Testing IOCTL discovery and fuzzing capabilities...")
		t.Logf("Trying %d device paths for testing...", len(testDevices))

		for _, path := range testDevices {
			// Use GENERIC_READ for safer testing
			h, err := OpenDevice(path, GENERIC_READ)
			if err == nil {
				hDevice = h
				devicePath = path
				t.Logf("[*] Successfully opened device: %s", path)
				break
			}
			t.Logf("  ✗ Failed to open %s: %v", path, err)
		}

		if hDevice == 0 {
			// Last resort: Try CLFS paths (may already be loaded by Windows)
			clfsDevicePaths := []string{
				`\\.\CLFS`,
				`\\.\ClfsLog`,
				`\\.\Global\CLFS`,
				`\\Device\CLFS`,
				`\\.\GLOBALROOT\Device\CLFS`,
			}

			discoveredPaths := GetDriverDevicePaths("CLFS")
			clfsDevicePaths = append(clfsDevicePaths, discoveredPaths...)

			t.Logf("Trying %d CLFS device paths...", len(clfsDevicePaths))
			for _, path := range clfsDevicePaths {
				h, err := OpenDevice(path, GENERIC_READ|GENERIC_WRITE)
				if err == nil {
					hDevice = h
					devicePath = path
					t.Logf("[*] Successfully opened CLFS device: %s", path)
					break
				}
			}
		}

		if hDevice == 0 {
			t.Skip("Could not open any device for IOCTL testing (may need admin privileges)")
			return
		}
		defer CloseHandle(hDevice)

		// Now perform IOCTL discovery
		t.Run("DiscoverIOCTLs", func(t *testing.T) {
			const FILE_DEVICE_FILE_SYSTEM = 0x00000009
			t.Logf("Discovering IOCTLs for device: %s", devicePath)
			t.Logf("Scanning device type: FILE_DEVICE_FILE_SYSTEM (0x%02X)", FILE_DEVICE_FILE_SYSTEM)

			results := DiscoverIOCTLsByDeviceType(hDevice, FILE_DEVICE_FILE_SYSTEM)
			t.Logf("[*] Found %d valid IOCTL codes", len(results))

			if len(results) > 0 {
				t.Log("\nValid IOCTLs discovered:")
				for i, result := range results {
					// if i >= 20 {
					// 	t.Logf("  ... and %d more IOCTLs", len(results)-20)
					// 	break
					// }
					components := DecodeIOCTL(result.Code)
					t.Logf("  [%2d] 0x%08X - %s", i+1, result.Code, FormatIOCTL(result.Code))
					t.Logf("       Function: 0x%03X, Method: %s, Access: %s, BytesReturned: %d",
						components.Function, components.MethodName, components.AccessName, result.BytesReturned)
					if result.ErrorCode != nil {
						t.Logf("       Error: %v", result.ErrorCode)
					}
				}
			} else {
				t.Log("No valid IOCTLs discovered (this might indicate the device doesn't respond to our probes)")
			}
		})

		// Test specific IOCTL codes
		t.Run("ProbeSpecificIOCTLs", func(t *testing.T) {
			const FILE_DEVICE_FILE_SYSTEM = 0x00000009

			testCases := []struct {
				function uint32
				method   uint32
				access   uint32
			}{
				{0x000, METHOD_BUFFERED, FILE_ANY_ACCESS},
				{0x001, METHOD_BUFFERED, FILE_ANY_ACCESS},
				{0x002, METHOD_BUFFERED, FILE_READ_ACCESS},
				{0x003, METHOD_BUFFERED, FILE_WRITE_ACCESS},
				{0x010, METHOD_BUFFERED, FILE_ANY_ACCESS},
				{0x020, METHOD_IN_DIRECT, FILE_ANY_ACCESS},
				{0x030, METHOD_OUT_DIRECT, FILE_ANY_ACCESS},
			}

			t.Logf("Probing %d specific IOCTL codes...", len(testCases))
			validCount := 0
			for _, tc := range testCases {
				code := CTL_CODE(FILE_DEVICE_FILE_SYSTEM, tc.function, tc.method, tc.access)
				result := ProbeIOCTL(hDevice, code)

				if result.Valid {
					validCount++
					t.Logf("[*] 0x%08X - VALID: %s", code, FormatIOCTL(code))
					if result.ErrorCode != nil {
						t.Logf("[X] Error: %v, BytesReturned: %d", result.ErrorCode, result.BytesReturned)
					}
				}
			}
			t.Logf("Result: %d/%d IOCTLs recognized as valid", validCount, len(testCases))
		})

		// Test with varying buffer sizes
		t.Run("TestBufferSizes", func(t *testing.T) {
			const FILE_DEVICE_FILE_SYSTEM = 0x00000009
			testCode := CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 0, METHOD_BUFFERED, FILE_ANY_ACCESS)

			bufferSizes := []int{0, 16, 64, 256, 1024, 4096}
			t.Logf("Testing IOCTL 0x%08X with varying buffer sizes...", testCode)

			for _, size := range bufferSizes {
				var inBuf, outBuf []byte
				if size > 0 {
					inBuf = make([]byte, size)
					outBuf = make([]byte, size)
				}

				var bytesReturned uint32
				var inPtr, outPtr unsafe.Pointer

				if size > 0 {
					inPtr = unsafe.Pointer(&inBuf[0])
					outPtr = unsafe.Pointer(&outBuf[0])
				}

				success, err := DeviceIoControl(
					hDevice, testCode,
					inPtr, uint32(size),
					outPtr, uint32(size),
					&bytesReturned, nil,
				)

				t.Logf("  Buffer size %5d: success=%v, bytesReturned=%d, err=%v",
					size, success, bytesReturned, err)
			}
		})
	})
}

// TestAfdDriver tests interactions with the Ancillary Function Driver for WinSock (afd.sys)
func TestAfdDriver(t *testing.T) {
	t.Run("DiscoverAfdDevice", func(t *testing.T) {
		t.Log("Searching for AFD driver...")

		// First, search for symbolic links with "afd" in the name
		t.Log("Searching DOS device symbolic links...")
		symlinks, err := FindSymbolicLinksByPattern("afd")
		if err != nil {
			t.Logf("Error searching symbolic links: %v", err)
		} else if len(symlinks) > 0 {
			t.Logf("Found %d AFD symbolic link(s):", len(symlinks))
			for device, targets := range symlinks {
				t.Logf("  %s ->", device)
				for _, target := range targets {
					t.Logf("    %s", target)
				}
			}
		} else {
			t.Log("No AFD symbolic links found")
		}

		// Also search for devices with service name "AFD"
		t.Log("Searching for AFD devices by service name...")
		devices, err := FindDevicesByService("AFD")
		if err != nil {
			t.Logf("Error searching for AFD devices: %v", err)
		}

		if len(devices) == 0 {
			t.Log("No AFD devices found through service name search")
		} else {
			t.Logf("Found %d AFD device(s):", len(devices))
			for i, device := range devices {
				t.Logf("  Device %d:", i)
				if device.Description != "" {
					t.Logf("    Description: %s", device.Description)
				}
				if device.Service != "" {
					t.Logf("    Service: %s", device.Service)
				}
				if device.Class != "" {
					t.Logf("    Class: %s", device.Class)
				}
				if device.FriendlyName != "" {
					t.Logf("    Friendly Name: %s", device.FriendlyName)
				}
			}
		}
	})

	t.Run("OpenAfdDevice", func(t *testing.T) {
		var workingPath string
		var hDevice handle.HANDLE
		var err error

		// Get common AFD device paths
		afdDevicePaths := GetDriverDevicePaths("AFD")
		// Add some known AFD-specific paths
		afdDevicePaths = append(afdDevicePaths, `\\Device\Afd\Endpoint`)

		// Try each known path
		for _, devicePath := range afdDevicePaths {
			t.Logf("Trying device path: %s", devicePath)

			hDevice, err = CreateFile(
				devicePath,
				GENERIC_READ|GENERIC_WRITE,
				FILE_SHARE_READ|FILE_SHARE_WRITE,
				nil,
				OPEN_EXISTING,
				0,
				0,
			)

			if err == nil && hDevice != 0 && uintptr(hDevice) != INVALID_HANDLE_VALUE {
				workingPath = devicePath
				t.Logf("Successfully opened AFD device at %s with handle: 0x%x", devicePath, hDevice)
				break
			}

			t.Logf("  Failed: %v", err)
		}

		if workingPath == "" {
			t.Log("Could not open any AFD device path (expected - AFD requires special access)")
			return
		}
		defer CloseHandle(hDevice)
	})

	t.Run("AfdDeviceIoControl", func(t *testing.T) {
		var hDevice handle.HANDLE
		var err error

		// Get common AFD device paths
		afdDevicePaths := GetDriverDevicePaths("AFD")
		afdDevicePaths = append(afdDevicePaths, `\\Device\Afd\Endpoint`)

		// Find a working device path
		for _, devicePath := range afdDevicePaths {
			hDevice, err = OpenDeviceReadWrite(devicePath)
			if err == nil && hDevice != 0 && uintptr(hDevice) != INVALID_HANDLE_VALUE {
				break
			}
		}

		if hDevice == 0 || uintptr(hDevice) == INVALID_HANDLE_VALUE {
			t.Skip("AFD device not accessible (expected)")
			return
		}
		defer CloseHandle(hDevice)

		// AFD IOCTL codes (FILE_DEVICE_NETWORK = 0x00000012)
		const FILE_DEVICE_NETWORK = 0x00000012

		ioctlCodes := []struct {
			name string
			code uint32
		}{
			// AFD IOCTLs - these typically require socket context
			{"AFD_BIND", CTL_CODE(FILE_DEVICE_NETWORK, 0x00, METHOD_BUFFERED, FILE_ANY_ACCESS)},
			{"AFD_CONNECT", CTL_CODE(FILE_DEVICE_NETWORK, 0x01, METHOD_BUFFERED, FILE_ANY_ACCESS)},
			{"AFD_START_LISTEN", CTL_CODE(FILE_DEVICE_NETWORK, 0x02, METHOD_BUFFERED, FILE_ANY_ACCESS)},
			{"AFD_RECV", CTL_CODE(FILE_DEVICE_NETWORK, 0x05, METHOD_BUFFERED, FILE_ANY_ACCESS)},
			{"AFD_SEND", CTL_CODE(FILE_DEVICE_NETWORK, 0x06, METHOD_BUFFERED, FILE_ANY_ACCESS)},
		}

		for _, ioctl := range ioctlCodes {
			t.Run(ioctl.name, func(t *testing.T) {
				var bytesReturned uint32
				outBuffer := make([]byte, 256)

				success, err := DeviceIoControl(
					hDevice,
					ioctl.code,
					nil,
					0,
					unsafe.Pointer(&outBuffer[0]),
					uint32(len(outBuffer)),
					&bytesReturned,
					nil,
				)

				if success {
					t.Logf("IOCTL %s (0x%08X) succeeded, returned %d bytes", ioctl.name, ioctl.code, bytesReturned)
					if bytesReturned > 0 {
						displayLen := bytesReturned
						if displayLen > 16 {
							displayLen = 16
						}
						t.Logf("  First bytes: % X", outBuffer[:displayLen])
					}
				} else {
					t.Logf("IOCTL %s (0x%08X) failed: %v (expected without socket context)", ioctl.name, ioctl.code, err)
				}
			})
		}
	})
}

// TestPhysicalDriveAccess tests accessing physical drives
func TestPhysicalDriveAccess(t *testing.T) {
	t.Run("OpenPhysicalDrive0", func(t *testing.T) {
		hDevice, err := CreateFile(
			`\\.\PhysicalDrive0`,
			GENERIC_READ,
			FILE_SHARE_READ|FILE_SHARE_WRITE,
			nil,
			OPEN_EXISTING,
			0,
			0,
		)

		if err != nil || hDevice == 0 || uintptr(hDevice) == INVALID_HANDLE_VALUE {
			t.Logf("Failed to open PhysicalDrive0 (requires admin): %v", err)
			t.Skip("Skipping physical drive tests - requires admin privileges")
			return
		}
		defer CloseHandle(hDevice)

		t.Logf("Successfully opened PhysicalDrive0 with handle: 0x%x", hDevice)
	})

	t.Run("GetDriveGeometry", func(t *testing.T) {
		hDevice, err := OpenDeviceReadOnly(`\\.\PhysicalDrive0`)
		if err != nil || hDevice == 0 || uintptr(hDevice) == INVALID_HANDLE_VALUE {
			t.Skip("PhysicalDrive0 not accessible")
			return
		}
		defer CloseHandle(hDevice)

		var geometry DISK_GEOMETRY
		var bytesReturned uint32

		success, err := DeviceIoControl(
			hDevice,
			IOCTL_DISK_GET_DRIVE_GEOMETRY,
			nil,
			0,
			unsafe.Pointer(&geometry),
			uint32(unsafe.Sizeof(geometry)),
			&bytesReturned,
			nil,
		)

		if !success {
			t.Fatalf("IOCTL_DISK_GET_DRIVE_GEOMETRY failed: %v", err)
		}

		t.Logf("Disk Geometry:")
		t.Logf("  Cylinders: %d", geometry.Cylinders)
		t.Logf("  Media Type: 0x%X", geometry.MediaType)
		t.Logf("  Tracks per Cylinder: %d", geometry.TracksPerCylinder)
		t.Logf("  Sectors per Track: %d", geometry.SectorsPerTrack)
		t.Logf("  Bytes per Sector: %d", geometry.BytesPerSector)

		// Calculate total size
		totalSize := uint64(geometry.Cylinders) *
			uint64(geometry.TracksPerCylinder) *
			uint64(geometry.SectorsPerTrack) *
			uint64(geometry.BytesPerSector)

		t.Logf("  Total Size: %d bytes (%.2f GB)", totalSize, float64(totalSize)/(1024*1024*1024))
	})

	t.Run("GetDeviceNumber", func(t *testing.T) {
		hDevice, err := OpenDeviceReadOnly(`\\.\PhysicalDrive0`)
		if err != nil || hDevice == 0 || uintptr(hDevice) == INVALID_HANDLE_VALUE {
			t.Skip("PhysicalDrive0 not accessible")
			return
		}
		defer CloseHandle(hDevice)

		var deviceNumber STORAGE_DEVICE_NUMBER
		var bytesReturned uint32

		success, err := DeviceIoControl(
			hDevice,
			IOCTL_STORAGE_GET_DEVICE_NUMBER,
			nil,
			0,
			unsafe.Pointer(&deviceNumber),
			uint32(unsafe.Sizeof(deviceNumber)),
			&bytesReturned,
			nil,
		)

		if !success {
			t.Fatalf("IOCTL_STORAGE_GET_DEVICE_NUMBER failed: %v", err)
		}

		t.Logf("Storage Device Number:")
		t.Logf("  Device Type: 0x%X", deviceNumber.DeviceType)
		t.Logf("  Device Number: %d", deviceNumber.DeviceNumber)
		t.Logf("  Partition Number: %d", deviceNumber.PartitionNumber)
	})

	t.Run("QueryStorageProperty", func(t *testing.T) {
		hDevice, err := OpenDeviceReadOnly(`\\.\PhysicalDrive0`)
		if err != nil || hDevice == 0 || uintptr(hDevice) == INVALID_HANDLE_VALUE {
			t.Skip("PhysicalDrive0 not accessible")
			return
		}
		defer CloseHandle(hDevice)

		// Prepare query structure
		query := STORAGE_PROPERTY_QUERY{
			PropertyId: StorageDeviceProperty,
			QueryType:  PropertyStandardQuery,
		}

		// Allocate output buffer
		outBuffer := make([]byte, 4096)
		var bytesReturned uint32

		success, err := DeviceIoControl(
			hDevice,
			IOCTL_STORAGE_QUERY_PROPERTY,
			unsafe.Pointer(&query),
			uint32(unsafe.Sizeof(query)),
			unsafe.Pointer(&outBuffer[0]),
			uint32(len(outBuffer)),
			&bytesReturned,
			nil,
		)

		if !success {
			t.Fatalf("IOCTL_STORAGE_QUERY_PROPERTY failed: %v", err)
		}

		// Parse descriptor
		descriptor := (*STORAGE_DEVICE_DESCRIPTOR)(unsafe.Pointer(&outBuffer[0]))
		t.Logf("Storage Device Descriptor:")
		t.Logf("  Version: %d", descriptor.Version)
		t.Logf("  Size: %d", descriptor.Size)
		t.Logf("  Device Type: 0x%X", descriptor.DeviceType)
		t.Logf("  Removable Media: %d", descriptor.RemovableMedia)
		t.Logf("  Command Queueing: %d", descriptor.CommandQueueing)
		t.Logf("  Bus Type: %d", descriptor.BusType)

		// Extract vendor and product strings if available
		if descriptor.VendorIdOffset > 0 && descriptor.VendorIdOffset < uint32(len(outBuffer)) {
			vendorBytes := outBuffer[descriptor.VendorIdOffset:]
			vendorStr := string(vendorBytes[:findNull(vendorBytes)])
			t.Logf("  Vendor ID: %s", vendorStr)
		}

		if descriptor.ProductIdOffset > 0 && descriptor.ProductIdOffset < uint32(len(outBuffer)) {
			productBytes := outBuffer[descriptor.ProductIdOffset:]
			productStr := string(productBytes[:findNull(productBytes)])
			t.Logf("  Product ID: %s", productStr)
		}
	})
}

// TestVolumeAccess tests accessing logical volumes
func TestVolumeAccess(t *testing.T) {
	t.Run("OpenVolumeC", func(t *testing.T) {
		hDevice, err := CreateFile(
			`\\.\C:`,
			GENERIC_READ,
			FILE_SHARE_READ|FILE_SHARE_WRITE,
			nil,
			OPEN_EXISTING,
			0,
			0,
		)

		if err != nil || hDevice == 0 || uintptr(hDevice) == INVALID_HANDLE_VALUE {
			t.Logf("Failed to open C: volume (requires admin): %v", err)
			t.Skip("Skipping volume tests")
			return
		}
		defer CloseHandle(hDevice)

		t.Logf("Successfully opened C: volume with handle: 0x%x", hDevice)
	})

	t.Run("GetVolumeDiskExtents", func(t *testing.T) {
		hDevice, err := OpenDeviceReadOnly(`\\.\C:`)
		if err != nil || hDevice == 0 || uintptr(hDevice) == INVALID_HANDLE_VALUE {
			t.Skip("C: volume not accessible")
			return
		}
		defer CloseHandle(hDevice)

		// Allocate buffer for extents
		buffer := make([]byte, unsafe.Sizeof(VOLUME_DISK_EXTENTS{})+10*unsafe.Sizeof(DISK_EXTENT{}))
		var bytesReturned uint32

		success, err := DeviceIoControl(
			hDevice,
			IOCTL_VOLUME_GET_VOLUME_DISK_EXTENTS,
			nil,
			0,
			unsafe.Pointer(&buffer[0]),
			uint32(len(buffer)),
			&bytesReturned,
			nil,
		)

		if !success {
			t.Logf("IOCTL_VOLUME_GET_VOLUME_DISK_EXTENTS failed: %v", err)
			return
		}

		extents := (*VOLUME_DISK_EXTENTS)(unsafe.Pointer(&buffer[0]))
		t.Logf("Volume Disk Extents:")
		t.Logf("  Number of Disk Extents: %d", extents.NumberOfDiskExtents)

		for i := uint32(0); i < extents.NumberOfDiskExtents && i < 10; i++ {
			extent := (*DISK_EXTENT)(unsafe.Pointer(
				uintptr(unsafe.Pointer(&extents.Extents[0])) +
					uintptr(i)*unsafe.Sizeof(DISK_EXTENT{}),
			))
			t.Logf("  Extent %d:", i)
			t.Logf("    Disk Number: %d", extent.DiskNumber)
			t.Logf("    Starting Offset: %d", extent.StartingOffset)
			t.Logf("    Extent Length: %d bytes", extent.ExtentLength)
		}
	})
}

// TestDeviceEnumeration tests device enumeration using SetupDi
func TestDeviceEnumeration(t *testing.T) {
	t.Run("EnumerateDiskDevices", func(t *testing.T) {
		// GUID_DEVINTERFACE_DISK
		diskGUID := GUID{
			Data1: 0x53f56307,
			Data2: 0xb6bf,
			Data3: 0x11d0,
			Data4: [8]byte{0x94, 0xf2, 0x00, 0xa0, 0xc9, 0x1e, 0xfb, 0x8b},
		}

		devices, err := EnumerateDevices(&diskGUID, DIGCF_PRESENT|DIGCF_DEVICEINTERFACE)
		if err != nil {
			t.Fatalf("Failed to enumerate disk devices: %v", err)
		}

		t.Logf("Found %d disk device(s):", len(devices))
		for i, device := range devices {
			t.Logf("  [%d] %s", i, device)
		}
	})

	t.Run("EnumerateAllDevices", func(t *testing.T) {
		// Enumerate all present devices (may return many results)
		deviceInfoSet, err := SetupDiGetClassDevs(nil, "", 0, DIGCF_PRESENT|DIGCF_ALLCLASSES)
		if err != nil {
			t.Fatalf("SetupDiGetClassDevs failed: %v", err)
		}
		defer SetupDiDestroyDeviceInfoList(deviceInfoSet)

		var deviceInfoData SP_DEVINFO_DATA
		index := uint32(0)
		deviceCount := 0

		for {
			success, err := SetupDiEnumDeviceInfo(deviceInfoSet, index, &deviceInfoData)
			if !success {
				if err != nil {
					t.Logf("Error enumerating devices: %v", err)
				}
				break
			}

			// Get device description
			desc, _ := SetupDiGetDeviceRegistryProperty(
				deviceInfoSet,
				&deviceInfoData,
				SPDRP_DEVICEDESC,
			)

			if deviceCount < 10000 { // Only log first 10000 to avoid spam
				t.Logf("Device %d: %s", index, desc)
			}

			deviceCount++
			index++
		}

		t.Logf("Total devices enumerated: %d", deviceCount)
	})
}

// TestIOCTLCodeConstruction tests building various IOCTL codes
func TestIOCTLCodeConstruction(t *testing.T) {
	tests := []struct {
		name       string
		deviceType uint32
		function   uint32
		method     uint32
		access     uint32
		expected   string
	}{
		{
			name:       "Disk Get Geometry",
			deviceType: FILE_DEVICE_DISK,
			function:   0,
			method:     METHOD_BUFFERED,
			access:     FILE_ANY_ACCESS,
			expected:   fmt.Sprintf("0x%08X", IOCTL_DISK_GET_DRIVE_GEOMETRY),
		},
		{
			name:       "Custom Device IOCTL",
			deviceType: 0x8000,
			function:   0x800,
			method:     METHOD_BUFFERED,
			access:     FILE_ANY_ACCESS,
			expected:   "0x80002000",
		},
		{
			name:       "Network Device IOCTL",
			deviceType: FILE_DEVICE_NETWORK,
			function:   0x100,
			method:     METHOD_NEITHER,
			access:     FILE_READ_ACCESS | FILE_WRITE_ACCESS,
			expected:   fmt.Sprintf("0x%08X", CTL_CODE(FILE_DEVICE_NETWORK, 0x100, METHOD_NEITHER, FILE_READ_ACCESS|FILE_WRITE_ACCESS)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := CTL_CODE(tt.deviceType, tt.function, tt.method, tt.access)
			t.Logf("IOCTL Code: 0x%08X (expected %s)", code, tt.expected)

			// Verify components can be extracted
			extractedDeviceType := (code >> 16) & 0xFFFF
			extractedMethod := code & 0x3

			if extractedDeviceType != tt.deviceType {
				t.Errorf("Device type mismatch: got 0x%X, want 0x%X", extractedDeviceType, tt.deviceType)
			}

			if extractedMethod != tt.method {
				t.Errorf("Method mismatch: got 0x%X, want 0x%X", extractedMethod, tt.method)
			}
		})
	}
}

// Helper function to find null terminator in byte slice
func findNull(data []byte) int {
	for i, b := range data {
		if b == 0 {
			return i
		}
	}
	return len(data)
}

// BenchmarkPhysicalDriveOpen benchmarks opening a physical drive
func BenchmarkPhysicalDriveOpen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hDevice, err := OpenDeviceReadOnly(`\\.\PhysicalDrive0`)
		if err != nil {
			b.Skip("PhysicalDrive0 not accessible")
			return
		}
		CloseHandle(hDevice)
	}
}

// BenchmarkDiskGeometryIOCTL benchmarks getting disk geometry
func BenchmarkDiskGeometryIOCTL(b *testing.B) {
	hDevice, err := OpenDeviceReadOnly(`\\.\PhysicalDrive0`)
	if err != nil {
		b.Skip("PhysicalDrive0 not accessible")
		return
	}
	defer CloseHandle(hDevice)

	var geometry DISK_GEOMETRY
	var bytesReturned uint32

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeviceIoControl(
			hDevice,
			IOCTL_DISK_GET_DRIVE_GEOMETRY,
			nil,
			0,
			unsafe.Pointer(&geometry),
			uint32(unsafe.Sizeof(geometry)),
			&bytesReturned,
			nil,
		)
	}
}

// ============================================================================
// Driver Loading and Management Tests
// ============================================================================

// TestDriverLoadBasic tests basic driver loading with default options
func TestDriverLoadBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping driver loading test in short mode")
	}

	driverPath := `C:\Windows\System32\drivers\null.sys`
	driverName := "NullDriver_BasicTest"

	hService, err := LoadDriver(driverPath, driverName)
	if err != nil {
		t.Skipf("Cannot load driver (need admin): %v", err)
		return
	}
	defer UnloadDriver(hService)

	t.Log("[*] LoadDriver works with default options")

	// Verify driver is running
	status, err := QueryDriverStatus(hService)
	if err != nil {
		t.Errorf("Failed to query driver status: %v", err)
	} else {
		t.Logf("Driver state: %d (4=RUNNING)", status.CurrentState)
	}
}

// TestDriverLoadEx tests extended driver loading with custom parameters
func TestDriverLoadEx(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping driver loading test in short mode")
	}

	driverPath := `C:\Windows\System32\drivers\null.sys`
	driverName := "NullDriver_ExTest"

	hService, err := LoadDriverEx(
		driverPath,
		driverName,
		service.SERVICE_ALL_ACCESS,
		service.SERVICE_DEMAND_START,
		service.SERVICE_ERROR_IGNORE, // Ignore errors for testing
	)
	if err != nil {
		t.Skipf("Cannot load driver (need admin): %v", err)
		return
	}
	defer UnloadDriver(hService)

	t.Log("[*] LoadDriverEx works with custom parameters")
}

// TestDriverLoadWithOptions tests driver loading with detailed options
func TestDriverLoadWithOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping driver loading test in short mode")
	}

	driverPath := `C:\Windows\System32\drivers\null.sys`
	driverName := "NullDriver_OptionsTest"

	options := DriverLoadOptions{
		DesiredAccess:    service.SERVICE_ALL_ACCESS,
		StartType:        service.SERVICE_DEMAND_START,
		ErrorControl:     service.SERVICE_ERROR_NORMAL,
		StartImmediately: false, // Don't start yet
		RecreateIfExists: true,  // Force recreation
	}

	hService, err := LoadDriverWithOptions(driverPath, driverName, options)
	if err != nil {
		t.Skipf("Cannot load driver (need admin): %v", err)
		return
	}
	defer UnloadDriver(hService)

	t.Log("[*] LoadDriverWithOptions works with custom options")

	// Manually start the driver
	if err := StartDriver(hService); err != nil {
		t.Errorf("Failed to start driver: %v", err)
	} else {
		t.Log("[*] StartDriver works")
	}

	// Query status
	status, err := QueryDriverStatus(hService)
	if err != nil {
		t.Errorf("Failed to query status: %v", err)
	} else {
		t.Logf("[*] QueryDriverStatus works: state=%d", status.CurrentState)
	}

	// Stop the driver
	if err := StopDriver(hService); err != nil {
		t.Errorf("Failed to stop driver: %v", err)
	} else {
		t.Log("[*] StopDriver works")
	}
}

// TestStartStopDriver tests starting and stopping a driver multiple times
func TestStartStopDriver(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping driver start/stop test in short mode")
	}

	driverPath := `C:\Windows\System32\drivers\null.sys`
	driverName := "NullDriver_StartStopTest"

	// Load driver without starting
	options := DefaultDriverLoadOptions()
	options.StartImmediately = false

	hService, err := LoadDriverWithOptions(driverPath, driverName, options)
	if err != nil {
		t.Skipf("Cannot load driver (need admin): %v", err)
		return
	}
	defer UnloadDriver(hService)

	// Test start/stop cycles
	for i := 0; i < 3; i++ {
		// Start the driver
		if err := StartDriver(hService); err != nil {
			t.Errorf("Cycle %d: Failed to start driver: %v", i+1, err)
			return
		}
		t.Logf("Cycle %d: Driver started", i+1)

		// Verify it's running
		status, err := QueryDriverStatus(hService)
		if err == nil && status.CurrentState != service.SERVICE_RUNNING {
			t.Errorf("Cycle %d: Driver not in RUNNING state: %d", i+1, status.CurrentState)
		}

		// Stop the driver
		if err := StopDriver(hService); err != nil {
			t.Errorf("Cycle %d: Failed to stop driver: %v", i+1, err)
			return
		}
		t.Logf("Cycle %d: Driver stopped", i+1)
	}

	t.Log("[*] Start/stop cycles completed successfully")
}

// TestUnloadDriverEx tests unloading with custom options
func TestUnloadDriverEx(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping driver unload test in short mode")
	}

	driverPath := `C:\Windows\System32\drivers\null.sys`
	driverName := "NullDriver_UnloadExTest"

	hService, err := LoadDriver(driverPath, driverName)
	if err != nil {
		t.Skipf("Cannot load driver (need admin): %v", err)
		return
	}

	// Stop driver but keep service registered (don't delete, don't close handle)
	if err := UnloadDriverEx(hService, false, false); err != nil {
		t.Errorf("Failed to stop driver: %v", err)
		return
	}
	t.Log("[*] Driver stopped (service still registered)")

	// Verify it's stopped
	status, err := QueryDriverStatus(hService)
	if err == nil && status.CurrentState != service.SERVICE_STOPPED {
		t.Logf("Warning: Driver state is %d, expected STOPPED (1)", status.CurrentState)
	}

	// Completely unload it (delete service, close handle)
	if err := UnloadDriverEx(hService, true, true); err != nil {
		t.Errorf("Failed to unload driver: %v", err)
		return
	}
	t.Log("[*] Driver completely unloaded and service deleted")
}

// TestOpenExistingDriver tests opening a handle to an existing system driver
func TestOpenExistingDriver(t *testing.T) {
	// Try to open Beep driver which is usually present on Windows
	hService, err := OpenExistingDriver("Beep", service.SERVICE_QUERY_STATUS)
	if err != nil {
		t.Skipf("Cannot open Beep driver: %v", err)
		return
	}
	defer service.CloseServiceHandle(hService)

	t.Log("[*] OpenExistingDriver successfully opened Beep driver")

	// Query its status
	status, err := QueryDriverStatus(hService)
	if err != nil {
		t.Errorf("Failed to query Beep driver status: %v", err)
		return
	}

	t.Logf("[*] Beep driver status: state=%d, type=%d", status.CurrentState, status.ServiceType)
}

// TestQueryDriverStatus tests querying driver status information
func TestQueryDriverStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping driver status query test in short mode")
	}

	driverPath := `C:\Windows\System32\drivers\null.sys`
	driverName := "NullDriver_StatusTest"

	hService, err := LoadDriver(driverPath, driverName)
	if err != nil {
		t.Skipf("Cannot load driver (need admin): %v", err)
		return
	}
	defer UnloadDriver(hService)

	// Query status
	status, err := QueryDriverStatus(hService)
	if err != nil {
		t.Errorf("Failed to query driver status: %v", err)
		return
	}

	t.Logf("Driver Status:")
	t.Logf("  Service Type: %d", status.ServiceType)
	t.Logf("  Current State: %d (1=STOPPED, 2=START_PENDING, 3=STOP_PENDING, 4=RUNNING)", status.CurrentState)
	t.Logf("  Controls Accepted: 0x%X", status.ControlsAccepted)
	t.Logf("  Win32 Exit Code: %d", status.Win32ExitCode)
	t.Logf("  Service Specific Exit Code: %d", status.ServiceSpecificExitCode)

	// Verify it's a kernel driver
	if status.ServiceType != service.SERVICE_KERNEL_DRIVER {
		t.Errorf("Expected SERVICE_KERNEL_DRIVER (%d), got %d", service.SERVICE_KERNEL_DRIVER, status.ServiceType)
	}

	// Verify it's running
	if status.CurrentState != service.SERVICE_RUNNING {
		t.Errorf("Expected SERVICE_RUNNING (4), got %d", status.CurrentState)
	}
}

// TestDefaultDriverLoadOptions tests the default options function
func TestDefaultDriverLoadOptions(t *testing.T) {
	options := DefaultDriverLoadOptions()

	if options.DesiredAccess != service.SERVICE_ALL_ACCESS {
		t.Errorf("Expected DesiredAccess=%d, got %d", service.SERVICE_ALL_ACCESS, options.DesiredAccess)
	}

	if options.StartType != service.SERVICE_DEMAND_START {
		t.Errorf("Expected StartType=%d, got %d", service.SERVICE_DEMAND_START, options.StartType)
	}

	if options.ErrorControl != service.SERVICE_ERROR_NORMAL {
		t.Errorf("Expected ErrorControl=%d, got %d", service.SERVICE_ERROR_NORMAL, options.ErrorControl)
	}

	if !options.StartImmediately {
		t.Error("Expected StartImmediately=true, got false")
	}

	if options.RecreateIfExists {
		t.Error("Expected RecreateIfExists=false, got true")
	}

	t.Log("[*] DefaultDriverLoadOptions returns correct defaults")
}

// TestDriverRecreateIfExists tests the RecreateIfExists option
func TestDriverRecreateIfExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping driver recreation test in short mode")
	}

	driverPath := `C:\Windows\System32\drivers\null.sys`
	driverName := "NullDriver_RecreateTest"

	// First load
	hService1, err := LoadDriver(driverPath, driverName)
	if err != nil {
		t.Skipf("Cannot load driver (need admin): %v", err)
		return
	}
	UnloadDriver(hService1)

	// Second load with RecreateIfExists=false (should reuse)
	options := DefaultDriverLoadOptions()
	options.RecreateIfExists = false

	hService2, err := LoadDriverWithOptions(driverPath, driverName, options)
	if err != nil {
		t.Errorf("Failed to load existing driver: %v", err)
		return
	}
	t.Log("[*] Successfully reused existing driver service")
	UnloadDriver(hService2)

	// Third load with RecreateIfExists=true (should recreate)
	options.RecreateIfExists = true
	hService3, err := LoadDriverWithOptions(driverPath, driverName, options)
	if err != nil {
		t.Errorf("Failed to recreate driver: %v", err)
		return
	}
	t.Log("[*] Successfully recreated driver service")
	UnloadDriver(hService3)
}

// TestDriverVariations demonstrates all driver loading variations
func TestDriverVariations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping driver variations test in short mode")
	}

	t.Run("BasicLoad", func(t *testing.T) {
		driverPath := `C:\Windows\System32\drivers\null.sys`
		driverName := "NullDriver_Var1"

		hService, err := LoadDriver(driverPath, driverName)
		if err != nil {
			t.Skipf("Cannot load driver (need admin): %v", err)
			return
		}
		defer UnloadDriver(hService)

		t.Log("[*] LoadDriver works")
	})

	t.Run("ExtendedLoad", func(t *testing.T) {
		driverPath := `C:\Windows\System32\drivers\null.sys`
		driverName := "NullDriver_Var2"

		hService, err := LoadDriverEx(
			driverPath,
			driverName,
			service.SERVICE_ALL_ACCESS,
			service.SERVICE_DEMAND_START,
			service.SERVICE_ERROR_IGNORE,
		)
		if err != nil {
			t.Skipf("Cannot load driver (need admin): %v", err)
			return
		}
		defer UnloadDriver(hService)

		t.Log("[*] LoadDriverEx works")
	})

	t.Run("OptionsLoad", func(t *testing.T) {
		driverPath := `C:\Windows\System32\drivers\null.sys`
		driverName := "NullDriver_Var3"

		options := DriverLoadOptions{
			DesiredAccess:    service.SERVICE_ALL_ACCESS,
			StartType:        service.SERVICE_AUTO_START,
			ErrorControl:     service.SERVICE_ERROR_NORMAL,
			StartImmediately: false,
			RecreateIfExists: true,
		}

		hService, err := LoadDriverWithOptions(driverPath, driverName, options)
		if err != nil {
			t.Skipf("Cannot load driver (need admin): %v", err)
			return
		}
		defer UnloadDriver(hService)

		// Manually start it
		if err := StartDriver(hService); err != nil {
			t.Logf("Failed to start: %v", err)
		} else {
			t.Log("[*] StartDriver works")
		}

		// Query status
		status, err := QueryDriverStatus(hService)
		if err == nil {
			t.Logf("[*] QueryDriverStatus works: state=%d", status.CurrentState)
		}

		// Stop it
		if err := StopDriver(hService); err != nil {
			t.Logf("Failed to stop: %v", err)
		} else {
			t.Log("[*] StopDriver works")
		}
	})

	t.Run("OpenExisting", func(t *testing.T) {
		hService, err := OpenExistingDriver("Beep", service.SERVICE_QUERY_STATUS)
		if err != nil {
			t.Skipf("Cannot open Beep driver: %v", err)
			return
		}
		defer service.CloseServiceHandle(hService)

		status, err := QueryDriverStatus(hService)
		if err != nil {
			t.Logf("Cannot query status: %v", err)
		} else {
			t.Logf("[*] OpenExistingDriver works: Beep driver state=%d", status.CurrentState)
		}
	})
}
