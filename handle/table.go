package handle

import "unsafe"

// SYSTEM_HANDLE_TABLE_ENTRY_INFO represents a single entry in the system handle table
type SYSTEM_HANDLE_TABLE_ENTRY_INFO struct {
	UniqueProcessId      uint16
	CreateBackTraceIndex uint8
	ObjectTypeIndex      uint8
	HandleAttributes     uint8
	HandleValue          uint16
	Object               uintptr
	GrantedAccess        uint32
}

// SYSTEM_HANDLE_INFORMATION represents the system handle information structure
type SYSTEM_HANDLE_INFORMATION struct {
	NumberOfHandles uint32
	Reserved        uint32
	Handles         unsafe.Pointer // SYSTEM_HANDLE_TABLE_ENTRY_INFO[1]
}

// SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX represents an extended entry in the system handle table
type SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX struct {
	Object                uintptr
	UniqueProcessId       uintptr
	HandleValue           uintptr
	GrantedAccess         uint32
	CreatorBackTraceIndex uint16
	ObjectTypeIndex       uint16
	HandleAttributes      uint32
	Reserved              uint32
}

// SYSTEM_HANDLE_INFORMATION_EX represents the extended system handle information structure
type SYSTEM_HANDLE_INFORMATION_EX struct {
	NumberOfHandles uint32
	Reserved        uint32
	Handles         unsafe.Pointer // SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX[1]
}

// HandlesSlice converts the handle table to a Go slice for easier iteration
func (table *SYSTEM_HANDLE_INFORMATION) HandlesSlice() []SYSTEM_HANDLE_TABLE_ENTRY_INFO {
	if table.NumberOfHandles == 0 {
		return nil
	}
	return unsafe.Slice((*SYSTEM_HANDLE_TABLE_ENTRY_INFO)(table.Handles), table.NumberOfHandles)
}

// HandlesSlice converts the extended handle table to a Go slice for easier iteration
func (table *SYSTEM_HANDLE_INFORMATION_EX) HandlesSlice() []SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX {
	if table.NumberOfHandles == 0 {
		return nil
	}
	return unsafe.Slice((*SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX)(table.Handles), table.NumberOfHandles)
}
