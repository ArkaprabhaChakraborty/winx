package exitcodes

import "fmt"

// NTStatusCode represents an NT status code with its symbolic name and description
type NTStatusCode struct {
	Code        uint32
	Name        string
	Description string
}

// NTStatusCodeMap contains common NT status codes
// Note: STATUS_SUCCESS and STATUS_WAIT_0 share the same value (0x00000000)
var NTStatusCodeMap = map[uint32]NTStatusCode{
	// Success codes
	0x00000000: {0x00000000, "STATUS_SUCCESS", "The operation completed successfully."},
	0x00000102: {0x00000102, "STATUS_TIMEOUT", "The timeout period has expired."},
	0x00000103: {0x00000103, "STATUS_PENDING", "The operation that was requested is pending completion."},

	// Informational codes
	0x40000000: {0x40000000, "STATUS_OBJECT_NAME_EXISTS", "The specified object name already exists."},
	0x40000006: {0x40000006, "STATUS_NO_MORE_FILES", "No more files were found which match the file specification."},

	// Warning codes
	0x80000001: {0x80000001, "STATUS_GUARD_PAGE_VIOLATION", "Guard page violation."},
	0x80000002: {0x80000002, "STATUS_DATATYPE_MISALIGNMENT", "Datatype misalignment."},
	0x80000005: {0x80000005, "STATUS_BUFFER_OVERFLOW", "The data was too large to fit into the specified buffer."},
	0x80000006: {0x80000006, "STATUS_NO_MORE_ENTRIES", "No more entries are available from an enumeration operation."},

	// Error codes
	0xC0000001: {0xC0000001, "STATUS_UNSUCCESSFUL", "The requested operation was unsuccessful."},
	0xC0000002: {0xC0000002, "STATUS_NOT_IMPLEMENTED", "The requested operation is not implemented."},
	0xC0000003: {0xC0000003, "STATUS_INVALID_INFO_CLASS", "The specified information class is not a valid information class for the specified object."},
	0xC0000004: {0xC0000004, "STATUS_INFO_LENGTH_MISMATCH", "The specified information record length does not match the length required for the specified information class."},
	0xC0000005: {0xC0000005, "STATUS_ACCESS_VIOLATION", "The instruction at 0x%p referenced memory at 0x%p. The memory could not be %s."},
	0xC0000008: {0xC0000008, "STATUS_INVALID_HANDLE", "An invalid HANDLE was specified."},
	0xC000000D: {0xC000000D, "STATUS_INVALID_PARAMETER", "An invalid parameter was passed to a service or function."},
	0xC000000E: {0xC000000E, "STATUS_NO_SUCH_DEVICE", "The specified device does not exist."},
	0xC000000F: {0xC000000F, "STATUS_NO_SUCH_FILE", "The file does not exist."},
	0xC0000010: {0xC0000010, "STATUS_INVALID_DEVICE_REQUEST", "The specified request is not a valid operation for the target device."},
	0xC0000017: {0xC0000017, "STATUS_NO_MEMORY", "Not enough virtual memory or paging file quota is available to complete the specified operation."},
	0xC0000022: {0xC0000022, "STATUS_ACCESS_DENIED", "A process has requested access to an object, but has not been granted those access rights."},
	0xC0000023: {0xC0000023, "STATUS_BUFFER_TOO_SMALL", "The buffer is too small to contain the entry."},
	0xC0000024: {0xC0000024, "STATUS_OBJECT_TYPE_MISMATCH", "There is a mismatch between the type of object required by the requested operation and the type of object that is specified in the request."},
	0xC0000033: {0xC0000033, "STATUS_OBJECT_NAME_INVALID", "The object name is invalid."},
	0xC0000034: {0xC0000034, "STATUS_OBJECT_NAME_NOT_FOUND", "The object name is not found."},
	0xC0000035: {0xC0000035, "STATUS_OBJECT_NAME_COLLISION", "The object name already exists."},
	0xC0000039: {0xC0000039, "STATUS_OBJECT_PATH_INVALID", "The path to the directory is invalid."},
	0xC000003A: {0xC000003A, "STATUS_OBJECT_PATH_NOT_FOUND", "The path does not exist."},
	0xC000003B: {0xC000003B, "STATUS_OBJECT_PATH_SYNTAX_BAD", "The path syntax is bad."},
	0xC0000043: {0xC0000043, "STATUS_SHARING_VIOLATION", "A file cannot be opened because the share access flags are incompatible."},
	0xC0000044: {0xC0000044, "STATUS_QUOTA_EXCEEDED", "The system quota has been exceeded."},
	0xC0000056: {0xC0000056, "STATUS_DELETE_PENDING", "The file cannot be opened because it is in the process of being deleted."},
	0xC0000061: {0xC0000061, "STATUS_PRIVILEGE_NOT_HELD", "A required privilege is not held by the client."},
	0xC000006D: {0xC000006D, "STATUS_LOGON_FAILURE", "The attempted logon is invalid."},
	0xC0000071: {0xC0000071, "STATUS_PASSWORD_EXPIRED", "The user account password has expired."},
	0xC0000072: {0xC0000072, "STATUS_ACCOUNT_DISABLED", "The user account is currently disabled."},
	0xC000007F: {0xC000007F, "STATUS_DISK_FULL", "There is not enough space on the disk."},
	0xC00000BA: {0xC00000BA, "STATUS_FILE_IS_A_DIRECTORY", "The file that was specified as a target is a directory."},
	0xC00000BB: {0xC00000BB, "STATUS_NOT_SUPPORTED", "The request is not supported."},
	0xC0000101: {0xC0000101, "STATUS_DIRECTORY_NOT_EMPTY", "The directory is not empty."},
	0xC0000103: {0xC0000103, "STATUS_NOT_A_DIRECTORY", "A requested opened file is not a directory."},
	0xC000010A: {0xC000010A, "STATUS_PROCESS_IS_TERMINATING", "An attempt was made to access an exiting process."},
	0xC0000120: {0xC0000120, "STATUS_CANCELLED", "The I/O request was canceled."},
	0xC0000121: {0xC0000121, "STATUS_CANNOT_DELETE", "The file cannot be deleted."},
	0xC0000128: {0xC0000128, "STATUS_FILE_INVALID", "The file is not a valid executable file."},
	0xC0000135: {0xC0000135, "STATUS_DLL_NOT_FOUND", "The required DLL was not found."},
	0xC000013A: {0xC000013A, "STATUS_CONTROL_C_EXIT", "The application terminated as a result of a CTRL+C."},
	0xC0000142: {0xC0000142, "STATUS_DLL_INIT_FAILED", "A DLL initialization routine failed."},
	0xC0000185: {0xC0000185, "STATUS_IO_DEVICE_ERROR", "An I/O device error has occurred."},
	0xC0000188: {0xC0000188, "STATUS_LOG_FILE_FULL", "The event log file is full."},
	0xC000019B: {0xC000019B, "STATUS_TOO_MANY_SECRETS", "The maximum number of secrets that may be stored in a single system has been exceeded."},
	0xC0000205: {0xC0000205, "STATUS_INSUFF_SERVER_RESOURCES", "Insufficient system resources exist to complete the API."},
	0xC0000225: {0xC0000225, "STATUS_NOT_FOUND", "The requested object was not found."},
	0xC0000243: {0xC0000243, "STATUS_USER_MAPPED_FILE", "The file cannot be opened because it is being used by another process."},
	0xC00002C5: {0xC00002C5, "STATUS_NOT_SAME_DEVICE", "An attempt was made to move a file or directory to a different device than that on which it was located."},
	0xC0000354: {0xC0000354, "STATUS_INVALID_LOCK_RANGE", "A requested file lock operation cannot be processed due to an invalid byte range."},
}

// GetNTStatusName returns the symbolic name for a given NTSTATUS code
func GetNTStatusName(code uint32) (string, error) {
	if statusCode, exists := NTStatusCodeMap[code]; exists {
		return statusCode.Name, nil
	}
	return "", fmt.Errorf("NTSTATUS code 0x%08X not found", code)
}

// GetNTStatusDescription returns the description for a given NTSTATUS code
func GetNTStatusDescription(code uint32) (string, error) {
	if statusCode, exists := NTStatusCodeMap[code]; exists {
		return statusCode.Description, nil
	}
	return "", fmt.Errorf("NTSTATUS code 0x%08X not found", code)
}

// GetNTStatusCode returns the full NTStatusCode struct for a given code
func GetNTStatusCode(code uint32) (NTStatusCode, error) {
	if statusCode, exists := NTStatusCodeMap[code]; exists {
		return statusCode, nil
	}
	return NTStatusCode{}, fmt.Errorf("NTSTATUS code 0x%08X not found", code)
}

// FormatNTStatus returns a formatted string containing all information about an NTSTATUS code
func FormatNTStatus(code uint32) string {
	if statusCode, exists := NTStatusCodeMap[code]; exists {
		return fmt.Sprintf("[NTSTATUS: 0x%08X] %s: %s", statusCode.Code, statusCode.Name, statusCode.Description)
	}
	return fmt.Sprintf("Unknown NTSTATUS code: 0x%08X", code)
}

// IsNTSuccess checks if the NTSTATUS code represents success (0x00000000)
func IsNTSuccess(code uint32) bool {
	return code == 0
}

// IsNTError checks if the NTSTATUS code is an error (severity 0x3)
func IsNTError(code uint32) bool {
	return (code >> 30) == 0x3
}

// IsNTWarning checks if the NTSTATUS code is a warning (severity 0x2)
func IsNTWarning(code uint32) bool {
	return (code >> 30) == 0x2
}

// IsNTInformational checks if the NTSTATUS code is informational (severity 0x1)
func IsNTInformational(code uint32) bool {
	return (code >> 30) == 0x1
}
