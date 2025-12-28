package ntdll

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/ArkaprabhaChakraborty/winx/exitcodes"
)

// _NtQuerySystemInformation is the low-level wrapper for NtQuerySystemInformation
func _NtQuerySystemInformation(
	SystemInformationClass uint32,
	SystemInformation unsafe.Pointer,
	SystemInformationLength uint32,
	ReturnLength *uint32,
	debug bool) uint32 {

	var ntdll = syscall.NewLazyDLL("ntdll.dll")
	var procNtQuerySystemInformation = ntdll.NewProc("NtQuerySystemInformation")

	ret_code, err_no, kerr := syscall.SyscallN(
		procNtQuerySystemInformation.Addr(),
		uintptr(SystemInformationClass),
		uintptr(SystemInformation),
		uintptr(SystemInformationLength),
		uintptr(unsafe.Pointer(ReturnLength)),
	)

	if debug {
		fmt.Printf("[DEBUG] === NtQuerySystemInformation Call ===\n")
		fmt.Printf("[DEBUG] Class: 0x%02X (%d)\n", SystemInformationClass, SystemInformationClass)
		fmt.Printf("[DEBUG] Buffer: %p, Length: %d bytes\n", SystemInformation, SystemInformationLength)
		fmt.Printf("[DEBUG] Return Code: 0x%08X (%s)\n", ret_code, exitcodes.FormatError(uint32(ret_code)))
		fmt.Printf("[DEBUG] Errno: %v\n", err_no)
		fmt.Printf("[DEBUG] Kernel Error: %v\n", kerr)
	}

	return uint32(ret_code)
}

// _NtQuerySystemInformationEx is the low-level wrapper for NtQuerySystemInformationEx
func _NtQuerySystemInformationEx(
	SystemInformationClass uint32,
	InputBuffer unsafe.Pointer,
	InputBufferLength uint32,
	SystemInformation unsafe.Pointer,
	SystemInformationLength uint32,
	ReturnLength *uint32,
	debug bool) uint32 {

	var ntdll = syscall.NewLazyDLL("ntdll.dll")
	var procNtQuerySystemInformationEx = ntdll.NewProc("NtQuerySystemInformationEx")

	// Use SyscallN with proper parameter array
	ret_code, _, kerr := syscall.SyscallN(
		procNtQuerySystemInformationEx.Addr(),
		uintptr(SystemInformationClass),
		uintptr(InputBuffer),
		uintptr(InputBufferLength),
		uintptr(SystemInformation),
		uintptr(SystemInformationLength),
		uintptr(unsafe.Pointer(ReturnLength)),
	)

	if debug {
		fmt.Printf("[DEBUG] === NtQuerySystemInformationEx Call ===\n")
		fmt.Printf("[DEBUG] Class: 0x%02X (%d)\n", SystemInformationClass, SystemInformationClass)
		fmt.Printf("[DEBUG] Input Buffer: %p, Length: %d bytes\n", InputBuffer, InputBufferLength)
		fmt.Printf("[DEBUG] Output Buffer: %p, Length: %d bytes\n", SystemInformation, SystemInformationLength)
		fmt.Printf("[DEBUG] Return Code: 0x%08X (%s)\n", ret_code, exitcodes.FormatError(uint32(ret_code)))
		fmt.Printf("[DEBUG] Kernel Error: %v\n", kerr)
		if ReturnLength != nil {
			fmt.Printf("[DEBUG] Return Length: %d bytes\n", *ReturnLength)
		}
	}
	return uint32(ret_code)
}

// NtQuerySystemInformation is a convenience wrapper around _NtQuerySystemInformation
// that automatically allocates and resizes a buffer when STATUS_INFO_LENGTH_MISMATCH
// is returned. It returns the filled byte slice and the NTSTATUS code.
func NtQuerySystemInformation(class uint32, initialSize uint32, debug bool) ([]byte, uint32) {
	var returnLen uint32
	size := initialSize
	if size == 0 {
		size = 65536 // start with 64 KiB
	}

	for attempts := 0; attempts < 8; attempts++ {
		buf := make([]byte, size)
		ret := _NtQuerySystemInformation(class, unsafe.Pointer(&buf[0]), size, &returnLen, debug)
		if ret == 0 {
			// success; trim to actual returned length when sensible
			if returnLen > 0 && returnLen <= uint32(len(buf)) {
				return buf[:returnLen], ret
			}
			return buf, ret
		}

		// STATUS_INFO_LENGTH_MISMATCH -> need larger buffer
		if ret == 0xC0000004 {
			// use returned length if provided and larger
			if returnLen > 0 {
				size = returnLen
			} else {
				// exponential backoff
				size *= 2
			}
			if debug {
				fmt.Printf("[DEBUG] STATUS_INFO_LENGTH_MISMATCH, increasing buffer to %d\n", size)
			}
			continue
		}

		// other error - return immediately
		return nil, ret
	}
	return nil, 0xC0000004 // give up with INFO_LENGTH_MISMATCH
}

// NtQuerySystemInformationEx is a convenience wrapper around _NtQuerySystemInformationEx
// that automatically allocates and resizes a buffer when STATUS_INFO_LENGTH_MISMATCH
// is returned. It returns the filled byte slice and the NTSTATUS code.
func NtQuerySystemInformationEx(
	class uint32,
	processorGroup uint16, // Add processor group parameter
	initialSize uint32,
	debug bool) ([]byte, uint32) {

	var returnLen uint32
	size := initialSize

	if size == 0 {
		size = 65536
	}

	// Create input buffer with processor group information
	inputBuffer := make([]byte, 4) // USHORT + padding
	inputBuffer[0] = byte(processorGroup & 0xFF)
	inputBuffer[1] = byte((processorGroup >> 8) & 0xFF)

	for attempts := 0; attempts < 8; attempts++ {
		buf := make([]byte, size)

		ret := _NtQuerySystemInformationEx(
			class,
			unsafe.Pointer(&inputBuffer[0]), // Provide input buffer
			uint32(len(inputBuffer)),        // Input buffer length
			unsafe.Pointer(&buf[0]),
			size,
			&returnLen,
			debug,
		)

		if ret == 0 {
			if returnLen > 0 && returnLen <= uint32(len(buf)) {
				result := make([]byte, returnLen)
				copy(result, buf[:returnLen])
				return result, ret
			}
			return buf, ret
		}

		if ret == 0xC0000004 {
			if returnLen > 0 {
				size = returnLen
			} else {
				size *= 2
				if size > 16*1024*1024 {
					break
				}
			}
			if debug {
				fmt.Printf("[DEBUG] STATUS_INFO_LENGTH_MISMATCH, increasing buffer to %d bytes\n", size)
			}
			continue
		}

		return nil, ret
	}

	return nil, 0xC0000004
}
