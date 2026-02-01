# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**winx** is a Go library extending `golang.org/x/sys/windows` with Windows NT Native API bindings. It provides low-level Windows programming capabilities including NT API access, handle management, device I/O, driver control, and service management.

- **Language:** Go 1.24.0+
- **Platform:** Windows only
- **Dependencies:** Only Go stdlib and `golang.org/x/sys/windows`

## Build & Test Commands

```bash
go test ./...           # Run all tests
go test -v ./...        # Verbose test output
go test -bench=. ./...  # Run benchmarks
go build ./...          # Build all packages
go mod tidy             # Sync dependencies
```

## Architecture

```
winx/
├── Root Package        # Core types: NTSTATUS, UNICODE_STRING, OBJECT_ATTRIBUTES
├── exitcodes/          # Win32 error codes and NTSTATUS code mappings
├── ntdll/              # NT Native API (NtQuerySystemInformation, etc.)
├── handle/             # HANDLE type with validation, system handle tables
├── device/             # Device I/O, IOCTL, driver loading/unloading
├── service/            # Windows Service Control Manager operations
└── heap/               # Heap allocation functions
```

### Key Design Patterns

**Error Handling:**
- `NTSTATUS` type with `IsSuccess()`, `IsError()`, `IsWarning()` methods
- Functions return `(result, error)` tuples
- `NTStatusError` wraps status codes with messages

**API Layering:**
- Low-level: `syscall.SyscallN()` for direct syscalls
- Mid-level: Safe Go wrappers with error handling
- High-level: `*Ex` and `*WithOptions` variants for advanced configuration

**Resource Management:**
- Creator closes resources (clear ownership)
- Use defer for cleanup
- Validate handles before operations with `IsValid()`, `IsValidHandle()`

**Windows API Integration:**
- `syscall.NewLazyDLL()` for dynamic DLL loading
- Proper `unsafe.Pointer` handling for system structures
- NT API functions auto-size buffers with retry logic

## Key Files

- [device/DRIVER_API.md](device/DRIVER_API.md) - Complete driver loading API reference
- [types.go](types.go) - Core types (NTSTATUS, UNICODE_STRING, OBJECT_ATTRIBUTES)
- [ntdll/info.go](ntdll/info.go) - NtQuerySystemInformation with auto-buffering

## Development Notes

- Many operations require Administrator/SYSTEM privileges
- Tests use safe operations (e.g., NUL device) for real Windows API validation
- Table-driven tests are the standard pattern
