package device

import (
	"github.com/ArkaprabhaChakraborty/winx/handle"
	"github.com/ArkaprabhaChakraborty/winx/service"
)

// DriverLoadOptions represents configuration options for loading a driver
type DriverLoadOptions struct {
	// Access rights for the service handle (default: SERVICE_ALL_ACCESS)
	DesiredAccess uint32

	// Service start type (default: SERVICE_DEMAND_START)
	// SERVICE_AUTO_START, SERVICE_BOOT_START, SERVICE_DEMAND_START, SERVICE_DISABLED, SERVICE_SYSTEM_START
	StartType uint32

	// Error control level (default: SERVICE_ERROR_NORMAL)
	// SERVICE_ERROR_IGNORE, SERVICE_ERROR_NORMAL, SERVICE_ERROR_SEVERE, SERVICE_ERROR_CRITICAL
	ErrorControl uint32

	// Whether to start the service immediately after creating it (default: true)
	StartImmediately bool

	// Whether to delete existing service before creating new one (default: false)
	RecreateIfExists bool
}

// DefaultDriverLoadOptions returns the default driver loading options
func DefaultDriverLoadOptions() DriverLoadOptions {
	return DriverLoadOptions{
		DesiredAccess:    service.SERVICE_ALL_ACCESS,
		StartType:        service.SERVICE_DEMAND_START,
		ErrorControl:     service.SERVICE_ERROR_NORMAL,
		StartImmediately: true,
		RecreateIfExists: false,
	}
}

// LoadDriver is a convenience function to load a kernel driver with default options.
//
// Parameters:
//   - driverPath: The full path to the driver .sys file
//   - driverName: The name of the driver service
//
// Returns:
//   - A handle to the driver service, or an error
func LoadDriver(driverPath string, driverName string) (handle.HANDLE, error) {
	// Open service control manager
	scm, err := service.OpenSCManager("", "", service.SC_MANAGER_ALL_ACCESS)
	if err != nil {
		return 0, err
	}
	defer service.CloseServiceHandle(scm)
	// Try to open existing service first
	svc, err := service.OpenService(scm, driverName, service.SERVICE_ALL_ACCESS)
	if err != nil {
		// Service doesn't exist, create it
		svc, err = service.CreateService(
			scm,
			driverName,
			driverName,
			service.SERVICE_ALL_ACCESS,
			service.SERVICE_KERNEL_DRIVER,
			service.SERVICE_DEMAND_START,
			service.SERVICE_ERROR_NORMAL,
			driverPath,
		)
		if err != nil {
			return 0, err
		}
	}

	// Start the service
	_, err = service.StartService(svc, nil)
	if err != nil {
		// Ignore error if service is already running
		var status service.SERVICE_STATUS
		service.QueryServiceStatus(svc, &status)
		if status.CurrentState != service.SERVICE_RUNNING {
			service.CloseServiceHandle(svc)
			return 0, err
		}
	}

	return svc, nil
}

// UnloadDriver is a convenience function to unload a kernel driver.
// This stops the driver and deletes the service.
//
// Parameters:
//   - hService: A handle to the driver service
//
// Returns:
//   - An error if the operation fails
func UnloadDriver(hService handle.HANDLE) error {
	var status service.SERVICE_STATUS

	// Stop the service
	_, err := service.ControlService(hService, service.SERVICE_CONTROL_STOP, &status)
	if err != nil {
		return err
	}

	// Delete the service
	_, err = service.DeleteService(hService)
	if err != nil {
		return err
	}

	service.CloseServiceHandle(hService)
	return nil
}

// LoadDriverEx loads a kernel driver with custom access rights and start type.
//
// Parameters:
//   - driverPath: The full path to the driver .sys file
//   - driverName: The name of the driver service
//   - desiredAccess: Access rights for the service handle (e.g., SERVICE_ALL_ACCESS)
//   - startType: Service start type (e.g., SERVICE_DEMAND_START, SERVICE_AUTO_START)
//   - errorControl: Error control level (e.g., SERVICE_ERROR_NORMAL)
//
// Returns:
//   - A handle to the driver service, or an error
func LoadDriverEx(driverPath string, driverName string, desiredAccess uint32, startType uint32, errorControl uint32) (handle.HANDLE, error) {
	scm, err := service.OpenSCManager("", "", service.SC_MANAGER_ALL_ACCESS)
	if err != nil {
		return 0, err
	}
	defer service.CloseServiceHandle(scm)

	// Try to open existing service
	svc, err := service.OpenService(scm, driverName, desiredAccess)
	if err != nil {
		// Service doesn't exist, create it
		svc, err = service.CreateService(
			scm,
			driverName,
			driverName,
			desiredAccess,
			service.SERVICE_KERNEL_DRIVER,
			startType,
			errorControl,
			driverPath,
		)
		if err != nil {
			return 0, err
		}
	}

	// Start the service
	_, err = service.StartService(svc, nil)
	if err != nil {
		// Ignore error if service is already running
		var status service.SERVICE_STATUS
		service.QueryServiceStatus(svc, &status)
		if status.CurrentState != service.SERVICE_RUNNING {
			service.CloseServiceHandle(svc)
			return 0, err
		}
	}

	return svc, nil
}

// LoadDriverWithOptions loads a kernel driver with detailed configuration options.
//
// Parameters:
//   - driverPath: The full path to the driver .sys file
//   - driverName: The name of the driver service
//   - options: Configuration options for loading the driver
//
// Returns:
//   - A handle to the driver service, or an error
func LoadDriverWithOptions(driverPath string, driverName string, options DriverLoadOptions) (handle.HANDLE, error) {
	scm, err := service.OpenSCManager("", "", service.SC_MANAGER_ALL_ACCESS)
	if err != nil {
		return 0, err
	}
	defer service.CloseServiceHandle(scm)

	// Try to open existing service
	svc, err := service.OpenService(scm, driverName, options.DesiredAccess)
	if err == nil {
		// Service exists
		if options.RecreateIfExists {
			// Stop and delete existing service
			var status service.SERVICE_STATUS
			service.ControlService(svc, service.SERVICE_CONTROL_STOP, &status)
			service.DeleteService(svc)
			service.CloseServiceHandle(svc)

			// Create new service
			svc, err = service.CreateService(
				scm,
				driverName,
				driverName,
				options.DesiredAccess,
				service.SERVICE_KERNEL_DRIVER,
				options.StartType,
				options.ErrorControl,
				driverPath,
			)
			if err != nil {
				return 0, err
			}
		}
	} else {
		// Service doesn't exist, create it
		svc, err = service.CreateService(
			scm,
			driverName,
			driverName,
			options.DesiredAccess,
			service.SERVICE_KERNEL_DRIVER,
			options.StartType,
			options.ErrorControl,
			driverPath,
		)
		if err != nil {
			return 0, err
		}
	}

	// Start the service if requested
	if options.StartImmediately {
		_, err = service.StartService(svc, nil)
		if err != nil {
			// Ignore error if service is already running
			var status service.SERVICE_STATUS
			service.QueryServiceStatus(svc, &status)
			if status.CurrentState != service.SERVICE_RUNNING {
				service.CloseServiceHandle(svc)
				return 0, err
			}
		}
	}

	return svc, nil
}

// StartDriver starts an existing driver service.
//
// Parameters:
//   - hService: A handle to the driver service
//
// Returns:
//   - An error if the operation fails
func StartDriver(hService handle.HANDLE) error {
	_, err := service.StartService(hService, nil)
	if err != nil {
		// Check if already running
		var status service.SERVICE_STATUS
		ok, _ := service.QueryServiceStatus(hService, &status)
		if ok && status.CurrentState == service.SERVICE_RUNNING {
			return nil // Already running, not an error
		}
		return err
	}
	return nil
}

// StopDriver stops a running driver service without deleting it.
//
// Parameters:
//   - hService: A handle to the driver service
//
// Returns:
//   - An error if the operation fails
func StopDriver(hService handle.HANDLE) error {
	var status service.SERVICE_STATUS
	_, err := service.ControlService(hService, service.SERVICE_CONTROL_STOP, &status)
	return err
}

// UnloadDriverEx unloads a kernel driver with options.
//
// Parameters:
//   - hService: A handle to the driver service
//   - deleteService: If true, delete the service; if false, only stop it
//   - closeHandle: If true, close the service handle
//
// Returns:
//   - An error if the operation fails
func UnloadDriverEx(hService handle.HANDLE, deleteService bool, closeHandle bool) error {
	var status service.SERVICE_STATUS

	// Stop the service
	_, err := service.ControlService(hService, service.SERVICE_CONTROL_STOP, &status)
	if err != nil {
		return err
	}

	// Delete the service if requested
	if deleteService {
		_, err = service.DeleteService(hService)
		if err != nil {
			return err
		}
	}

	// Close handle if requested
	if closeHandle {
		service.CloseServiceHandle(hService)
	}

	return nil
}

// OpenExistingDriver opens a handle to an existing driver service.
//
// Parameters:
//   - driverName: The name of the driver service
//   - desiredAccess: Access rights for the service handle
//
// Returns:
//   - A handle to the driver service, or an error
func OpenExistingDriver(driverName string, desiredAccess uint32) (handle.HANDLE, error) {
	scm, err := service.OpenSCManager("", "", service.SC_MANAGER_CONNECT)
	if err != nil {
		return 0, err
	}
	defer service.CloseServiceHandle(scm)

	svc, err := service.OpenService(scm, driverName, desiredAccess)
	if err != nil {
		return 0, err
	}

	return svc, nil
}

// QueryDriverStatus queries the current status of a driver service.
//
// Parameters:
//   - hService: A handle to the driver service
//
// Returns:
//   - The service status structure and an error if the operation fails
func QueryDriverStatus(hService handle.HANDLE) (service.SERVICE_STATUS, error) {
	var status service.SERVICE_STATUS
	ok, err := service.QueryServiceStatus(hService, &status)
	if !ok {
		return status, err
	}
	return status, nil
}
