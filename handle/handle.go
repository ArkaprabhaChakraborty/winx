package handle

// HANDLE represents a Windows handle value.
// Common handle values:
//   - 0 (NULL): Invalid handle
//   - ^HANDLE(0) (INVALID_HANDLE_VALUE): Invalid handle returned by some APIs
type HANDLE uintptr

// InvalidHandleValue is the constant representing an invalid handle (-1).
const InvalidHandleValue = ^HANDLE(0)

// IsValidHandle returns true if the handle is valid (non-zero and not INVALID_HANDLE_VALUE).
func (h HANDLE) IsValidHandle() bool {
	return h != 0 && h != InvalidHandleValue
}

// IsValid is an alias for IsValidHandle for consistency with other packages.
func (h HANDLE) IsValid() bool {
	return h.IsValidHandle()
}
