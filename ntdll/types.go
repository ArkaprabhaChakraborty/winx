package ntdll

import "unsafe"

// SYSTEM_PROCESS_INFORMATION represents process information from NtQuerySystemInformation
type SYSTEM_PROCESS_INFORMATION struct {
	NextEntryOffset              uint32
	NumberOfThreads              uint32
	Reserved1                    [48]byte
	CreateTime                   int64
	UserTime                     int64
	KernelTime                   int64
	ImageName                    string
	BasePriority                 int32
	ProcessId                    uintptr
	InheritedFromProcessId       uintptr
	HandleCount                  uint32
	SessionId                    uint32
	PageDirectoryBase            uintptr
	PeakVirtualSize              uintptr
	VirtualSize                  uintptr
	PeakWorkingSetSize           uint32
	WorkingSetSize               uint32
	QuotaPeakPagedPoolUsage      uint32
	QuotaPagedPoolUsage          uint32
	QuotaPeakNonPagedPoolUsage   uint32
	QuotaNonPagedPoolUsage       uint32
	PagefileUsage                uint32
	PeakPagefileUsage            uint32
	PrivatePageCount             uint32
	ReadOperationCount           int64
	WriteOperationCount          int64
	OtherOperationCount          int64
	ReadTransferCount            int64
	WriteTransferCount           int64
	OtherTransferCount           int64
	Threads                      unsafe.Pointer
}
