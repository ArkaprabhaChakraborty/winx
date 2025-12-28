package winx

import "unsafe"

// NTSTATUS represents an NT status code returned by NT Native API functions.
type NTSTATUS uint32

// Common NTSTATUS severity codes
const (
	STATUS_SEVERITY_SUCCESS       = 0x0
	STATUS_SEVERITY_INFORMATIONAL = 0x1
	STATUS_SEVERITY_WARNING       = 0x2
	STATUS_SEVERITY_ERROR         = 0x3
)

// IsSuccess returns true if the NTSTATUS indicates success.
func (status NTSTATUS) IsSuccess() bool {
	return status == 0
}

// IsError returns true if the NTSTATUS indicates an error.
func (status NTSTATUS) IsError() bool {
	return (status >> 30) == STATUS_SEVERITY_ERROR
}

// IsWarning returns true if the NTSTATUS indicates a warning.
func (status NTSTATUS) IsWarning() bool {
	return (status >> 30) == STATUS_SEVERITY_WARNING
}

// IsInformational returns true if the NTSTATUS indicates informational status.
func (status NTSTATUS) IsInformational() bool {
	return (status >> 30) == STATUS_SEVERITY_INFORMATIONAL
}

// UNICODE_STRING represents a Unicode string structure used by NT APIs.
type UNICODE_STRING struct {
	Length        uint16
	MaximumLength uint16
	Buffer        *uint16
}

// OBJECT_ATTRIBUTES represents object attributes used in NT APIs.
type OBJECT_ATTRIBUTES struct {
	Length                   uint32
	RootDirectory            uintptr // Using uintptr instead of HANDLE to avoid circular dependency
	ObjectName               *UNICODE_STRING
	Attributes               uint32
	SecurityDescriptor       unsafe.Pointer
	SecurityQualityOfService unsafe.Pointer
}

// Common object attribute flags
const (
	OBJ_INHERIT             = 0x00000002
	OBJ_PERMANENT           = 0x00000010
	OBJ_EXCLUSIVE           = 0x00000020
	OBJ_CASE_INSENSITIVE    = 0x00000040
	OBJ_OPENIF              = 0x00000080
	OBJ_OPENLINK            = 0x00000100
	OBJ_KERNEL_HANDLE       = 0x00000200
	OBJ_FORCE_ACCESS_CHECK  = 0x00000400
	OBJ_VALID_ATTRIBUTES    = 0x000007F2
)
