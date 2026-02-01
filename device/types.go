package device

import "unsafe"

// OVERLAPPED structure for asynchronous I/O operations
type OVERLAPPED struct {
	Internal     uintptr
	InternalHigh uintptr
	Offset       uint32
	OffsetHigh   uint32
	HEvent       uintptr
}

// SECURITY_ATTRIBUTES structure
type SECURITY_ATTRIBUTES struct {
	Length             uint32
	SecurityDescriptor unsafe.Pointer
	InheritHandle      int32
}

// DISK_GEOMETRY structure
type DISK_GEOMETRY struct {
	Cylinders         int64
	MediaType         uint32
	TracksPerCylinder uint32
	SectorsPerTrack   uint32
	BytesPerSector    uint32
}

// PARTITION_INFORMATION structure
type PARTITION_INFORMATION struct {
	StartingOffset     int64
	PartitionLength    int64
	HiddenSectors      uint32
	PartitionNumber    uint32
	PartitionType      byte
	BootIndicator      byte
	RecognizedPartition byte
	RewritePartition   byte
}

// STORAGE_DEVICE_NUMBER structure
type STORAGE_DEVICE_NUMBER struct {
	DeviceType      uint32
	DeviceNumber    uint32
	PartitionNumber uint32
}

// STORAGE_PROPERTY_ID enumeration
type STORAGE_PROPERTY_ID uint32

const (
	StorageDeviceProperty STORAGE_PROPERTY_ID = iota
	StorageAdapterProperty
	StorageDeviceIdProperty
	StorageDeviceUniqueIdProperty
	StorageDeviceWriteCacheProperty
	StorageMiniportProperty
	StorageAccessAlignmentProperty
	StorageDeviceSeekPenaltyProperty
	StorageDeviceTrimProperty
	StorageDeviceWriteAggregationProperty
)

// STORAGE_QUERY_TYPE enumeration
type STORAGE_QUERY_TYPE uint32

const (
	PropertyStandardQuery STORAGE_QUERY_TYPE = iota
	PropertyExistsQuery
	PropertyMaskQuery
	PropertyQueryMaxDefined
)

// STORAGE_PROPERTY_QUERY structure
type STORAGE_PROPERTY_QUERY struct {
	PropertyId           STORAGE_PROPERTY_ID
	QueryType            STORAGE_QUERY_TYPE
	AdditionalParameters [1]byte
}

// STORAGE_DESCRIPTOR_HEADER structure
type STORAGE_DESCRIPTOR_HEADER struct {
	Version uint32
	Size    uint32
}

// STORAGE_DEVICE_DESCRIPTOR structure
type STORAGE_DEVICE_DESCRIPTOR struct {
	Version               uint32
	Size                  uint32
	DeviceType            byte
	DeviceTypeModifier    byte
	RemovableMedia        byte
	CommandQueueing       byte
	VendorIdOffset        uint32
	ProductIdOffset       uint32
	ProductRevisionOffset uint32
	SerialNumberOffset    uint32
	BusType               uint32
	RawPropertiesLength   uint32
	RawDeviceProperties   [1]byte
}

// VOLUME_DISK_EXTENTS structure
type VOLUME_DISK_EXTENTS struct {
	NumberOfDiskExtents uint32
	Extents             [1]DISK_EXTENT
}

// DISK_EXTENT structure
type DISK_EXTENT struct {
	DiskNumber     uint32
	StartingOffset int64
	ExtentLength   int64
}

// IOCTLComponents represents the decoded components of an IOCTL code.
// IOCTL codes are 32-bit values with the following structure:
//   Bits 31-16: Device Type (e.g., FILE_DEVICE_DISK = 0x07)
//   Bits 15-14: Access Required (0=ANY, 1=READ, 2=WRITE, 3=READ_WRITE)
//   Bits 13-2:  Function Code (0-4095, driver-specific operation)
//   Bits 1-0:   Transfer Method (0=BUFFERED, 1=IN_DIRECT, 2=OUT_DIRECT, 3=NEITHER)
type IOCTLComponents struct {
	IOCTLCode      uint32 // Original IOCTL code
	DeviceType     uint32 // Device type (bits 31-16)
	DeviceTypeName string // Human-readable device type name
	Function       uint32 // Function code (bits 13-2)
	Method         uint32 // Transfer method (bits 1-0)
	MethodName     string // Human-readable method name
	Access         uint32 // Access required (bits 15-14)
	AccessName     string // Human-readable access name
	KnownName      string // Known IOCTL name if available (e.g., "IOCTL_DISK_GET_DRIVE_GEOMETRY")
}

// IOCTLComparison represents a comparison between two IOCTL codes
type IOCTLComparison struct {
	Code1          uint32
	Code2          uint32
	SameDeviceType bool
	SameFunction   bool
	SameMethod     bool
	SameAccess     bool
	Identical      bool
}
