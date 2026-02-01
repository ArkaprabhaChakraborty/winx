# Driver Loading API Reference

This document describes the various driver loading and management functions available in the `device` package.

## Quick Reference

| Function | Use Case | Parameters |
|----------|----------|------------|
| `LoadDriver` | Simple driver loading with defaults | path, name |
| `LoadDriverEx` | Custom access rights and start type | path, name, access, startType, errorControl |
| `LoadDriverWithOptions` | Full configuration control | path, name, options struct |
| `UnloadDriver` | Stop and delete driver service | handle |
| `UnloadDriverEx` | Unload with options | handle, deleteService, closeHandle |
| `StartDriver` | Start an existing driver | handle |
| `StopDriver` | Stop without deleting | handle |
| `OpenExistingDriver` | Open handle to existing driver | name, access |
| `QueryDriverStatus` | Get driver status | handle |

## Function Details

### LoadDriver

**Simple driver loading with default options.**

```go
func LoadDriver(driverPath string, driverName string) (handle.HANDLE, error)
```

**Example:**
```go
hService, err := LoadDriver(`C:\Windows\System32\drivers\null.sys`, "MyDriver")
if err != nil {
    log.Fatal(err)
}
defer UnloadDriver(hService)
```

**Default behavior:**
- Access: `SERVICE_ALL_ACCESS`
- Start type: `SERVICE_DEMAND_START`
- Error control: `SERVICE_ERROR_NORMAL`
- Starts immediately: Yes
- Recreates if exists: No

---

### LoadDriverEx

**Load driver with custom access rights, start type, and error control.**

```go
func LoadDriverEx(driverPath string, driverName string, desiredAccess uint32,
                  startType uint32, errorControl uint32) (handle.HANDLE, error)
```

**Example:**
```go
hService, err := LoadDriverEx(
    `C:\Windows\System32\drivers\null.sys`,
    "MyDriver",
    service.SERVICE_ALL_ACCESS,
    service.SERVICE_AUTO_START,     // Start automatically on boot
    service.SERVICE_ERROR_CRITICAL, // Critical error handling
)
```

**Access Rights:**
- `SERVICE_ALL_ACCESS` - Full control
- `SERVICE_QUERY_STATUS` - Query status only
- `SERVICE_START` - Start service only
- `SERVICE_STOP` - Stop service only

**Start Types:**
- `SERVICE_AUTO_START` - Start automatically at boot
- `SERVICE_BOOT_START` - Start during boot sequence
- `SERVICE_DEMAND_START` - Start on demand (default)
- `SERVICE_DISABLED` - Disabled
- `SERVICE_SYSTEM_START` - Start during system initialization

**Error Control:**
- `SERVICE_ERROR_IGNORE` - Ignore errors
- `SERVICE_ERROR_NORMAL` - Normal error handling (default)
- `SERVICE_ERROR_SEVERE` - Severe error handling
- `SERVICE_ERROR_CRITICAL` - Critical error handling

---

### LoadDriverWithOptions

**Load driver with detailed configuration using an options structure.**

```go
func LoadDriverWithOptions(driverPath string, driverName string,
                          options DriverLoadOptions) (handle.HANDLE, error)
```

**Options Structure:**
```go
type DriverLoadOptions struct {
    DesiredAccess    uint32  // Access rights for service handle
    StartType        uint32  // Service start type
    ErrorControl     uint32  // Error control level
    StartImmediately bool    // Start driver after creating service
    RecreateIfExists bool    // Delete and recreate if service exists
}
```

**Example:**
```go
options := DriverLoadOptions{
    DesiredAccess:    service.SERVICE_ALL_ACCESS,
    StartType:        service.SERVICE_DEMAND_START,
    ErrorControl:     service.SERVICE_ERROR_NORMAL,
    StartImmediately: false,  // Don't start yet
    RecreateIfExists: true,   // Force recreation
}

hService, err := LoadDriverWithOptions(
    `C:\Windows\System32\drivers\null.sys`,
    "MyDriver",
    options,
)

// Manually start later
if err := StartDriver(hService); err != nil {
    log.Fatal(err)
}
```

**Get default options:**
```go
options := DefaultDriverLoadOptions()
options.StartImmediately = false
hService, err := LoadDriverWithOptions(path, name, options)
```

---

### UnloadDriver

**Stop the driver and delete the service completely.**

```go
func UnloadDriver(hService handle.HANDLE) error
```

**Example:**
```go
if err := UnloadDriver(hService); err != nil {
    log.Printf("Failed to unload: %v", err)
}
```

**What it does:**
1. Stops the driver service
2. Deletes the service from registry
3. Closes the service handle

---

### UnloadDriverEx

**Unload driver with granular control over stop/delete/close operations.**

```go
func UnloadDriverEx(hService handle.HANDLE, deleteService bool,
                   closeHandle bool) error
```

**Example:**
```go
// Stop driver but keep service registered for later use
err := UnloadDriverEx(hService, false, false)

// Later, completely remove it
err = UnloadDriverEx(hService, true, true)
```

**Parameters:**
- `deleteService`: If true, deletes the service from registry
- `closeHandle`: If true, closes the service handle

**Common patterns:**
```go
// Just stop, keep service and handle open
UnloadDriverEx(hService, false, false)

// Stop and delete, keep handle open
UnloadDriverEx(hService, true, false)

// Stop, delete, and close (same as UnloadDriver)
UnloadDriverEx(hService, true, true)
```

---

### StartDriver

**Start an existing driver service.**

```go
func StartDriver(hService handle.HANDLE) error
```

**Example:**
```go
if err := StartDriver(hService); err != nil {
    log.Printf("Failed to start: %v", err)
}
```

**Notes:**
- Returns `nil` if driver is already running
- Can be called multiple times safely

---

### StopDriver

**Stop a running driver without deleting the service.**

```go
func StopDriver(hService handle.HANDLE) error
```

**Example:**
```go
// Stop driver temporarily
if err := StopDriver(hService); err != nil {
    log.Printf("Failed to stop: %v", err)
}

// Can restart later
if err := StartDriver(hService); err != nil {
    log.Printf("Failed to restart: %v", err)
}
```

**Use cases:**
- Temporarily pause driver operation
- Reload driver configuration
- Testing start/stop cycles

---

### OpenExistingDriver

**Open a handle to an already-loaded driver service.**

```go
func OpenExistingDriver(driverName string, desiredAccess uint32) (handle.HANDLE, error)
```

**Example:**
```go
// Open handle to Windows Beep driver
hService, err := OpenExistingDriver("Beep", service.SERVICE_QUERY_STATUS)
if err != nil {
    log.Fatal(err)
}
defer service.CloseServiceHandle(hService)

// Query its status
status, _ := QueryDriverStatus(hService)
fmt.Printf("Beep driver state: %d\n", status.CurrentState)
```

**Common use cases:**
- Query status of system drivers
- Control existing services
- Attach to drivers loaded by other processes

---

### QueryDriverStatus

**Query the current status of a driver service.**

```go
func QueryDriverStatus(hService handle.HANDLE) (service.SERVICE_STATUS, error)
```

**Example:**
```go
status, err := QueryDriverStatus(hService)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("State: %d\n", status.CurrentState)
fmt.Printf("Type: %d\n", status.ServiceType)
fmt.Printf("Controls Accepted: 0x%X\n", status.ControlsAccepted)
```

**Service States:**
- `1` = `SERVICE_STOPPED`
- `2` = `SERVICE_START_PENDING`
- `3` = `SERVICE_STOP_PENDING`
- `4` = `SERVICE_RUNNING`
- `5` = `SERVICE_CONTINUE_PENDING`
- `6` = `SERVICE_PAUSE_PENDING`
- `7` = `SERVICE_PAUSED`

---

## Common Patterns

### Pattern 1: Simple Load and Unload
```go
hService, err := LoadDriver(path, name)
if err != nil {
    return err
}
defer UnloadDriver(hService)

// Use driver...
```

### Pattern 2: Load Without Auto-Start
```go
options := DefaultDriverLoadOptions()
options.StartImmediately = false

hService, err := LoadDriverWithOptions(path, name, options)
// Driver is registered but not started

// Start when needed
StartDriver(hService)
```

### Pattern 3: Persistent Driver (Don't Delete)
```go
hService, err := LoadDriver(path, name)
defer func() {
    StopDriver(hService)  // Just stop, don't delete
    service.CloseServiceHandle(hService)
}()
```

### Pattern 4: Force Recreation
```go
options := DefaultDriverLoadOptions()
options.RecreateIfExists = true

hService, err := LoadDriverWithOptions(path, name, options)
// Old service deleted and recreated
```

### Pattern 5: Monitor Existing System Driver
```go
hService, err := OpenExistingDriver("CLFS", service.SERVICE_QUERY_STATUS)
if err != nil {
    return err
}
defer service.CloseServiceHandle(hService)

status, _ := QueryDriverStatus(hService)
if status.CurrentState == 4 { // SERVICE_RUNNING
    fmt.Println("CLFS driver is running")
}
```

### Pattern 6: Start/Stop Cycles for Testing
```go
hService, err := LoadDriver(path, name)
defer UnloadDriver(hService)

for i := 0; i < 10; i++ {
    StopDriver(hService)
    time.Sleep(100 * time.Millisecond)
    StartDriver(hService)
    time.Sleep(100 * time.Millisecond)
}
```

---

## Error Handling

All functions return errors that can be checked:

```go
hService, err := LoadDriver(path, name)
if err != nil {
    if strings.Contains(err.Error(), "Access is denied") {
        log.Fatal("Need administrator privileges")
    }
    if strings.Contains(err.Error(), "already exists") {
        // Service already registered
    }
    return err
}
```

Common errors:
- "Access is denied" - Need administrator/SYSTEM privileges
- "The service already exists" - Service is already registered
- "The specified service does not exist" - Service not found
- "The service cannot be started" - Driver initialization failed

---

## Requirements

- **Administrator Privileges**: Required for loading/unloading drivers
- **Valid Driver Path**: Must point to a valid `.sys` kernel driver file
- **Unique Service Name**: Service names must be unique across the system

---

## Safety Notes

1. **Always unload drivers you load** - Use `defer UnloadDriver(hService)`
2. **Check driver compatibility** - Only load drivers compatible with your OS version
3. **Handle errors gracefully** - Driver loading can fail for many reasons
4. **Don't load untrusted drivers** - Kernel drivers run with full system privileges
5. **Test in VM first** - Bad drivers can crash the system

---

## Advanced Example: IOCTL Fuzzing Setup

```go
// Load driver for fuzzing
options := DefaultDriverLoadOptions()
options.RecreateIfExists = true

hService, err := LoadDriverWithOptions(driverPath, "FuzzTarget", options)
if err != nil {
    log.Fatal(err)
}
defer UnloadDriver(hService)

// Wait for driver to initialize
time.Sleep(100 * time.Millisecond)

// Open device for IOCTL testing
hDevice, err := OpenDevice(`\\.\FuzzTarget`, GENERIC_READ|GENERIC_WRITE)
if err != nil {
    log.Fatal(err)
}
defer CloseHandle(hDevice)

// Discover valid IOCTLs
results := DiscoverIOCTLsByDeviceType(hDevice, FILE_DEVICE_UNKNOWN)
for _, result := range results {
    fmt.Printf("Valid IOCTL: 0x%08X\n", result.Code)
}
```
