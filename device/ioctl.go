package device

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/ArkaprabhaChakraborty/winx/handle"
)

var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	procCreateFileW       = kernel32.NewProc("CreateFileW")
	procCloseHandle       = kernel32.NewProc("CloseHandle")
	procDeviceIoControl   = kernel32.NewProc("DeviceIoControl")
	procGetLastError      = kernel32.NewProc("GetLastError")
	procReadFile          = kernel32.NewProc("ReadFile")
	procWriteFile         = kernel32.NewProc("WriteFile")
	procGetFileSizeEx     = kernel32.NewProc("GetFileSizeEx")
	procQueryDosDeviceW   = kernel32.NewProc("QueryDosDeviceW")
)

// CreateFile opens or creates a file or I/O device.
// This is the primary function for accessing devices.
//
// Parameters:
//   - fileName: The name of the file or device to be created or opened
//   - desiredAccess: The requested access to the file or device
//   - shareMode: The requested sharing mode of the file or device
//   - securityAttributes: A pointer to a SECURITY_ATTRIBUTES structure (can be nil)
//   - creationDisposition: An action to take on a file or device that exists or does not exist
//   - flagsAndAttributes: The file or device attributes and flags
//   - templateFile: A handle to a template file (usually 0)
//
// Returns:
//   - A handle to the file or device if successful, INVALID_HANDLE_VALUE otherwise
func CreateFile(
	fileName string,
	desiredAccess uint32,
	shareMode uint32,
	securityAttributes *SECURITY_ATTRIBUTES,
	creationDisposition uint32,
	flagsAndAttributes uint32,
	templateFile handle.HANDLE,
) (handle.HANDLE, error) {

	fileNamePtr, err := syscall.UTF16PtrFromString(fileName)
	if err != nil {
		return handle.HANDLE(INVALID_HANDLE_VALUE), err
	}

	var secAttrPtr uintptr
	if securityAttributes != nil {
		secAttrPtr = uintptr(unsafe.Pointer(securityAttributes))
	}

	ret, _, err := syscall.SyscallN(
		procCreateFileW.Addr(),
		uintptr(unsafe.Pointer(fileNamePtr)),
		uintptr(desiredAccess),
		uintptr(shareMode),
		secAttrPtr,
		uintptr(creationDisposition),
		uintptr(flagsAndAttributes),
		uintptr(templateFile),
	)

	if ret == INVALID_HANDLE_VALUE || ret == 0 {
		if err != nil {
			return handle.HANDLE(INVALID_HANDLE_VALUE), err
		}
		return handle.HANDLE(INVALID_HANDLE_VALUE), syscall.GetLastError()
	}

	return handle.HANDLE(ret), nil
}

// CloseHandle closes an open object handle.
//
// Parameters:
//   - hObject: A valid handle to an open object
//
// Returns:
//   - true if successful, false otherwise
func CloseHandle(hObject handle.HANDLE) bool {
	ret, _, _ := syscall.SyscallN(
		procCloseHandle.Addr(),
		uintptr(hObject),
	)
	return ret != 0
}

// DeviceIoControl sends a control code directly to a specified device driver,
// causing the corresponding device to perform the corresponding operation.
//
// Parameters:
//   - hDevice: A handle to the device on which the operation is to be performed
//   - ioControlCode: The control code for the operation
//   - inBuffer: A pointer to the input buffer (can be nil)
//   - inBufferSize: The size of the input buffer in bytes
//   - outBuffer: A pointer to the output buffer (can be nil)
//   - outBufferSize: The size of the output buffer in bytes
//   - bytesReturned: A pointer to a variable that receives the size of the data stored in the output buffer
//   - overlapped: A pointer to an OVERLAPPED structure (can be nil)
//
// Returns:
//   - true if successful, false otherwise
func DeviceIoControl(
	hDevice handle.HANDLE,
	ioControlCode uint32,
	inBuffer unsafe.Pointer,
	inBufferSize uint32,
	outBuffer unsafe.Pointer,
	outBufferSize uint32,
	bytesReturned *uint32,
	overlapped *OVERLAPPED,
) (bool, error) {

	var overlappedPtr uintptr
	if overlapped != nil {
		overlappedPtr = uintptr(unsafe.Pointer(overlapped))
	}

	ret, _, _ := syscall.SyscallN(
		procDeviceIoControl.Addr(),
		uintptr(hDevice),
		uintptr(ioControlCode),
		uintptr(inBuffer),
		uintptr(inBufferSize),
		uintptr(outBuffer),
		uintptr(outBufferSize),
		uintptr(unsafe.Pointer(bytesReturned)),
		overlappedPtr,
	)

	if ret == 0 {
		return false, syscall.GetLastError()
	}

	return true, nil
}

// DeviceIoControlBytes is a convenience wrapper around DeviceIoControl that works with byte slices.
//
// Parameters:
//   - hDevice: A handle to the device
//   - ioControlCode: The control code for the operation
//   - inBuffer: Input data as a byte slice (can be nil)
//   - outBufferSize: The size of the output buffer to allocate
//
// Returns:
//   - The output data as a byte slice, the number of bytes returned, and any error
func DeviceIoControlBytes(
	hDevice handle.HANDLE,
	ioControlCode uint32,
	inBuffer []byte,
	outBufferSize uint32,
) ([]byte, uint32, error) {

	var inPtr unsafe.Pointer
	var inSize uint32

	if len(inBuffer) > 0 {
		inPtr = unsafe.Pointer(&inBuffer[0])
		inSize = uint32(len(inBuffer))
	}

	outBuffer := make([]byte, outBufferSize)
	var outPtr unsafe.Pointer
	if outBufferSize > 0 {
		outPtr = unsafe.Pointer(&outBuffer[0])
	}

	var bytesReturned uint32
	success, err := DeviceIoControl(
		hDevice,
		ioControlCode,
		inPtr,
		inSize,
		outPtr,
		outBufferSize,
		&bytesReturned,
		nil,
	)

	if !success {
		return nil, 0, err
	}

	return outBuffer[:bytesReturned], bytesReturned, nil
}

// ReadFile reads data from a file or device.
//
// Parameters:
//   - hFile: A handle to the file or device
//   - buffer: A buffer to receive the data
//   - numberOfBytesToRead: The maximum number of bytes to read
//   - numberOfBytesRead: A pointer to a variable that receives the number of bytes read
//   - overlapped: A pointer to an OVERLAPPED structure (can be nil)
//
// Returns:
//   - true if successful, false otherwise
func ReadFile(
	hFile handle.HANDLE,
	buffer []byte,
	numberOfBytesToRead uint32,
	numberOfBytesRead *uint32,
	overlapped *OVERLAPPED,
) (bool, error) {

	var bufferPtr unsafe.Pointer
	if len(buffer) > 0 {
		bufferPtr = unsafe.Pointer(&buffer[0])
	}

	var overlappedPtr uintptr
	if overlapped != nil {
		overlappedPtr = uintptr(unsafe.Pointer(overlapped))
	}

	ret, _, _ := syscall.SyscallN(
		procReadFile.Addr(),
		uintptr(hFile),
		uintptr(bufferPtr),
		uintptr(numberOfBytesToRead),
		uintptr(unsafe.Pointer(numberOfBytesRead)),
		overlappedPtr,
	)

	if ret == 0 {
		return false, syscall.GetLastError()
	}

	return true, nil
}

// WriteFile writes data to a file or device.
//
// Parameters:
//   - hFile: A handle to the file or device
//   - buffer: The data to be written
//   - numberOfBytesToWrite: The number of bytes to write
//   - numberOfBytesWritten: A pointer to a variable that receives the number of bytes written
//   - overlapped: A pointer to an OVERLAPPED structure (can be nil)
//
// Returns:
//   - true if successful, false otherwise
func WriteFile(
	hFile handle.HANDLE,
	buffer []byte,
	numberOfBytesToWrite uint32,
	numberOfBytesWritten *uint32,
	overlapped *OVERLAPPED,
) (bool, error) {

	var bufferPtr unsafe.Pointer
	if len(buffer) > 0 {
		bufferPtr = unsafe.Pointer(&buffer[0])
	}

	var overlappedPtr uintptr
	if overlapped != nil {
		overlappedPtr = uintptr(unsafe.Pointer(overlapped))
	}

	ret, _, _ := syscall.SyscallN(
		procWriteFile.Addr(),
		uintptr(hFile),
		uintptr(bufferPtr),
		uintptr(numberOfBytesToWrite),
		uintptr(unsafe.Pointer(numberOfBytesWritten)),
		overlappedPtr,
	)

	if ret == 0 {
		return false, syscall.GetLastError()
	}

	return true, nil
}

// GetFileSize gets the size of a file.
//
// Parameters:
//   - hFile: A handle to the file
//
// Returns:
//   - The file size in bytes, or 0 and an error if the operation fails
func GetFileSize(hFile handle.HANDLE) (int64, error) {
	var fileSize int64

	ret, _, _ := syscall.SyscallN(
		procGetFileSizeEx.Addr(),
		uintptr(hFile),
		uintptr(unsafe.Pointer(&fileSize)),
	)

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return fileSize, nil
}

// OpenDevice is a convenience function to open a device for IOCTL operations.
//
// Parameters:
//   - devicePath: The device path (e.g., "\\\\.\\PhysicalDrive0")
//   - desiredAccess: The requested access (typically GENERIC_READ | GENERIC_WRITE)
//
// Returns:
//   - A handle to the device and any error
func OpenDevice(devicePath string, desiredAccess uint32) (handle.HANDLE, error) {
	return CreateFile(
		devicePath,
		desiredAccess,
		FILE_SHARE_READ|FILE_SHARE_WRITE,
		nil,
		OPEN_EXISTING,
		0,
		0,
	)
}

// OpenDeviceReadOnly is a convenience function to open a device for read-only access.
//
// Parameters:
//   - devicePath: The device path
//
// Returns:
//   - A handle to the device and any error
func OpenDeviceReadOnly(devicePath string) (handle.HANDLE, error) {
	return OpenDevice(devicePath, GENERIC_READ)
}

// OpenDeviceReadWrite is a convenience function to open a device for read-write access.
//
// Parameters:
//   - devicePath: The device path
//
// Returns:
//   - A handle to the device and any error
func OpenDeviceReadWrite(devicePath string) (handle.HANDLE, error) {
	return OpenDevice(devicePath, GENERIC_READ|GENERIC_WRITE)
}

// QueryDosDevice retrieves information about MS-DOS device names.
// This can be used to query symbolic links in the \DosDevices namespace.
//
// Parameters:
//   - deviceName: The MS-DOS device name (e.g., "C:", "COM1", or empty for all devices)
//
// Returns:
//   - A slice of target paths for the device, and any error
func QueryDosDevice(deviceName string) ([]string, error) {
	var deviceNamePtr *uint16
	var err error

	if deviceName != "" {
		deviceNamePtr, err = syscall.UTF16PtrFromString(deviceName)
		if err != nil {
			return nil, err
		}
	}

	// Initial buffer size
	bufferSize := uint32(65536) // 64KB should be enough for most cases
	buffer := make([]uint16, bufferSize)

	ret, _, err := syscall.SyscallN(
		procQueryDosDeviceW.Addr(),
		uintptr(unsafe.Pointer(deviceNamePtr)),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(bufferSize),
	)

	if ret == 0 {
		if err != nil {
			return nil, err
		}
		return nil, syscall.GetLastError()
	}

	// Parse the multi-string result
	// QueryDosDevice returns a list of null-terminated strings, terminated by an additional null
	var results []string
	start := 0
	for i := 0; i < len(buffer); i++ {
		if buffer[i] == 0 {
			if i == start {
				// Double null means end of list
				break
			}
			results = append(results, syscall.UTF16ToString(buffer[start:i]))
			start = i + 1
		}
	}

	return results, nil
}

// FindSymbolicLinksByPattern searches for DOS device symbolic links matching a pattern
//
// Parameters:
//   - pattern: A substring to search for in device names (case-insensitive, e.g., "clfs")
//
// Returns:
//   - A map of device names to their target paths
func FindSymbolicLinksByPattern(pattern string) (map[string][]string, error) {
	// Get all DOS devices
	allDevices, err := QueryDosDevice("")
	if err != nil {
		return nil, err
	}

	results := make(map[string][]string)

	// Search for devices matching the pattern (case-insensitive)
	patternLower := strings.ToLower(pattern)
	for _, device := range allDevices {
		// Case-insensitive search using strings.Contains
		if len(pattern) == 0 || strings.Contains(strings.ToLower(device), patternLower) {
			targets, err := QueryDosDevice(device)
			if err == nil && len(targets) > 0 {
				results[device] = targets
			}
		}
	}

	return results, nil
}


// ============================================================================
// IOCTL Code Decoding Functions
// ============================================================================

// ExtractDeviceType extracts the device type from an IOCTL code.
// The device type occupies bits 31-16 (upper 16 bits).
//
// Parameters:
//   - ioctlCode: The IOCTL code to decode
//
// Returns:
//   - The device type as a 16-bit value (e.g., 0x07 for FILE_DEVICE_DISK)
func ExtractDeviceType(ioctlCode uint32) uint32 {
	return (ioctlCode >> 16) & 0xFFFF
}

// ExtractFunction extracts the function code from an IOCTL code.
// The function code occupies bits 13-2 (12 bits).
//
// Parameters:
//   - ioctlCode: The IOCTL code to decode
//
// Returns:
//   - The function code as a value from 0 to 4095
func ExtractFunction(ioctlCode uint32) uint32 {
	return (ioctlCode >> 2) & 0xFFF
}

// ExtractMethod extracts the transfer method from an IOCTL code.
// The method occupies bits 1-0 (lower 2 bits).
//
// Parameters:
//   - ioctlCode: The IOCTL code to decode
//
// Returns:
//   - The transfer method (0=BUFFERED, 1=IN_DIRECT, 2=OUT_DIRECT, 3=NEITHER)
func ExtractMethod(ioctlCode uint32) uint32 {
	return ioctlCode & 0x3
}

// ExtractAccess extracts the access required from an IOCTL code.
// The access level occupies bits 15-14.
//
// Parameters:
//   - ioctlCode: The IOCTL code to decode
//
// Returns:
//   - The access level (0=ANY, 1=READ, 2=WRITE, 3=READ_WRITE)
func ExtractAccess(ioctlCode uint32) uint32 {
	return (ioctlCode >> 14) & 0x3
}

// ============================================================================
// Lookup Maps for IOCTL Components
// ============================================================================

// deviceTypeNames maps device type codes to human-readable names
// Uses constants already defined in constants.go
var deviceTypeNames = map[uint32]string{
	FILE_DEVICE_BEEP:               "BEEP",
	FILE_DEVICE_CD_ROM:             "CD_ROM",
	FILE_DEVICE_CD_ROM_FILE_SYSTEM: "CD_ROM_FILE_SYSTEM",
	FILE_DEVICE_CONTROLLER:         "CONTROLLER",
	FILE_DEVICE_DATALINK:           "DATALINK",
	FILE_DEVICE_DFS:                "DFS",
	FILE_DEVICE_DISK:               "DISK",
	FILE_DEVICE_DISK_FILE_SYSTEM:   "DISK_FILE_SYSTEM",
	FILE_DEVICE_FILE_SYSTEM:        "FILE_SYSTEM",
	FILE_DEVICE_INPORT_PORT:        "INPORT_PORT",
	FILE_DEVICE_KEYBOARD:           "KEYBOARD",
	FILE_DEVICE_MAILSLOT:           "MAILSLOT",
	FILE_DEVICE_MIDI_IN:            "MIDI_IN",
	FILE_DEVICE_MIDI_OUT:           "MIDI_OUT",
	FILE_DEVICE_MOUSE:              "MOUSE",
	FILE_DEVICE_MULTI_UNC_PROVIDER: "MULTI_UNC_PROVIDER",
	FILE_DEVICE_NAMED_PIPE:         "NAMED_PIPE",
	FILE_DEVICE_NETWORK:            "NETWORK",
	FILE_DEVICE_NETWORK_BROWSER:    "NETWORK_BROWSER",
	FILE_DEVICE_NETWORK_FILE_SYSTEM: "NETWORK_FILE_SYSTEM",
	FILE_DEVICE_NULL:               "NULL",
	FILE_DEVICE_PARALLEL_PORT:      "PARALLEL_PORT",
	FILE_DEVICE_PHYSICAL_NETCARD:   "PHYSICAL_NETCARD",
	FILE_DEVICE_PRINTER:            "PRINTER",
	FILE_DEVICE_SCANNER:            "SCANNER",
	FILE_DEVICE_SERIAL_MOUSE_PORT:  "SERIAL_MOUSE_PORT",
	FILE_DEVICE_SERIAL_PORT:        "SERIAL_PORT",
	FILE_DEVICE_SCREEN:             "SCREEN",
	FILE_DEVICE_SOUND:              "SOUND",
	FILE_DEVICE_STREAMS:            "STREAMS",
	FILE_DEVICE_TAPE:               "TAPE",
	FILE_DEVICE_TAPE_FILE_SYSTEM:   "TAPE_FILE_SYSTEM",
	FILE_DEVICE_TRANSPORT:          "TRANSPORT",
	FILE_DEVICE_UNKNOWN:            "UNKNOWN",
	FILE_DEVICE_VIDEO:              "VIDEO",
	FILE_DEVICE_VIRTUAL_DISK:       "VIRTUAL_DISK",
	FILE_DEVICE_WAVE_IN:            "WAVE_IN",
	FILE_DEVICE_WAVE_OUT:           "WAVE_OUT",
	FILE_DEVICE_8042_PORT:          "8042_PORT",
	FILE_DEVICE_NETWORK_REDIRECTOR: "NETWORK_REDIRECTOR",
	FILE_DEVICE_BATTERY:            "BATTERY",
	FILE_DEVICE_BUS_EXTENDER:       "BUS_EXTENDER",
	FILE_DEVICE_MODEM:              "MODEM",
	FILE_DEVICE_VDM:                "VDM",
	FILE_DEVICE_MASS_STORAGE:       "MASS_STORAGE",
	FILE_DEVICE_SMB:                "SMB",
	FILE_DEVICE_KS:                 "KS",
	FILE_DEVICE_CHANGER:            "CHANGER",
	FILE_DEVICE_SMARTCARD:          "SMARTCARD",
	FILE_DEVICE_ACPI:               "ACPI",
	FILE_DEVICE_DVD:                "DVD",
	FILE_DEVICE_FULLSCREEN_VIDEO:   "FULLSCREEN_VIDEO",
	FILE_DEVICE_DFS_FILE_SYSTEM:    "DFS_FILE_SYSTEM",
	FILE_DEVICE_DFS_VOLUME:         "DFS_VOLUME",
	FILE_DEVICE_SERENUM:            "SERENUM",
	FILE_DEVICE_TERMSRV:            "TERMSRV",
	FILE_DEVICE_KSEC:               "KSEC",
	FILE_DEVICE_FIPS:               "FIPS",
	FILE_DEVICE_INFINIBAND:         "INFINIBAND",
}

// methodNames maps transfer method codes to human-readable names
var methodNames = map[uint32]string{
	METHOD_BUFFERED:   "BUFFERED",
	METHOD_IN_DIRECT:  "IN_DIRECT",
	METHOD_OUT_DIRECT: "OUT_DIRECT",
	METHOD_NEITHER:    "NEITHER",
}

// accessNames maps access level codes to human-readable names
var accessNames = map[uint32]string{
	FILE_ANY_ACCESS:     "ANY",
	FILE_READ_ACCESS:    "READ",
	FILE_WRITE_ACCESS:   "WRITE",
	3:                   "READ_WRITE", // FILE_READ_ACCESS | FILE_WRITE_ACCESS
}

// knownIOCTLs maps known IOCTL codes to their symbolic names
// Uses constants already defined in constants.go
var knownIOCTLs = map[uint32]string{
	IOCTL_DISK_GET_DRIVE_GEOMETRY:       "IOCTL_DISK_GET_DRIVE_GEOMETRY",
	IOCTL_DISK_GET_PARTITION_INFO:       "IOCTL_DISK_GET_PARTITION_INFO",
	IOCTL_DISK_GET_DRIVE_LAYOUT:         "IOCTL_DISK_GET_DRIVE_LAYOUT",
	IOCTL_STORAGE_GET_DEVICE_NUMBER:     "IOCTL_STORAGE_GET_DEVICE_NUMBER",
	IOCTL_STORAGE_QUERY_PROPERTY:        "IOCTL_STORAGE_QUERY_PROPERTY",
	IOCTL_VOLUME_GET_VOLUME_DISK_EXTENTS: "IOCTL_VOLUME_GET_VOLUME_DISK_EXTENTS",
}

// ============================================================================
// Name Lookup Functions
// ============================================================================

// GetDeviceTypeName returns the human-readable name for a device type code.
// If the device type is unknown, it returns "UNKNOWN" with the hex value.
//
// Parameters:
//   - deviceType: The device type code (0x00-0x3B or custom 0x8000+)
//
// Returns:
//   - The device type name as a string
func GetDeviceTypeName(deviceType uint32) string {
	if name, ok := deviceTypeNames[deviceType]; ok {
		return name
	}
	// Custom device types start at 0x8000
	if deviceType >= 0x8000 {
		return "CUSTOM"
	}
	return "UNKNOWN"
}

// GetMethodName returns the human-readable name for a transfer method.
//
// Parameters:
//   - method: The method code (0-3)
//
// Returns:
//   - The method name as a string
func GetMethodName(method uint32) string {
	if name, ok := methodNames[method]; ok {
		return name
	}
	return "INVALID"
}

// GetAccessName returns the human-readable name for an access level.
//
// Parameters:
//   - access: The access code (0-3)
//
// Returns:
//   - The access name as a string
func GetAccessName(access uint32) string {
	if name, ok := accessNames[access]; ok {
		return name
	}
	return "INVALID"
}

// LookupKnownIOCTL searches for a known IOCTL code and returns its name.
//
// Parameters:
//   - ioctlCode: The IOCTL code to look up
//
// Returns:
//   - The IOCTL name if found, and a boolean indicating success
func LookupKnownIOCTL(ioctlCode uint32) (string, bool) {
	name, found := knownIOCTLs[ioctlCode]
	return name, found
}

// ============================================================================
// IOCTL Decoder
// ============================================================================

// DecodeIOCTL decodes an IOCTL code into its component parts with human-readable names.
// This is the main entry point for IOCTL code analysis.
//
// The function extracts and interprets all components of the 32-bit IOCTL code:
//   - Device Type: The type of device (bits 31-16)
//   - Function: The specific operation code (bits 13-2)
//   - Method: The data transfer method (bits 1-0)
//   - Access: Required access rights (bits 15-14)
//
// It also provides human-readable names for all components and checks if the
// IOCTL is a known standard Windows IOCTL.
//
// Parameters:
//   - ioctlCode: The 32-bit IOCTL code to decode
//
// Returns:
//   - A pointer to IOCTLComponents containing all decoded information
//
// Example:
//
//	components := DecodeIOCTL(IOCTL_DISK_GET_DRIVE_GEOMETRY)
//	fmt.Printf("Device: %s\n", components.DeviceTypeName)
//	fmt.Printf("Function: 0x%X\n", components.Function)
//	fmt.Printf("Method: %s\n", components.MethodName)
//	if components.KnownName != "" {
//	    fmt.Printf("Known IOCTL: %s\n", components.KnownName)
//	}
func DecodeIOCTL(ioctlCode uint32) *IOCTLComponents {
	// Extract all bit fields
	deviceType := ExtractDeviceType(ioctlCode)
	function := ExtractFunction(ioctlCode)
	method := ExtractMethod(ioctlCode)
	access := ExtractAccess(ioctlCode)

	// Create the components structure
	components := &IOCTLComponents{
		IOCTLCode:      ioctlCode,
		DeviceType:     deviceType,
		DeviceTypeName: GetDeviceTypeName(deviceType),
		Function:       function,
		Method:         method,
		MethodName:     GetMethodName(method),
		Access:         access,
		AccessName:     GetAccessName(access),
	}

	// Check if this is a known IOCTL
	if name, found := LookupKnownIOCTL(ioctlCode); found {
		components.KnownName = name
	}

	return components
}

// ============================================================================
// IOCTL Formatting Functions
// ============================================================================

// FormatIOCTL formats an IOCTL code into a compact human-readable string.
// This provides a brief, single-line description of the IOCTL.
//
// Format: "DEVICE_TYPE(function:0xXXX, method:METHOD, access:ACCESS)"
// If the IOCTL is known, the known name is prepended.
//
// Parameters:
//   - ioctlCode: The IOCTL code to format
//
// Returns:
//   - A formatted string representation of the IOCTL
//
// Example:
//
//	fmt.Println(FormatIOCTL(IOCTL_DISK_GET_DRIVE_GEOMETRY))
//	// Output: "IOCTL_DISK_GET_DRIVE_GEOMETRY: DISK(function:0x0, method:BUFFERED, access:ANY)"
func FormatIOCTL(ioctlCode uint32) string {
	components := DecodeIOCTL(ioctlCode)

	var result strings.Builder

	// Add known name if available
	if components.KnownName != "" {
		result.WriteString(components.KnownName)
		result.WriteString(": ")
	}

	// Add device type
	result.WriteString(components.DeviceTypeName)
	result.WriteString("(function:0x")
	result.WriteString(fmt.Sprintf("%X", components.Function))
	result.WriteString(", method:")
	result.WriteString(components.MethodName)
	result.WriteString(", access:")
	result.WriteString(components.AccessName)
	result.WriteString(")")

	return result.String()
}

// FormatIOCTLVerbose formats an IOCTL code into a detailed multi-line string.
// This provides complete information about all IOCTL components.
//
// The output includes:
//   - IOCTL code in hexadecimal
//   - Known name (if available)
//   - Device type with code
//   - Function code
//   - Transfer method
//   - Access requirements
//   - Binary breakdown of the IOCTL code
//
// Parameters:
//   - ioctlCode: The IOCTL code to format
//
// Returns:
//   - A detailed multi-line formatted string
//
// Example:
//
//	fmt.Println(FormatIOCTLVerbose(IOCTL_DISK_GET_DRIVE_GEOMETRY))
func FormatIOCTLVerbose(ioctlCode uint32) string {
	components := DecodeIOCTL(ioctlCode)

	var result strings.Builder

	result.WriteString("IOCTL Code: 0x")
	result.WriteString(fmt.Sprintf("%08X", components.IOCTLCode))
	result.WriteString("\n")

	if components.KnownName != "" {
		result.WriteString("Known Name: ")
		result.WriteString(components.KnownName)
		result.WriteString("\n")
	}

	result.WriteString("Device Type: ")
	result.WriteString(components.DeviceTypeName)
	result.WriteString(" (0x")
	result.WriteString(fmt.Sprintf("%04X", components.DeviceType))
	result.WriteString(")\n")

	result.WriteString("Function: 0x")
	result.WriteString(fmt.Sprintf("%03X", components.Function))
	result.WriteString(" (")
	result.WriteString(fmt.Sprintf("%d", components.Function))
	result.WriteString(")\n")

	result.WriteString("Method: ")
	result.WriteString(components.MethodName)
	result.WriteString(" (")
	result.WriteString(fmt.Sprintf("%d", components.Method))
	result.WriteString(")\n")

	result.WriteString("Access: ")
	result.WriteString(components.AccessName)
	result.WriteString(" (")
	result.WriteString(fmt.Sprintf("%d", components.Access))
	result.WriteString(")\n")

	// Add binary breakdown
	result.WriteString("\nBit Layout:\n")
	result.WriteString("  [31:16] Device Type: 0x")
	result.WriteString(fmt.Sprintf("%04X", components.DeviceType))
	result.WriteString(" (")
	result.WriteString(components.DeviceTypeName)
	result.WriteString(")\n")

	result.WriteString("  [15:14] Access:      0x")
	result.WriteString(fmt.Sprintf("%01X", components.Access))
	result.WriteString(" (")
	result.WriteString(components.AccessName)
	result.WriteString(")\n")

	result.WriteString("  [13:2]  Function:    0x")
	result.WriteString(fmt.Sprintf("%03X", components.Function))
	result.WriteString("\n")

	result.WriteString("  [1:0]   Method:      0x")
	result.WriteString(fmt.Sprintf("%01X", components.Method))
	result.WriteString(" (")
	result.WriteString(components.MethodName)
	result.WriteString(")\n")

	return result.String()
}

// FormatIOCTLHex formats an IOCTL code showing its hexadecimal representation
// with bit field annotations.
//
// Format: "0xDDDDAAFF M" where:
//   - DDDD = Device type (16 bits)
//   - AA = Access (2 bits)
//   - FF = Function (12 bits, shown as 3 hex digits)
//   - M = Method (2 bits)
//
// Parameters:
//   - ioctlCode: The IOCTL code to format
//
// Returns:
//   - A hex string with bit field breakdown
//
// Example:
//
//	fmt.Println(FormatIOCTLHex(IOCTL_DISK_GET_DRIVE_GEOMETRY))
//	// Output: "0x00070000 [Device:0x0007 Access:0x0 Func:0x000 Method:0x0]"
func FormatIOCTLHex(ioctlCode uint32) string {
	components := DecodeIOCTL(ioctlCode)

	return fmt.Sprintf("0x%08X [Device:0x%04X Access:0x%X Func:0x%03X Method:0x%X]",
		components.IOCTLCode,
		components.DeviceType,
		components.Access,
		components.Function,
		components.Method)
}

// ============================================================================
// IOCTL Discovery and Probing Functions
// ============================================================================

// Windows error codes not in syscall package
const (
	ERROR_INVALID_FUNCTION     syscall.Errno = 1
	ERROR_INVALID_PARAMETER    syscall.Errno = 87
	ERROR_BAD_LENGTH           syscall.Errno = 24
	ERROR_INVALID_USER_BUFFER  syscall.Errno = 1784
	ERROR_NOT_SUPPORTED        syscall.Errno = 50
	ERROR_CALL_NOT_IMPLEMENTED syscall.Errno = 120
)

// IOCTLProbeResult represents the result of probing an IOCTL code
type IOCTLProbeResult struct {
	Code          uint32 // The IOCTL code that was tested
	Valid         bool   // Whether the IOCTL is accepted by the driver
	ErrorCode     error  // The error returned (if any)
	BytesReturned uint32 // Number of bytes returned
}

// ProbeIOCTL tests whether a specific IOCTL code is valid for a device.
// This sends the IOCTL with minimal buffers and checks if it's accepted.
//
// A valid IOCTL may still return an error (e.g., buffer too small), but
// the error code will indicate the IOCTL was recognized. Invalid IOCTLs
// typically return ERROR_INVALID_FUNCTION or ERROR_NOT_SUPPORTED.
//
// Parameters:
//   - hDevice: Handle to the device
//   - ioctlCode: The IOCTL code to probe
//
// Returns:
//   - IOCTLProbeResult containing the probe results
func ProbeIOCTL(hDevice handle.HANDLE, ioctlCode uint32) IOCTLProbeResult {
	// Small test buffers
	inBuf := make([]byte, 16)
	outBuf := make([]byte, 256)
	var bytesReturned uint32

	result := IOCTLProbeResult{
		Code: ioctlCode,
	}

	// Try the IOCTL
	_, err := DeviceIoControl(
		hDevice,
		ioctlCode,
		unsafe.Pointer(&inBuf[0]),
		uint32(len(inBuf)),
		unsafe.Pointer(&outBuf[0]),
		uint32(len(outBuf)),
		&bytesReturned,
		nil,
	)

	result.BytesReturned = bytesReturned
	result.ErrorCode = err

	// Determine if the IOCTL is valid based on the error
	if err == nil {
		// Success means the IOCTL is definitely valid
		result.Valid = true
	} else {
		// Check the error code
		errno, ok := err.(syscall.Errno)
		if ok {
			// These errors indicate the IOCTL was recognized but failed for other reasons
			switch errno {
			case syscall.ERROR_INSUFFICIENT_BUFFER,
				syscall.ERROR_MORE_DATA,
				ERROR_INVALID_PARAMETER,
				ERROR_BAD_LENGTH,
				ERROR_INVALID_USER_BUFFER,
				syscall.ERROR_ACCESS_DENIED:
				result.Valid = true
			// These errors indicate the IOCTL is not recognized
			case ERROR_INVALID_FUNCTION,
				ERROR_NOT_SUPPORTED,
				ERROR_CALL_NOT_IMPLEMENTED:
				result.Valid = false
			default:
				// For other errors, consider it potentially valid
				result.Valid = true
			}
		}
	}

	return result
}

// ScanIOCTLRange scans a range of IOCTL codes to find valid ones.
// This systematically tests IOCTL codes within the specified range
// and returns those that are accepted by the driver.
//
// WARNING: This can be slow for large ranges. Consider using smaller
// ranges or specific function code ranges for efficiency.
//
// Parameters:
//   - hDevice: Handle to the device
//   - startCode: Starting IOCTL code (inclusive)
//   - endCode: Ending IOCTL code (inclusive)
//
// Returns:
//   - A slice of IOCTLProbeResult for all valid IOCTLs found
func ScanIOCTLRange(hDevice handle.HANDLE, startCode, endCode uint32) []IOCTLProbeResult {
	var results []IOCTLProbeResult

	for code := startCode; code <= endCode; code++ {
		result := ProbeIOCTL(hDevice, code)
		if result.Valid {
			results = append(results, result)
		}
	}

	return results
}

// DiscoverIOCTLsByDeviceType discovers IOCTLs for a device by scanning
// common function code ranges for the device's type.
//
// This is more efficient than scanning all possible codes, as it focuses
// on likely function codes (0-0x100, 0x400-0x600, 0x800-0x900, 0xFFF).
//
// Parameters:
//   - hDevice: Handle to the device
//   - deviceType: The device type code (e.g., FILE_DEVICE_DISK)
//
// Returns:
//   - A slice of IOCTLProbeResult for all valid IOCTLs found
func DiscoverIOCTLsByDeviceType(hDevice handle.HANDLE, deviceType uint32) []IOCTLProbeResult {
	var results []IOCTLProbeResult

	// Common function code ranges to scan
	functionRanges := []struct {
		start uint32
		end   uint32
	}{
		{0x000, 0x100}, // Common low range (0-256)
		{0x400, 0x600}, // Storage commands range
		{0x800, 0x900}, // Extended range
		{0xF00, 0xFFF}, // High range
	}

	// Try all methods and access combinations
	methods := []uint32{METHOD_BUFFERED, METHOD_IN_DIRECT, METHOD_OUT_DIRECT, METHOD_NEITHER}
	accessLevels := []uint32{FILE_ANY_ACCESS, FILE_READ_ACCESS, FILE_WRITE_ACCESS, 3}

	for _, funcRange := range functionRanges {
		for function := funcRange.start; function <= funcRange.end; function++ {
			for _, method := range methods {
				for _, access := range accessLevels {
					code := CTL_CODE(deviceType, function, method, access)
					result := ProbeIOCTL(hDevice, code)
					if result.Valid {
						results = append(results, result)
						// Once we find a valid combination for this function,
						// we can skip other method/access combos for efficiency
						goto nextFunction
					}
				}
			}
		nextFunction:
		}
	}

	return results
}
