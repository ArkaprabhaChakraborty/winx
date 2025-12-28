# winx - Windows Extended Library for Go

A comprehensive Windows API library extending `golang.org/x/sys/windows` with additional native API functions, types, and utilities not available in the standard library.

## Overview

`winx` provides Go bindings for Windows Native API (NT API) functions and additional Windows functionality that isn't covered by the standard `sys/windows` package. This library is designed for advanced Windows programming scenarios including:

- System programming and debugging
- Process and thread manipulation
- Low-level memory management
- Handle enumeration and inspection
- NT Native API access

## Installation

```bash
go get github.com/ArkaprabhaChakraborty/winx
```

## Project Structure

```
winx/
├── types.go              # Common Windows types (NTSTATUS, UNICODE_STRING, etc.)
├── errors.go             # Error handling and NTSTATUS definitions
├── constants.go          # System constants and information classes
│
├── exitcodes/            # Windows error codes and NTSTATUS codes
│   ├── exitcodes.go      # Win32 error code definitions and utilities
│   └── ntstatus.go       # NT status code definitions and utilities
│
├── ntdll/                # NT Native API (ntdll.dll) functions
│   ├── info.go           # NtQuerySystemInformation and related functions
│   ├── info_test.go      # Tests for system information functions
│   └── types.go          # NT API specific types and structures
│
├── handle/               # Handle management
│   ├── handle.go         # HANDLE type and validation methods
│   ├── handle_test.go    # Tests for HANDLE type
│   ├── table.go          # System handle table structures
│   └── table_test.go     # Tests for handle table operations
│
├── heap/                 # Heap management (future expansion)
│   └── heap.go           # Heap-related constants and functions
│
├── internal/             # Internal utilities (not exported)
│   └── syscall/          # Common syscall helpers
│
└── examples/             # Usage examples
```

## Packages

### Root Package (`winx`)

Provides common types and constants:

- `NTSTATUS` - NT status code type with helper methods
- `UNICODE_STRING` - Unicode string structure for NT APIs
- `OBJECT_ATTRIBUTES` - Object attributes for NT APIs
- System Information Classes constants

### `exitcodes`

Windows error code and NTSTATUS code definitions:

```go
import "github.com/ArkaprabhaChakraborty/winx/exitcodes"

// Win32 error codes
msg, err := exitcodes.GetErrorMessage(5)
fmt.Println(msg) // "Access is denied."

// NTSTATUS codes
name, _ := exitcodes.GetNTStatusName(0xC0000005)
fmt.Println(name) // "STATUS_ACCESS_VIOLATION"

// Format error for debugging
fmt.Println(exitcodes.FormatError(5))
// Output: [Return Value: 5] ERROR_ACCESS_DENIED: Access is denied.
```

### `ntdll`

NT Native API functions from ntdll.dll:

```go
import "github.com/ArkaprabhaChakraborty/winx/ntdll"

// Query system information with automatic buffer sizing
buf, status := ntdll.NtQuerySystemInformation(
    winx.SystemHandleInformation,
    0,     // initial size (0 = auto)
    false, // debug mode
)

if status != 0 {
    fmt.Printf("Error: 0x%08X\n", status)
}
```

### `handle`

Handle types and system handle table operations:

```go
import "github.com/ArkaprabhaChakraborty/winx/handle"

// Validate handles
h := handle.HANDLE(0x1234)
if h.IsValid() {
    fmt.Println("Valid handle")
}

// Work with system handle tables
var table handle.SYSTEM_HANDLE_INFORMATION_EX
// ... populate table from NtQuerySystemInformation ...
handles := table.HandlesSlice()
for _, entry := range handles {
    fmt.Printf("PID: %d, Handle: 0x%X\n",
        entry.UniqueProcessId,
        entry.HandleValue)
}
```

### `heap`

Heap management constants and functions (placeholder for future expansion):

```go
import "github.com/ArkaprabhaChakraborty/winx/heap"

// Heap flags available for future heap operations
const flags = heap.HEAP_ZERO_MEMORY | heap.HEAP_GENERATE_EXCEPTIONS
```

## Features

### Error Handling

Comprehensive error code support:
- Win32 error codes (0-499 range)
- NTSTATUS codes with severity checking
- Formatted error messages for debugging

```go
status := winx.STATUS_INFO_LENGTH_MISMATCH
fmt.Printf("Is Error: %v\n", status.IsError())     // true
fmt.Printf("Is Success: %v\n", status.IsSuccess()) // false
```

### Type Safety

Strong typing for Windows handles and status codes:
- `handle.HANDLE` with validation methods
- `winx.NTSTATUS` with severity checking
- Proper unsafe.Pointer handling for system structures

### Automatic Buffer Management

NT API functions handle buffer sizing automatically:

```go
// No need to guess buffer sizes
data, status := ntdll.NtQuerySystemInformation(
    winx.SystemProcessInformation,
    0,    // Will automatically size the buffer
    true, // Enable debug output
)
```

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## Requirements

- Go 1.24.0 or later
- Windows operating system
- Administrator privileges may be required for certain operations

## Safety and Best Practices

1. **Administrator Privileges**: Many NT API functions require elevated privileges
2. **Error Checking**: Always check NTSTATUS return values
3. **Buffer Validation**: Validate buffer sizes and pointers before use
4. **Handle Lifecycle**: Properly close handles when done
5. **Documentation**: Refer to Microsoft's official NT API documentation

## Examples

See the `examples/` directory for complete usage examples:

- Handle enumeration
- Process information retrieval
- System information queries

## Contributing

Contributions are welcome! Please ensure:

1. Code follows Go best practices
2. All tests pass
3. New features include tests
4. Documentation is updated

## License

See [LICENSE](LICENSE) file for details.

## Related Projects

- [golang.org/x/sys/windows](https://pkg.go.dev/golang.org/x/sys/windows) - Standard Windows API bindings
- [Microsoft Windows Driver Kit Documentation](https://docs.microsoft.com/en-us/windows-hardware/drivers/) - Official NT API documentation

## Disclaimer

This library provides access to undocumented and low-level Windows APIs. Use with caution and always test thoroughly. The NT Native API is subject to change between Windows versions without notice.
