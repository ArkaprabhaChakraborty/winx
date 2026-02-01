package service

import (
	"syscall"
	"unsafe"

	"github.com/ArkaprabhaChakraborty/winx/handle"
)

var (
	advapi32                 = syscall.NewLazyDLL("advapi32.dll")
	procOpenSCManagerW       = advapi32.NewProc("OpenSCManagerW")
	procCreateServiceW       = advapi32.NewProc("CreateServiceW")
	procOpenServiceW         = advapi32.NewProc("OpenServiceW")
	procStartServiceW        = advapi32.NewProc("StartServiceW")
	procControlService       = advapi32.NewProc("ControlService")
	procDeleteService        = advapi32.NewProc("DeleteService")
	procCloseServiceHandle   = advapi32.NewProc("CloseServiceHandle")
	procQueryServiceStatus   = advapi32.NewProc("QueryServiceStatus")
)

// Service Control Manager access rights
const (
	SC_MANAGER_CONNECT            = 0x0001
	SC_MANAGER_CREATE_SERVICE     = 0x0002
	SC_MANAGER_ENUMERATE_SERVICE  = 0x0004
	SC_MANAGER_LOCK               = 0x0008
	SC_MANAGER_QUERY_LOCK_STATUS  = 0x0010
	SC_MANAGER_MODIFY_BOOT_CONFIG = 0x0020
	SC_MANAGER_ALL_ACCESS         = 0xF003F
)

// Service access rights
const (
	SERVICE_QUERY_CONFIG         = 0x0001
	SERVICE_CHANGE_CONFIG        = 0x0002
	SERVICE_QUERY_STATUS         = 0x0004
	SERVICE_ENUMERATE_DEPENDENTS = 0x0008
	SERVICE_START                = 0x0010
	SERVICE_STOP                 = 0x0020
	SERVICE_PAUSE_CONTINUE       = 0x0040
	SERVICE_INTERROGATE          = 0x0080
	SERVICE_USER_DEFINED_CONTROL = 0x0100
	SERVICE_ALL_ACCESS           = 0xF01FF
)

// Service types
const (
	SERVICE_KERNEL_DRIVER       = 0x00000001
	SERVICE_FILE_SYSTEM_DRIVER  = 0x00000002
	SERVICE_ADAPTER             = 0x00000004
	SERVICE_RECOGNIZER_DRIVER   = 0x00000008
	SERVICE_WIN32_OWN_PROCESS   = 0x00000010
	SERVICE_WIN32_SHARE_PROCESS = 0x00000020
	SERVICE_INTERACTIVE_PROCESS = 0x00000100
)

// Service start types
const (
	SERVICE_BOOT_START   = 0x00000000
	SERVICE_SYSTEM_START = 0x00000001
	SERVICE_AUTO_START   = 0x00000002
	SERVICE_DEMAND_START = 0x00000003
	SERVICE_DISABLED     = 0x00000004
)

// Service error control
const (
	SERVICE_ERROR_IGNORE   = 0x00000000
	SERVICE_ERROR_NORMAL   = 0x00000001
	SERVICE_ERROR_SEVERE   = 0x00000002
	SERVICE_ERROR_CRITICAL = 0x00000003
)

// Service control codes
const (
	SERVICE_CONTROL_STOP                  = 0x00000001
	SERVICE_CONTROL_PAUSE                 = 0x00000002
	SERVICE_CONTROL_CONTINUE              = 0x00000003
	SERVICE_CONTROL_INTERROGATE           = 0x00000004
	SERVICE_CONTROL_SHUTDOWN              = 0x00000005
	SERVICE_CONTROL_PARAMCHANGE           = 0x00000006
	SERVICE_CONTROL_NETBINDADD            = 0x00000007
	SERVICE_CONTROL_NETBINDREMOVE         = 0x00000008
	SERVICE_CONTROL_NETBINDENABLE         = 0x00000009
	SERVICE_CONTROL_NETBINDDISABLE        = 0x0000000A
	SERVICE_CONTROL_DEVICEEVENT           = 0x0000000B
	SERVICE_CONTROL_HARDWAREPROFILECHANGE = 0x0000000C
	SERVICE_CONTROL_POWEREVENT            = 0x0000000D
	SERVICE_CONTROL_SESSIONCHANGE         = 0x0000000E
)

// Service states
const (
	SERVICE_STOPPED          = 0x00000001
	SERVICE_START_PENDING    = 0x00000002
	SERVICE_STOP_PENDING     = 0x00000003
	SERVICE_RUNNING          = 0x00000004
	SERVICE_CONTINUE_PENDING = 0x00000005
	SERVICE_PAUSE_PENDING    = 0x00000006
	SERVICE_PAUSED           = 0x00000007
)

// SERVICE_STATUS structure
type SERVICE_STATUS struct {
	ServiceType             uint32
	CurrentState            uint32
	ControlsAccepted        uint32
	Win32ExitCode           uint32
	ServiceSpecificExitCode uint32
	CheckPoint              uint32
	WaitHint                uint32
}

// OpenSCManager establishes a connection to the service control manager.
//
// Parameters:
//   - machineName: The name of the target computer (empty string for local)
//   - databaseName: The name of the service control manager database (nil for default)
//   - desiredAccess: The access to the service control manager
//
// Returns:
//   - A handle to the service control manager, or 0 on failure
func OpenSCManager(machineName string, databaseName string, desiredAccess uint32) (handle.HANDLE, error) {
	var machineNamePtr uintptr
	if machineName != "" {
		ptr, err := syscall.UTF16PtrFromString(machineName)
		if err != nil {
			return 0, err
		}
		machineNamePtr = uintptr(unsafe.Pointer(ptr))
	}

	var databaseNamePtr uintptr
	if databaseName != "" {
		ptr, err := syscall.UTF16PtrFromString(databaseName)
		if err != nil {
			return 0, err
		}
		databaseNamePtr = uintptr(unsafe.Pointer(ptr))
	}

	ret, _, _ := syscall.SyscallN(
		procOpenSCManagerW.Addr(),
		machineNamePtr,
		databaseNamePtr,
		uintptr(desiredAccess),
	)

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return handle.HANDLE(ret), nil
}

// CreateService creates a service object and adds it to the service control manager database.
//
// Parameters:
//   - hSCManager: A handle to the service control manager
//   - serviceName: The name of the service
//   - displayName: The display name of the service
//   - desiredAccess: The access to the service
//   - serviceType: The service type
//   - startType: The service start options
//   - errorControl: The severity of the error
//   - binaryPathName: The fully qualified path to the service binary
//
// Returns:
//   - A handle to the service, or 0 on failure
func CreateService(
	hSCManager handle.HANDLE,
	serviceName string,
	displayName string,
	desiredAccess uint32,
	serviceType uint32,
	startType uint32,
	errorControl uint32,
	binaryPathName string,
) (handle.HANDLE, error) {

	serviceNamePtr, err := syscall.UTF16PtrFromString(serviceName)
	if err != nil {
		return 0, err
	}

	displayNamePtr, err := syscall.UTF16PtrFromString(displayName)
	if err != nil {
		return 0, err
	}

	binaryPathPtr, err := syscall.UTF16PtrFromString(binaryPathName)
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.SyscallN(
		procCreateServiceW.Addr(),
		uintptr(hSCManager),
		uintptr(unsafe.Pointer(serviceNamePtr)),
		uintptr(unsafe.Pointer(displayNamePtr)),
		uintptr(desiredAccess),
		uintptr(serviceType),
		uintptr(startType),
		uintptr(errorControl),
		uintptr(unsafe.Pointer(binaryPathPtr)),
		0, // lpLoadOrderGroup
		0, // lpdwTagId
		0, // lpDependencies
		0, // lpServiceStartName
		0, // lpPassword
	)

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return handle.HANDLE(ret), nil
}

// OpenService opens an existing service.
//
// Parameters:
//   - hSCManager: A handle to the service control manager
//   - serviceName: The name of the service
//   - desiredAccess: The access to the service
//
// Returns:
//   - A handle to the service, or 0 on failure
func OpenService(hSCManager handle.HANDLE, serviceName string, desiredAccess uint32) (handle.HANDLE, error) {
	serviceNamePtr, err := syscall.UTF16PtrFromString(serviceName)
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.SyscallN(
		procOpenServiceW.Addr(),
		uintptr(hSCManager),
		uintptr(unsafe.Pointer(serviceNamePtr)),
		uintptr(desiredAccess),
	)

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return handle.HANDLE(ret), nil
}

// StartService starts a service.
//
// Parameters:
//   - hService: A handle to the service
//   - args: Arguments to pass to the service (can be nil)
//
// Returns:
//   - true if successful, false otherwise
func StartService(hService handle.HANDLE, args []string) (bool, error) {
	var numArgs uint32
	var argPtrs uintptr

	if len(args) > 0 {
		numArgs = uint32(len(args))
		// Convert args to UTF-16 pointers
		utf16Args := make([]*uint16, len(args))
		for i, arg := range args {
			ptr, err := syscall.UTF16PtrFromString(arg)
			if err != nil {
				return false, err
			}
			utf16Args[i] = ptr
		}
		argPtrs = uintptr(unsafe.Pointer(&utf16Args[0]))
	}

	ret, _, _ := syscall.SyscallN(
		procStartServiceW.Addr(),
		uintptr(hService),
		uintptr(numArgs),
		argPtrs,
	)

	if ret == 0 {
		return false, syscall.GetLastError()
	}

	return true, nil
}

// ControlService sends a control code to a service.
//
// Parameters:
//   - hService: A handle to the service
//   - control: The control code
//   - serviceStatus: A pointer to SERVICE_STATUS to receive the status
//
// Returns:
//   - true if successful, false otherwise
func ControlService(hService handle.HANDLE, control uint32, serviceStatus *SERVICE_STATUS) (bool, error) {
	ret, _, _ := syscall.SyscallN(
		procControlService.Addr(),
		uintptr(hService),
		uintptr(control),
		uintptr(unsafe.Pointer(serviceStatus)),
	)

	if ret == 0 {
		return false, syscall.GetLastError()
	}

	return true, nil
}

// DeleteService marks a service for deletion from the service control manager database.
//
// Parameters:
//   - hService: A handle to the service
//
// Returns:
//   - true if successful, false otherwise
func DeleteService(hService handle.HANDLE) (bool, error) {
	ret, _, _ := syscall.SyscallN(
		procDeleteService.Addr(),
		uintptr(hService),
	)

	if ret == 0 {
		return false, syscall.GetLastError()
	}

	return true, nil
}

// CloseServiceHandle closes a handle to a service control manager or service.
//
// Parameters:
//   - hSCObject: A handle to the service control manager or service
//
// Returns:
//   - true if successful, false otherwise
func CloseServiceHandle(hSCObject handle.HANDLE) bool {
	ret, _, _ := syscall.SyscallN(
		procCloseServiceHandle.Addr(),
		uintptr(hSCObject),
	)
	return ret != 0
}

// QueryServiceStatus retrieves the current status of the specified service.
//
// Parameters:
//   - hService: A handle to the service
//   - serviceStatus: A pointer to SERVICE_STATUS to receive the status
//
// Returns:
//   - true if successful, false otherwise
func QueryServiceStatus(hService handle.HANDLE, serviceStatus *SERVICE_STATUS) (bool, error) {
	ret, _, _ := syscall.SyscallN(
		procQueryServiceStatus.Addr(),
		uintptr(hService),
		uintptr(unsafe.Pointer(serviceStatus)),
	)

	if ret == 0 {
		return false, syscall.GetLastError()
	}

	return true, nil
}
