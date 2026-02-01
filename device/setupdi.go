package device

import (
	"strings"
	"syscall"
	"unsafe"

	"github.com/ArkaprabhaChakraborty/winx/handle"
)

var (
	setupapi                               = syscall.NewLazyDLL("setupapi.dll")
	procSetupDiGetClassDevsW               = setupapi.NewProc("SetupDiGetClassDevsW")
	procSetupDiEnumDeviceInterfaces        = setupapi.NewProc("SetupDiEnumDeviceInterfaces")
	procSetupDiGetDeviceInterfaceDetailW   = setupapi.NewProc("SetupDiGetDeviceInterfaceDetailW")
	procSetupDiDestroyDeviceInfoList       = setupapi.NewProc("SetupDiDestroyDeviceInfoList")
	procSetupDiGetDeviceRegistryPropertyW  = setupapi.NewProc("SetupDiGetDeviceRegistryPropertyW")
	procSetupDiEnumDeviceInfo              = setupapi.NewProc("SetupDiEnumDeviceInfo")
)

// SetupDi flags
const (
	DIGCF_DEFAULT          = 0x00000001
	DIGCF_PRESENT          = 0x00000002
	DIGCF_ALLCLASSES       = 0x00000004
	DIGCF_PROFILE          = 0x00000008
	DIGCF_DEVICEINTERFACE  = 0x00000010
)

// Windows error codes
const (
	ERROR_NO_MORE_ITEMS syscall.Errno = 259
)

// Device registry property codes
const (
	SPDRP_DEVICEDESC                  = 0x00000000
	SPDRP_HARDWAREID                  = 0x00000001
	SPDRP_COMPATIBLEIDS               = 0x00000002
	SPDRP_SERVICE                     = 0x00000004
	SPDRP_CLASS                       = 0x00000007
	SPDRP_CLASSGUID                   = 0x00000008
	SPDRP_DRIVER                      = 0x00000009
	SPDRP_CONFIGFLAGS                 = 0x0000000A
	SPDRP_MFG                         = 0x0000000B
	SPDRP_FRIENDLYNAME                = 0x0000000C
	SPDRP_LOCATION_INFORMATION        = 0x0000000D
	SPDRP_PHYSICAL_DEVICE_OBJECT_NAME = 0x0000000E
	SPDRP_CAPABILITIES                = 0x0000000F
	SPDRP_UI_NUMBER                   = 0x00000010
	SPDRP_UPPERFILTERS                = 0x00000011
	SPDRP_LOWERFILTERS                = 0x00000012
	SPDRP_ENUMERATOR_NAME             = 0x00000016
)

// GUID structure
type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// SP_DEVINFO_DATA structure
type SP_DEVINFO_DATA struct {
	CbSize    uint32
	ClassGuid GUID
	DevInst   uint32
	Reserved  uintptr
}

// SP_DEVICE_INTERFACE_DATA structure
type SP_DEVICE_INTERFACE_DATA struct {
	CbSize             uint32
	InterfaceClassGuid GUID
	Flags              uint32
	Reserved           uintptr
}

// SP_DEVICE_INTERFACE_DETAIL_DATA structure
type SP_DEVICE_INTERFACE_DETAIL_DATA struct {
	CbSize     uint32
	DevicePath [1]uint16 // Variable-length
}

// SetupDiGetClassDevs retrieves a handle to a device information set that
// contains requested device information elements for a local computer.
//
// Parameters:
//   - classGuid: A pointer to the GUID for a device setup class or interface class (can be nil)
//   - enumerator: A string that specifies a PnP enumerator (can be empty)
//   - hwndParent: A handle to the top-level window (usually 0)
//   - flags: Flags that control what is included in the device information set
//
// Returns:
//   - A handle to a device information set, or INVALID_HANDLE_VALUE on failure
func SetupDiGetClassDevs(
	classGuid *GUID,
	enumerator string,
	hwndParent uintptr,
	flags uint32,
) (handle.HANDLE, error) {

	var classGuidPtr uintptr
	if classGuid != nil {
		classGuidPtr = uintptr(unsafe.Pointer(classGuid))
	}

	var enumeratorPtr uintptr
	if enumerator != "" {
		enumPtr, err := syscall.UTF16PtrFromString(enumerator)
		if err != nil {
			return handle.HANDLE(INVALID_HANDLE_VALUE), err
		}
		enumeratorPtr = uintptr(unsafe.Pointer(enumPtr))
	}

	ret, _, _ := syscall.SyscallN(
		procSetupDiGetClassDevsW.Addr(),
		classGuidPtr,
		enumeratorPtr,
		hwndParent,
		uintptr(flags),
	)

	if ret == INVALID_HANDLE_VALUE {
		return handle.HANDLE(INVALID_HANDLE_VALUE), syscall.GetLastError()
	}

	return handle.HANDLE(ret), nil
}

// SetupDiEnumDeviceInterfaces enumerates the device interfaces that are
// contained in a device information set.
//
// Parameters:
//   - deviceInfoSet: A handle to a device information set
//   - deviceInfoData: A pointer to SP_DEVINFO_DATA (can be nil)
//   - interfaceClassGuid: A pointer to a GUID that specifies the interface class
//   - memberIndex: A zero-based index of the device interface to retrieve
//   - deviceInterfaceData: A pointer to SP_DEVICE_INTERFACE_DATA to receive information
//
// Returns:
//   - true if successful, false if there are no more interfaces
func SetupDiEnumDeviceInterfaces(
	deviceInfoSet handle.HANDLE,
	deviceInfoData *SP_DEVINFO_DATA,
	interfaceClassGuid *GUID,
	memberIndex uint32,
	deviceInterfaceData *SP_DEVICE_INTERFACE_DATA,
) (bool, error) {

	var deviceInfoDataPtr uintptr
	if deviceInfoData != nil {
		deviceInfoDataPtr = uintptr(unsafe.Pointer(deviceInfoData))
	}

	// Initialize cbSize field
	deviceInterfaceData.CbSize = uint32(unsafe.Sizeof(*deviceInterfaceData))

	ret, _, _ := syscall.SyscallN(
		procSetupDiEnumDeviceInterfaces.Addr(),
		uintptr(deviceInfoSet),
		deviceInfoDataPtr,
		uintptr(unsafe.Pointer(interfaceClassGuid)),
		uintptr(memberIndex),
		uintptr(unsafe.Pointer(deviceInterfaceData)),
	)

	if ret == 0 {
		err := syscall.GetLastError()
		if err == ERROR_NO_MORE_ITEMS {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// SetupDiGetDeviceInterfaceDetail retrieves details about a particular device interface.
//
// Parameters:
//   - deviceInfoSet: A handle to a device information set
//   - deviceInterfaceData: A pointer to SP_DEVICE_INTERFACE_DATA
//   - deviceInfoData: A pointer to SP_DEVINFO_DATA (can be nil)
//
// Returns:
//   - The device path as a string, and any error
func SetupDiGetDeviceInterfaceDetail(
	deviceInfoSet handle.HANDLE,
	deviceInterfaceData *SP_DEVICE_INTERFACE_DATA,
	deviceInfoData *SP_DEVINFO_DATA,
) (string, error) {

	// First call to get required size
	var requiredSize uint32
	syscall.SyscallN(
		procSetupDiGetDeviceInterfaceDetailW.Addr(),
		uintptr(deviceInfoSet),
		uintptr(unsafe.Pointer(deviceInterfaceData)),
		0,
		0,
		uintptr(unsafe.Pointer(&requiredSize)),
		0,
	)

	// Allocate buffer
	detailData := make([]byte, requiredSize)
	detailDataPtr := (*SP_DEVICE_INTERFACE_DETAIL_DATA)(unsafe.Pointer(&detailData[0]))
	detailDataPtr.CbSize = uint32(unsafe.Sizeof(SP_DEVICE_INTERFACE_DETAIL_DATA{}))

	var deviceInfoDataPtr uintptr
	if deviceInfoData != nil {
		deviceInfoData.CbSize = uint32(unsafe.Sizeof(*deviceInfoData))
		deviceInfoDataPtr = uintptr(unsafe.Pointer(deviceInfoData))
	}

	// Second call to get actual data
	ret, _, _ := syscall.SyscallN(
		procSetupDiGetDeviceInterfaceDetailW.Addr(),
		uintptr(deviceInfoSet),
		uintptr(unsafe.Pointer(deviceInterfaceData)),
		uintptr(unsafe.Pointer(detailDataPtr)),
		uintptr(requiredSize),
		uintptr(unsafe.Pointer(&requiredSize)),
		deviceInfoDataPtr,
	)

	if ret == 0 {
		return "", syscall.GetLastError()
	}

	// Extract device path from the structure
	pathStart := unsafe.Pointer(&detailDataPtr.DevicePath[0])
	devicePath := syscall.UTF16ToString((*[260]uint16)(pathStart)[:])

	return devicePath, nil
}

// SetupDiDestroyDeviceInfoList destroys a device information set and frees all associated memory.
//
// Parameters:
//   - deviceInfoSet: A handle to the device information set to destroy
//
// Returns:
//   - true if successful, false otherwise
func SetupDiDestroyDeviceInfoList(deviceInfoSet handle.HANDLE) bool {
	ret, _, _ := syscall.SyscallN(
		procSetupDiDestroyDeviceInfoList.Addr(),
		uintptr(deviceInfoSet),
	)
	return ret != 0
}

// SetupDiGetDeviceRegistryProperty retrieves a specified device property.
//
// Parameters:
//   - deviceInfoSet: A handle to the device information set
//   - deviceInfoData: A pointer to SP_DEVINFO_DATA
//   - property: The property to retrieve (one of SPDRP_* constants)
//
// Returns:
//   - The property value as a string, and any error
func SetupDiGetDeviceRegistryProperty(
	deviceInfoSet handle.HANDLE,
	deviceInfoData *SP_DEVINFO_DATA,
	property uint32,
) (string, error) {

	// First call to get required size
	var requiredSize uint32
	var regDataType uint32

	syscall.SyscallN(
		procSetupDiGetDeviceRegistryPropertyW.Addr(),
		uintptr(deviceInfoSet),
		uintptr(unsafe.Pointer(deviceInfoData)),
		uintptr(property),
		uintptr(unsafe.Pointer(&regDataType)),
		0,
		0,
		uintptr(unsafe.Pointer(&requiredSize)),
	)

	// Check if buffer size is reasonable
	if requiredSize == 0 || requiredSize > 65536 {
		return "", syscall.GetLastError()
	}

	// Allocate buffer
	buffer := make([]byte, requiredSize)

	// Second call to get actual data
	ret, _, _ := syscall.SyscallN(
		procSetupDiGetDeviceRegistryPropertyW.Addr(),
		uintptr(deviceInfoSet),
		uintptr(unsafe.Pointer(deviceInfoData)),
		uintptr(property),
		uintptr(unsafe.Pointer(&regDataType)),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(requiredSize),
		uintptr(unsafe.Pointer(&requiredSize)),
	)

	if ret == 0 {
		return "", syscall.GetLastError()
	}

	// Convert to string
	if len(buffer) == 0 {
		return "", nil
	}
	utf16Slice := (*[260]uint16)(unsafe.Pointer(&buffer[0]))[:]
	return syscall.UTF16ToString(utf16Slice), nil
}

// SetupDiEnumDeviceInfo enumerates the device information elements in a device information set.
//
// Parameters:
//   - deviceInfoSet: A handle to the device information set
//   - memberIndex: A zero-based index of the device information element to retrieve
//   - deviceInfoData: A pointer to SP_DEVINFO_DATA to receive information
//
// Returns:
//   - true if successful, false if there are no more devices
func SetupDiEnumDeviceInfo(
	deviceInfoSet handle.HANDLE,
	memberIndex uint32,
	deviceInfoData *SP_DEVINFO_DATA,
) (bool, error) {

	deviceInfoData.CbSize = uint32(unsafe.Sizeof(*deviceInfoData))

	ret, _, _ := syscall.SyscallN(
		procSetupDiEnumDeviceInfo.Addr(),
		uintptr(deviceInfoSet),
		uintptr(memberIndex),
		uintptr(unsafe.Pointer(deviceInfoData)),
	)

	if ret == 0 {
		err := syscall.GetLastError()
		if err == ERROR_NO_MORE_ITEMS {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// EnumerateDevices is a convenience function that enumerates all devices of a specific class.
//
// Parameters:
//   - classGuid: A pointer to the GUID for a device interface class (can be nil for all devices)
//   - flags: Flags that control what is included (typically DIGCF_PRESENT | DIGCF_DEVICEINTERFACE)
//
// Returns:
//   - A slice of device paths, and any error
func EnumerateDevices(classGuid *GUID, flags uint32) ([]string, error) {
	deviceInfoSet, err := SetupDiGetClassDevs(classGuid, "", 0, flags)
	if err != nil {
		return nil, err
	}
	defer SetupDiDestroyDeviceInfoList(deviceInfoSet)

	var devices []string
	var deviceInterfaceData SP_DEVICE_INTERFACE_DATA
	index := uint32(0)

	for {
		success, err := SetupDiEnumDeviceInterfaces(
			deviceInfoSet,
			nil,
			classGuid,
			index,
			&deviceInterfaceData,
		)

		if !success {
			if err != nil {
				return nil, err
			}
			break // No more devices
		}

		devicePath, err := SetupDiGetDeviceInterfaceDetail(
			deviceInfoSet,
			&deviceInterfaceData,
			nil,
		)

		if err != nil {
			return nil, err
		}

		devices = append(devices, devicePath)
		index++
	}

	return devices, nil
}

// DeviceInfo represents information about a device
type DeviceInfo struct {
	Description  string
	HardwareID   string
	Service      string
	Class        string
	FriendlyName string
	Enumerator   string
}

// EnumerateDevicesWithInfo enumerates all devices and returns detailed information
//
// Parameters:
//   - flags: Flags that control what is included (typically DIGCF_PRESENT | DIGCF_ALLCLASSES)
//
// Returns:
//   - A slice of DeviceInfo structs, and any error
func EnumerateDevicesWithInfo(flags uint32) ([]DeviceInfo, error) {
	deviceInfoSet, err := SetupDiGetClassDevs(nil, "", 0, flags)
	if err != nil {
		return nil, err
	}
	defer SetupDiDestroyDeviceInfoList(deviceInfoSet)

	var devices []DeviceInfo
	var deviceInfoData SP_DEVINFO_DATA
	index := uint32(0)

	for index < 10000 { // Limit iterations
		success, _ := SetupDiEnumDeviceInfo(deviceInfoSet, index, &deviceInfoData)
		if !success {
			break
		}

		info := DeviceInfo{}

		// Get device description
		if desc, err := SetupDiGetDeviceRegistryProperty(deviceInfoSet, &deviceInfoData, SPDRP_DEVICEDESC); err == nil {
			info.Description = desc
		}

		// Get hardware ID
		if hwid, err := SetupDiGetDeviceRegistryProperty(deviceInfoSet, &deviceInfoData, SPDRP_HARDWAREID); err == nil {
			info.HardwareID = hwid
		}

		// Get service name
		if service, err := SetupDiGetDeviceRegistryProperty(deviceInfoSet, &deviceInfoData, SPDRP_SERVICE); err == nil {
			info.Service = service
		}

		// Get device class
		if class, err := SetupDiGetDeviceRegistryProperty(deviceInfoSet, &deviceInfoData, SPDRP_CLASS); err == nil {
			info.Class = class
		}

		// Get friendly name
		if friendly, err := SetupDiGetDeviceRegistryProperty(deviceInfoSet, &deviceInfoData, SPDRP_FRIENDLYNAME); err == nil {
			info.FriendlyName = friendly
		}

		// Get enumerator name
		if enum, err := SetupDiGetDeviceRegistryProperty(deviceInfoSet, &deviceInfoData, SPDRP_ENUMERATOR_NAME); err == nil {
			info.Enumerator = enum
		}

		devices = append(devices, info)
		index++
	}

	return devices, nil
}

// FindDevicesByService searches for devices by service name (driver name)
//
// Parameters:
//   - serviceName: The service/driver name to search for (case-insensitive)
//
// Returns:
//   - A slice of DeviceInfo structs matching the service name
func FindDevicesByService(serviceName string) ([]DeviceInfo, error) {
	allDevices, err := EnumerateDevicesWithInfo(DIGCF_PRESENT | DIGCF_ALLCLASSES)
	if err != nil {
		return nil, err
	}

	var matchingDevices []DeviceInfo

	for _, device := range allDevices {
		if len(device.Service) > 0 && strings.EqualFold(device.Service, serviceName) {
			matchingDevices = append(matchingDevices, device)
		}
	}

	return matchingDevices, nil
}

// GetDriverDevicePaths generates common device path patterns for a driver service name
//
// Parameters:
//   - serviceName: The service/driver name (e.g., "CLFS", "AFD")
//
// Returns:
//   - A slice of possible device paths to try
func GetDriverDevicePaths(serviceName string) []string {
	paths := []string{
		`\\.\` + serviceName,
		`\\Device\` + serviceName,
	}

	// Add lowercase variant (use simple inline conversion)
	lowerName := ""
	for i := 0; i < len(serviceName); i++ {
		c := serviceName[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		lowerName += string(c)
	}

	if lowerName != serviceName {
		paths = append(paths, `\\.\`+lowerName)
		paths = append(paths, `\\Device\`+lowerName)
	}

	return paths
}

