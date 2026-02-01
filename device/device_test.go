package device

import (
	"strings"
	"testing"
)

// TestCTL_CODE tests the IOCTL code construction macro
func TestCTL_CODE(t *testing.T) {
	tests := []struct {
		name       string
		deviceType uint32
		function   uint32
		method     uint32
		access     uint32
		expected   uint32
	}{
		{
			name:       "IOCTL_DISK_GET_DRIVE_GEOMETRY",
			deviceType: FILE_DEVICE_DISK,
			function:   0,
			method:     METHOD_BUFFERED,
			access:     FILE_ANY_ACCESS,
			expected:   IOCTL_DISK_GET_DRIVE_GEOMETRY,
		},
		{
			name:       "Custom IOCTL",
			deviceType: 0x8000,
			function:   0x800,
			method:     METHOD_BUFFERED,
			access:     FILE_ANY_ACCESS,
			expected:   (0x8000 << 16) | (0 << 14) | (0x800 << 2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CTL_CODE(tt.deviceType, tt.function, tt.method, tt.access)
			if result != tt.expected {
				t.Errorf("CTL_CODE() = 0x%08X, want 0x%08X", result, tt.expected)
			}
		})
	}
}

// TestOpenDeviceNul tests opening the NUL device (should always work)
func TestOpenDeviceNul(t *testing.T) {
	hDevice, err := CreateFile(
		"NUL",
		GENERIC_READ|GENERIC_WRITE,
		FILE_SHARE_READ|FILE_SHARE_WRITE,
		nil,
		OPEN_EXISTING,
		0,
		0,
	)

	if err != nil {
		t.Fatalf("Failed to open NUL device: %v", err)
	}

	if hDevice == 0 || uintptr(hDevice) == INVALID_HANDLE_VALUE {
		t.Fatal("Invalid handle returned for NUL device")
	}

	if !CloseHandle(hDevice) {
		t.Error("Failed to close NUL device handle")
	}

	t.Logf("Successfully opened and closed NUL device with handle: 0x%x", hDevice)
}

// TestOpenDeviceReadOnly tests the convenience function for read-only access
func TestOpenDeviceReadOnly(t *testing.T) {
	// Try to open the NUL device read-only
	hDevice, err := OpenDeviceReadOnly("NUL")
	if err != nil {
		t.Fatalf("OpenDeviceReadOnly() failed: %v", err)
	}

	if hDevice == 0 {
		t.Fatal("OpenDeviceReadOnly() returned null handle")
	}

	if !CloseHandle(hDevice) {
		t.Error("Failed to close device handle")
	}

	t.Logf("Successfully opened NUL device read-only with handle: 0x%x", hDevice)
}

// TestGetFileSize tests getting file size (using NUL device)
func TestGetFileSize(t *testing.T) {
	hDevice, err := OpenDeviceReadOnly("NUL")
	if err != nil {
		t.Fatalf("Failed to open NUL device: %v", err)
	}
	defer CloseHandle(hDevice)

	size, err := GetFileSize(hDevice)
	if err != nil {
		t.Logf("GetFileSize() returned error (expected for NUL): %v", err)
	} else {
		t.Logf("NUL device size: %d bytes", size)
	}
}

// TestReadWriteFile tests reading and writing to a device (using NUL)
func TestReadWriteFile(t *testing.T) {
	hDevice, err := OpenDeviceReadWrite("NUL")
	if err != nil {
		t.Fatalf("Failed to open NUL device: %v", err)
	}
	defer CloseHandle(hDevice)

	// Test write
	testData := []byte("Hello, winx device package!")
	var bytesWritten uint32

	success, err := WriteFile(hDevice, testData, uint32(len(testData)), &bytesWritten, nil)
	if !success {
		t.Logf("WriteFile() failed (might be expected for NUL): %v", err)
	} else {
		t.Logf("Successfully wrote %d bytes to NUL device", bytesWritten)
	}

	// Test read
	readBuffer := make([]byte, 100)
	var bytesRead uint32

	success, err = ReadFile(hDevice, readBuffer, 100, &bytesRead, nil)
	if !success {
		t.Logf("ReadFile() failed (expected for NUL): %v", err)
	} else {
		t.Logf("Read %d bytes from NUL device", bytesRead)
	}
}

// TestDeviceIoControlBytes tests the byte slice wrapper for DeviceIoControl
func TestDeviceIoControlBytes(t *testing.T) {
	// This is a basic test - actual IOCTL testing would require specific devices
	hDevice, err := OpenDeviceReadOnly("NUL")
	if err != nil {
		t.Fatalf("Failed to open NUL device: %v", err)
	}
	defer CloseHandle(hDevice)

	// Try a simple IOCTL (will likely fail on NUL, but tests the function)
	_, _, err = DeviceIoControlBytes(hDevice, 0x12345678, nil, 0)
	if err != nil {
		t.Logf("DeviceIoControlBytes() failed as expected: %v", err)
	}
}

// BenchmarkCreateFileClose benchmarks opening and closing a device
func BenchmarkCreateFileClose(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hDevice, err := CreateFile(
			"NUL",
			GENERIC_READ,
			FILE_SHARE_READ,
			nil,
			OPEN_EXISTING,
			0,
			0,
		)
		if err != nil {
			b.Fatalf("CreateFile failed: %v", err)
		}
		CloseHandle(hDevice)
	}
}

// BenchmarkCTL_CODE benchmarks the IOCTL code construction
func BenchmarkCTL_CODE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CTL_CODE(FILE_DEVICE_DISK, 0x100, METHOD_BUFFERED, FILE_ANY_ACCESS)
	}
}

// ============================================================================
// IOCTL Decoder Tests
// ============================================================================

// TestExtractDeviceType tests the device type extraction
func TestExtractDeviceType(t *testing.T) {
	tests := []struct {
		name       string
		ioctlCode  uint32
		expected   uint32
	}{
		{
			name:      "IOCTL_DISK_GET_DRIVE_GEOMETRY",
			ioctlCode: IOCTL_DISK_GET_DRIVE_GEOMETRY, // 0x70000
			expected:  FILE_DEVICE_DISK,               // 0x07
		},
		{
			name:      "IOCTL_STORAGE_GET_DEVICE_NUMBER",
			ioctlCode: IOCTL_STORAGE_GET_DEVICE_NUMBER,
			expected:  FILE_DEVICE_MASS_STORAGE, // 0x2D
		},
		{
			name:      "Custom device type",
			ioctlCode: 0x80000000,
			expected:  0x8000,
		},
		{
			name:      "Zero IOCTL",
			ioctlCode: 0x00000000,
			expected:  0x0000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractDeviceType(tt.ioctlCode)
			if result != tt.expected {
				t.Errorf("ExtractDeviceType(0x%08X) = 0x%04X, want 0x%04X",
					tt.ioctlCode, result, tt.expected)
			}
		})
	}
}

// TestExtractFunction tests the function code extraction
func TestExtractFunction(t *testing.T) {
	tests := []struct {
		name       string
		ioctlCode  uint32
		expected   uint32
	}{
		{
			name:      "IOCTL_DISK_GET_DRIVE_GEOMETRY (function 0)",
			ioctlCode: IOCTL_DISK_GET_DRIVE_GEOMETRY,
			expected:  0,
		},
		{
			name:      "IOCTL_DISK_GET_PARTITION_INFO (function 1)",
			ioctlCode: IOCTL_DISK_GET_PARTITION_INFO,
			expected:  1,
		},
		{
			name:      "IOCTL_DISK_GET_DRIVE_LAYOUT (function 3)",
			ioctlCode: IOCTL_DISK_GET_DRIVE_LAYOUT,
			expected:  3,
		},
		{
			name:      "IOCTL_STORAGE_QUERY_PROPERTY (function 0x500)",
			ioctlCode: IOCTL_STORAGE_QUERY_PROPERTY,
			expected:  0x500,
		},
		{
			name:      "Custom function code",
			ioctlCode: (FILE_DEVICE_DISK << 16) | (0xFFF << 2), // Max function code
			expected:  0xFFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractFunction(tt.ioctlCode)
			if result != tt.expected {
				t.Errorf("ExtractFunction(0x%08X) = 0x%03X, want 0x%03X",
					tt.ioctlCode, result, tt.expected)
			}
		})
	}
}

// TestExtractMethod tests the transfer method extraction
func TestExtractMethod(t *testing.T) {
	tests := []struct {
		name       string
		ioctlCode  uint32
		expected   uint32
	}{
		{
			name:      "METHOD_BUFFERED",
			ioctlCode: IOCTL_DISK_GET_DRIVE_GEOMETRY,
			expected:  METHOD_BUFFERED,
		},
		{
			name:      "METHOD_IN_DIRECT",
			ioctlCode: (FILE_DEVICE_DISK << 16) | METHOD_IN_DIRECT,
			expected:  METHOD_IN_DIRECT,
		},
		{
			name:      "METHOD_OUT_DIRECT",
			ioctlCode: (FILE_DEVICE_DISK << 16) | METHOD_OUT_DIRECT,
			expected:  METHOD_OUT_DIRECT,
		},
		{
			name:      "METHOD_NEITHER",
			ioctlCode: (FILE_DEVICE_DISK << 16) | METHOD_NEITHER,
			expected:  METHOD_NEITHER,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractMethod(tt.ioctlCode)
			if result != tt.expected {
				t.Errorf("ExtractMethod(0x%08X) = %d, want %d",
					tt.ioctlCode, result, tt.expected)
			}
		})
	}
}

// TestExtractAccess tests the access level extraction
func TestExtractAccess(t *testing.T) {
	tests := []struct {
		name       string
		ioctlCode  uint32
		expected   uint32
	}{
		{
			name:      "FILE_ANY_ACCESS",
			ioctlCode: IOCTL_DISK_GET_DRIVE_GEOMETRY,
			expected:  FILE_ANY_ACCESS,
		},
		{
			name:      "FILE_READ_ACCESS",
			ioctlCode: IOCTL_DISK_GET_PARTITION_INFO,
			expected:  FILE_READ_ACCESS,
		},
		{
			name:      "FILE_WRITE_ACCESS",
			ioctlCode: (FILE_DEVICE_DISK << 16) | (FILE_WRITE_ACCESS << 14),
			expected:  FILE_WRITE_ACCESS,
		},
		{
			name:      "FILE_READ_WRITE",
			ioctlCode: (FILE_DEVICE_DISK << 16) | (3 << 14), // Both read and write
			expected:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAccess(tt.ioctlCode)
			if result != tt.expected {
				t.Errorf("ExtractAccess(0x%08X) = %d, want %d",
					tt.ioctlCode, result, tt.expected)
			}
		})
	}
}

// TestGetDeviceTypeName tests the device type name lookup
func TestGetDeviceTypeName(t *testing.T) {
	tests := []struct {
		name       string
		deviceType uint32
		expected   string
	}{
		{
			name:       "FILE_DEVICE_DISK",
			deviceType: FILE_DEVICE_DISK,
			expected:   "DISK",
		},
		{
			name:       "FILE_DEVICE_KEYBOARD",
			deviceType: FILE_DEVICE_KEYBOARD,
			expected:   "KEYBOARD",
		},
		{
			name:       "FILE_DEVICE_MASS_STORAGE",
			deviceType: FILE_DEVICE_MASS_STORAGE,
			expected:   "MASS_STORAGE",
		},
		{
			name:       "Custom device type",
			deviceType: 0x8000,
			expected:   "CUSTOM",
		},
		{
			name:       "Unknown device type",
			deviceType: 0x003C, // Just beyond the known range (0x00-0x3B)
			expected:   "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDeviceTypeName(tt.deviceType)
			if result != tt.expected {
				t.Errorf("GetDeviceTypeName(0x%04X) = %q, want %q",
					tt.deviceType, result, tt.expected)
			}
		})
	}
}

// TestGetMethodName tests the method name lookup
func TestGetMethodName(t *testing.T) {
	tests := []struct {
		name     string
		method   uint32
		expected string
	}{
		{
			name:     "METHOD_BUFFERED",
			method:   METHOD_BUFFERED,
			expected: "BUFFERED",
		},
		{
			name:     "METHOD_IN_DIRECT",
			method:   METHOD_IN_DIRECT,
			expected: "IN_DIRECT",
		},
		{
			name:     "METHOD_OUT_DIRECT",
			method:   METHOD_OUT_DIRECT,
			expected: "OUT_DIRECT",
		},
		{
			name:     "METHOD_NEITHER",
			method:   METHOD_NEITHER,
			expected: "NEITHER",
		},
		{
			name:     "Invalid method",
			method:   99,
			expected: "INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMethodName(tt.method)
			if result != tt.expected {
				t.Errorf("GetMethodName(%d) = %q, want %q",
					tt.method, result, tt.expected)
			}
		})
	}
}

// TestGetAccessName tests the access name lookup
func TestGetAccessName(t *testing.T) {
	tests := []struct {
		name     string
		access   uint32
		expected string
	}{
		{
			name:     "FILE_ANY_ACCESS",
			access:   FILE_ANY_ACCESS,
			expected: "ANY",
		},
		{
			name:     "FILE_READ_ACCESS",
			access:   FILE_READ_ACCESS,
			expected: "READ",
		},
		{
			name:     "FILE_WRITE_ACCESS",
			access:   FILE_WRITE_ACCESS,
			expected: "WRITE",
		},
		{
			name:     "FILE_READ_WRITE",
			access:   3,
			expected: "READ_WRITE",
		},
		{
			name:     "Invalid access",
			access:   99,
			expected: "INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAccessName(tt.access)
			if result != tt.expected {
				t.Errorf("GetAccessName(%d) = %q, want %q",
					tt.access, result, tt.expected)
			}
		})
	}
}

// TestLookupKnownIOCTL tests the known IOCTL lookup
func TestLookupKnownIOCTL(t *testing.T) {
	tests := []struct {
		name       string
		ioctlCode  uint32
		expectName string
		expectFound bool
	}{
		{
			name:       "IOCTL_DISK_GET_DRIVE_GEOMETRY",
			ioctlCode:  IOCTL_DISK_GET_DRIVE_GEOMETRY,
			expectName: "IOCTL_DISK_GET_DRIVE_GEOMETRY",
			expectFound: true,
		},
		{
			name:       "IOCTL_STORAGE_QUERY_PROPERTY",
			ioctlCode:  IOCTL_STORAGE_QUERY_PROPERTY,
			expectName: "IOCTL_STORAGE_QUERY_PROPERTY",
			expectFound: true,
		},
		{
			name:       "Unknown IOCTL",
			ioctlCode:  0x12345678,
			expectName: "",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, found := LookupKnownIOCTL(tt.ioctlCode)
			if found != tt.expectFound {
				t.Errorf("LookupKnownIOCTL(0x%08X) found = %v, want %v",
					tt.ioctlCode, found, tt.expectFound)
			}
			if found && name != tt.expectName {
				t.Errorf("LookupKnownIOCTL(0x%08X) = %q, want %q",
					tt.ioctlCode, name, tt.expectName)
			}
		})
	}
}

// TestIOCTLRoundTrip tests encoding and decoding consistency
func TestIOCTLRoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		deviceType uint32
		function   uint32
		method     uint32
		access     uint32
	}{
		{
			name:       "Standard disk IOCTL",
			deviceType: FILE_DEVICE_DISK,
			function:   0x123,
			method:     METHOD_BUFFERED,
			access:     FILE_ANY_ACCESS,
		},
		{
			name:       "Network IOCTL with read access",
			deviceType: FILE_DEVICE_NETWORK,
			function:   0x456,
			method:     METHOD_IN_DIRECT,
			access:     FILE_READ_ACCESS,
		},
		{
			name:       "Custom device IOCTL",
			deviceType: 0x8000,
			function:   0xFFF,
			method:     METHOD_NEITHER,
			access:     3, // READ_WRITE
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			code := CTL_CODE(tt.deviceType, tt.function, tt.method, tt.access)

			// Decode
			deviceType := ExtractDeviceType(code)
			function := ExtractFunction(code)
			method := ExtractMethod(code)
			access := ExtractAccess(code)

			// Verify
			if deviceType != tt.deviceType {
				t.Errorf("DeviceType: got 0x%04X, want 0x%04X", deviceType, tt.deviceType)
			}
			if function != tt.function {
				t.Errorf("Function: got 0x%03X, want 0x%03X", function, tt.function)
			}
			if method != tt.method {
				t.Errorf("Method: got %d, want %d", method, tt.method)
			}
			if access != tt.access {
				t.Errorf("Access: got %d, want %d", access, tt.access)
			}
		})
	}
}

// BenchmarkExtractDeviceType benchmarks the device type extraction
func BenchmarkExtractDeviceType(b *testing.B) {
	code := uint32(IOCTL_DISK_GET_DRIVE_GEOMETRY)
	for i := 0; i < b.N; i++ {
		_ = ExtractDeviceType(code)
	}
}

// BenchmarkExtractFunction benchmarks the function extraction
func BenchmarkExtractFunction(b *testing.B) {
	code := uint32(IOCTL_DISK_GET_DRIVE_GEOMETRY)
	for i := 0; i < b.N; i++ {
		_ = ExtractFunction(code)
	}
}

// BenchmarkGetDeviceTypeName benchmarks the device type name lookup
func BenchmarkGetDeviceTypeName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetDeviceTypeName(FILE_DEVICE_DISK)
	}
}

// BenchmarkLookupKnownIOCTL benchmarks the known IOCTL lookup
func BenchmarkLookupKnownIOCTL(b *testing.B) {
	code := uint32(IOCTL_DISK_GET_DRIVE_GEOMETRY)
	for i := 0; i < b.N; i++ {
		_, _ = LookupKnownIOCTL(code)
	}
}

// TestDecodeIOCTL tests the complete IOCTL decoder
func TestDecodeIOCTL(t *testing.T) {
	tests := []struct {
		name               string
		ioctlCode          uint32
		expectedDeviceType uint32
		expectedDeviceName string
		expectedFunction   uint32
		expectedMethod     uint32
		expectedMethodName string
		expectedAccess     uint32
		expectedAccessName string
		expectedKnownName  string
	}{
		{
			name:               "IOCTL_DISK_GET_DRIVE_GEOMETRY",
			ioctlCode:          IOCTL_DISK_GET_DRIVE_GEOMETRY,
			expectedDeviceType: FILE_DEVICE_DISK,
			expectedDeviceName: "DISK",
			expectedFunction:   0,
			expectedMethod:     METHOD_BUFFERED,
			expectedMethodName: "BUFFERED",
			expectedAccess:     FILE_ANY_ACCESS,
			expectedAccessName: "ANY",
			expectedKnownName:  "IOCTL_DISK_GET_DRIVE_GEOMETRY",
		},
		{
			name:               "IOCTL_DISK_GET_PARTITION_INFO",
			ioctlCode:          IOCTL_DISK_GET_PARTITION_INFO,
			expectedDeviceType: FILE_DEVICE_DISK,
			expectedDeviceName: "DISK",
			expectedFunction:   1,
			expectedMethod:     METHOD_BUFFERED,
			expectedMethodName: "BUFFERED",
			expectedAccess:     FILE_READ_ACCESS,
			expectedAccessName: "READ",
			expectedKnownName:  "IOCTL_DISK_GET_PARTITION_INFO",
		},
		{
			name:               "IOCTL_STORAGE_QUERY_PROPERTY",
			ioctlCode:          IOCTL_STORAGE_QUERY_PROPERTY,
			expectedDeviceType: FILE_DEVICE_MASS_STORAGE,
			expectedDeviceName: "MASS_STORAGE",
			expectedFunction:   0x500,
			expectedMethod:     METHOD_BUFFERED,
			expectedMethodName: "BUFFERED",
			expectedAccess:     FILE_ANY_ACCESS,
			expectedAccessName: "ANY",
			expectedKnownName:  "IOCTL_STORAGE_QUERY_PROPERTY",
		},
		{
			name:               "Custom IOCTL - Network device",
			ioctlCode:          CTL_CODE(FILE_DEVICE_NETWORK, 0x123, METHOD_IN_DIRECT, FILE_WRITE_ACCESS),
			expectedDeviceType: FILE_DEVICE_NETWORK,
			expectedDeviceName: "NETWORK",
			expectedFunction:   0x123,
			expectedMethod:     METHOD_IN_DIRECT,
			expectedMethodName: "IN_DIRECT",
			expectedAccess:     FILE_WRITE_ACCESS,
			expectedAccessName: "WRITE",
			expectedKnownName:  "",
		},
		{
			name:               "Custom IOCTL - Custom device type",
			ioctlCode:          CTL_CODE(0x8000, 0xABC, METHOD_NEITHER, FILE_READ_ACCESS),
			expectedDeviceType: 0x8000,
			expectedDeviceName: "CUSTOM",
			expectedFunction:   0xABC,
			expectedMethod:     METHOD_NEITHER,
			expectedMethodName: "NEITHER",
			expectedAccess:     FILE_READ_ACCESS,
			expectedAccessName: "READ",
			expectedKnownName:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			components := DecodeIOCTL(tt.ioctlCode)

			// Verify all components
			if components.IOCTLCode != tt.ioctlCode {
				t.Errorf("IOCTLCode: got 0x%08X, want 0x%08X", components.IOCTLCode, tt.ioctlCode)
			}
			if components.DeviceType != tt.expectedDeviceType {
				t.Errorf("DeviceType: got 0x%04X, want 0x%04X", components.DeviceType, tt.expectedDeviceType)
			}
			if components.DeviceTypeName != tt.expectedDeviceName {
				t.Errorf("DeviceTypeName: got %q, want %q", components.DeviceTypeName, tt.expectedDeviceName)
			}
			if components.Function != tt.expectedFunction {
				t.Errorf("Function: got 0x%03X, want 0x%03X", components.Function, tt.expectedFunction)
			}
			if components.Method != tt.expectedMethod {
				t.Errorf("Method: got %d, want %d", components.Method, tt.expectedMethod)
			}
			if components.MethodName != tt.expectedMethodName {
				t.Errorf("MethodName: got %q, want %q", components.MethodName, tt.expectedMethodName)
			}
			if components.Access != tt.expectedAccess {
				t.Errorf("Access: got %d, want %d", components.Access, tt.expectedAccess)
			}
			if components.AccessName != tt.expectedAccessName {
				t.Errorf("AccessName: got %q, want %q", components.AccessName, tt.expectedAccessName)
			}
			if components.KnownName != tt.expectedKnownName {
				t.Errorf("KnownName: got %q, want %q", components.KnownName, tt.expectedKnownName)
			}
		})
	}
}

// BenchmarkDecodeIOCTL benchmarks the complete IOCTL decoder
func BenchmarkDecodeIOCTL(b *testing.B) {
	code := uint32(IOCTL_DISK_GET_DRIVE_GEOMETRY)
	for i := 0; i < b.N; i++ {
		_ = DecodeIOCTL(code)
	}
}

// TestFormatIOCTL tests the compact IOCTL formatting
func TestFormatIOCTL(t *testing.T) {
	tests := []struct {
		name     string
		ioctl    uint32
		expected string
	}{
		{
			name:     "Known IOCTL with name",
			ioctl:    IOCTL_DISK_GET_DRIVE_GEOMETRY,
			expected: "IOCTL_DISK_GET_DRIVE_GEOMETRY: DISK(function:0x0, method:BUFFERED, access:ANY)",
		},
		{
			name:     "Known IOCTL - partition info",
			ioctl:    IOCTL_DISK_GET_PARTITION_INFO,
			expected: "IOCTL_DISK_GET_PARTITION_INFO: DISK(function:0x1, method:BUFFERED, access:READ)",
		},
		{
			name:     "Unknown IOCTL",
			ioctl:    CTL_CODE(FILE_DEVICE_NETWORK, 0x123, METHOD_IN_DIRECT, FILE_WRITE_ACCESS),
			expected: "NETWORK(function:0x123, method:IN_DIRECT, access:WRITE)",
		},
		{
			name:     "Custom device type",
			ioctl:    CTL_CODE(0x8000, 0xABC, METHOD_NEITHER, FILE_READ_ACCESS),
			expected: "CUSTOM(function:0xABC, method:NEITHER, access:READ)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatIOCTL(tt.ioctl)
			if result != tt.expected {
				t.Errorf("FormatIOCTL() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestFormatIOCTLVerbose tests the verbose IOCTL formatting
func TestFormatIOCTLVerbose(t *testing.T) {
	// Test with a known IOCTL
	result := FormatIOCTLVerbose(IOCTL_DISK_GET_DRIVE_GEOMETRY)

	// Check for key components in the output
	expectedStrings := []string{
		"IOCTL Code: 0x00070000",
		"Known Name: IOCTL_DISK_GET_DRIVE_GEOMETRY",
		"Device Type: DISK (0x0007)",
		"Function: 0x000 (0)",
		"Method: BUFFERED (0)",
		"Access: ANY (0)",
		"Bit Layout:",
		"[31:16] Device Type:",
		"[15:14] Access:",
		"[13:2]  Function:",
		"[1:0]   Method:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("FormatIOCTLVerbose() missing expected string %q\nGot:\n%s", expected, result)
		}
	}
}

// TestFormatIOCTLVerboseCustom tests verbose formatting with custom IOCTL
func TestFormatIOCTLVerboseCustom(t *testing.T) {
	// Test with a custom IOCTL
	customIOCTL := CTL_CODE(FILE_DEVICE_NETWORK, 0x123, METHOD_IN_DIRECT, FILE_WRITE_ACCESS)
	result := FormatIOCTLVerbose(customIOCTL)

	// Should not have Known Name section for custom IOCTLs
	if strings.Contains(result, "Known Name:") {
		t.Errorf("FormatIOCTLVerbose() should not have 'Known Name' for custom IOCTL\nGot:\n%s", result)
	}

	// Check for key components
	expectedStrings := []string{
		"IOCTL Code: 0x0012848D",
		"Device Type: NETWORK",
		"Function: 0x123",
		"Method: IN_DIRECT",
		"Access: WRITE",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("FormatIOCTLVerbose() missing expected string %q\nGot:\n%s", expected, result)
		}
	}
}

// TestFormatIOCTLHex tests the hexadecimal IOCTL formatting
func TestFormatIOCTLHex(t *testing.T) {
	tests := []struct {
		name     string
		ioctl    uint32
		expected string
	}{
		{
			name:     "IOCTL_DISK_GET_DRIVE_GEOMETRY",
			ioctl:    IOCTL_DISK_GET_DRIVE_GEOMETRY,
			expected: "0x00070000 [Device:0x0007 Access:0x0 Func:0x000 Method:0x0]",
		},
		{
			name:     "IOCTL_DISK_GET_PARTITION_INFO",
			ioctl:    IOCTL_DISK_GET_PARTITION_INFO,
			expected: "0x00074004 [Device:0x0007 Access:0x1 Func:0x001 Method:0x0]",
		},
		{
			name:     "Custom Network IOCTL",
			ioctl:    CTL_CODE(FILE_DEVICE_NETWORK, 0x123, METHOD_IN_DIRECT, FILE_WRITE_ACCESS),
			expected: "0x0012848D [Device:0x0012 Access:0x2 Func:0x123 Method:0x1]",
		},
		{
			name:     "Custom device with all fields",
			ioctl:    CTL_CODE(0x8000, 0xFFF, METHOD_NEITHER, 3), // READ_WRITE access
			expected: "0x8000FFFF [Device:0x8000 Access:0x3 Func:0xFFF Method:0x3]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatIOCTLHex(tt.ioctl)
			if result != tt.expected {
				t.Errorf("FormatIOCTLHex() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// BenchmarkFormatIOCTL benchmarks the compact formatting
func BenchmarkFormatIOCTL(b *testing.B) {
	code := uint32(IOCTL_DISK_GET_DRIVE_GEOMETRY)
	for i := 0; i < b.N; i++ {
		_ = FormatIOCTL(code)
	}
}

// BenchmarkFormatIOCTLVerbose benchmarks the verbose formatting
func BenchmarkFormatIOCTLVerbose(b *testing.B) {
	code := uint32(IOCTL_DISK_GET_DRIVE_GEOMETRY)
	for i := 0; i < b.N; i++ {
		_ = FormatIOCTLVerbose(code)
	}
}

// BenchmarkFormatIOCTLHex benchmarks the hex formatting
func BenchmarkFormatIOCTLHex(b *testing.B) {
	code := uint32(IOCTL_DISK_GET_DRIVE_GEOMETRY)
	for i := 0; i < b.N; i++ {
		_ = FormatIOCTLHex(code)
	}
}
