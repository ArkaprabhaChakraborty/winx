package heap

import (
	"syscall"
	"unsafe"

	"github.com/ArkaprabhaChakraborty/winx/handle"
)

// This package contains heap-related Windows API functions and types.

// Common heap flags
const (
	HEAP_NO_SERIALIZE         = 0x00000001
	HEAP_GROWABLE             = 0x00000002
	HEAP_GENERATE_EXCEPTIONS  = 0x00000004
	HEAP_ZERO_MEMORY          = 0x00000008
	HEAP_REALLOC_IN_PLACE_ONLY = 0x00000010
	HEAP_TAIL_CHECKING_ENABLED = 0x00000020
	HEAP_FREE_CHECKING_ENABLED = 0x00000040
	HEAP_DISABLE_COALESCE_ON_FREE = 0x00000080
	HEAP_CREATE_ALIGN_16      = 0x00010000
	HEAP_CREATE_ENABLE_TRACING = 0x00020000
	HEAP_CREATE_ENABLE_EXECUTE = 0x00040000
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procHeapCreate     = kernel32.NewProc("HeapCreate")
	procHeapDestroy    = kernel32.NewProc("HeapDestroy")
	procHeapAlloc      = kernel32.NewProc("HeapAlloc")
	procHeapReAlloc    = kernel32.NewProc("HeapReAlloc")
	procHeapFree       = kernel32.NewProc("HeapFree")
	procHeapSize       = kernel32.NewProc("HeapSize")
	procHeapValidate   = kernel32.NewProc("HeapValidate")
	procGetProcessHeap = kernel32.NewProc("GetProcessHeap")
	procGetProcessHeaps = kernel32.NewProc("GetProcessHeaps")
)

// HeapCreate creates a private heap object that can be used by the calling process.
// The function reserves space in the virtual address space of the process and allocates physical storage for a specified initial portion of this block.
//
// Parameters:
//   - flOptions: The heap allocation options. These options specify special behaviors for the heap.
//   - dwInitialSize: The initial size of the heap, in bytes.
//   - dwMaximumSize: The maximum size of the heap, in bytes. If zero, the heap can grow in size.
//
// Returns:
//   - A handle to the newly created heap if successful, 0 otherwise.
func HeapCreate(flOptions uint32, dwInitialSize uintptr, dwMaximumSize uintptr) handle.HANDLE {
	ret, _, _ := syscall.SyscallN(
		procHeapCreate.Addr(),
		uintptr(flOptions),
		dwInitialSize,
		dwMaximumSize,
	)
	return handle.HANDLE(ret)
}

// HeapDestroy destroys the specified heap object.
// It decommits and releases all the pages of a private heap object, and it invalidates the handle to the heap.
//
// Parameters:
//   - hHeap: A handle to the heap to be destroyed.
//
// Returns:
//   - true if successful, false otherwise.
func HeapDestroy(hHeap handle.HANDLE) bool {
	ret, _, _ := syscall.SyscallN(
		procHeapDestroy.Addr(),
		uintptr(hHeap),
	)
	return ret != 0
}

// HeapAlloc allocates a block of memory from a heap.
// The allocated memory is not movable.
//
// Parameters:
//   - hHeap: A handle to the heap from which the memory will be allocated.
//   - dwFlags: The heap allocation options.
//   - dwBytes: The number of bytes to be allocated.
//
// Returns:
//   - A pointer to the allocated memory block if successful, nil otherwise.
func HeapAlloc(hHeap handle.HANDLE, dwFlags uint32, dwBytes uintptr) unsafe.Pointer {
	ret, _, _ := syscall.SyscallN(
		procHeapAlloc.Addr(),
		uintptr(hHeap),
		uintptr(dwFlags),
		dwBytes,
	)
	if ret == 0 {
		return nil
	}
	return unsafe.Pointer(ret)
}

// HeapReAlloc reallocates a block of memory from a heap.
// This function enables you to resize a memory block and change other memory block properties.
//
// Parameters:
//   - hHeap: A handle to the heap from which the memory is to be reallocated.
//   - dwFlags: The heap reallocation options.
//   - lpMem: A pointer to the block of memory to be reallocated.
//   - dwBytes: The new size of the memory block, in bytes.
//
// Returns:
//   - A pointer to the reallocated memory block if successful, nil otherwise.
func HeapReAlloc(hHeap handle.HANDLE, dwFlags uint32, lpMem unsafe.Pointer, dwBytes uintptr) unsafe.Pointer {
	ret, _, _ := syscall.SyscallN(
		procHeapReAlloc.Addr(),
		uintptr(hHeap),
		uintptr(dwFlags),
		uintptr(lpMem),
		dwBytes,
	)
	if ret == 0 {
		return nil
	}
	return unsafe.Pointer(ret)
}

// HeapFree frees a memory block allocated from a heap by the HeapAlloc or HeapReAlloc function.
//
// Parameters:
//   - hHeap: A handle to the heap whose memory block is to be freed.
//   - dwFlags: The heap free options.
//   - lpMem: A pointer to the memory block to be freed.
//
// Returns:
//   - true if successful, false otherwise.
func HeapFree(hHeap handle.HANDLE, dwFlags uint32, lpMem unsafe.Pointer) bool {
	ret, _, _ := syscall.SyscallN(
		procHeapFree.Addr(),
		uintptr(hHeap),
		uintptr(dwFlags),
		uintptr(lpMem),
	)
	return ret != 0
}

// HeapSize returns the size of a memory block allocated from a heap by the HeapAlloc or HeapReAlloc function.
//
// Parameters:
//   - hHeap: A handle to the heap in which the memory block resides.
//   - dwFlags: The heap size options.
//   - lpMem: A pointer to the memory block whose size the function will obtain.
//
// Returns:
//   - The size of the allocated memory block, in bytes, or ^uintptr(0) on failure.
func HeapSize(hHeap handle.HANDLE, dwFlags uint32, lpMem unsafe.Pointer) uintptr {
	ret, _, _ := syscall.SyscallN(
		procHeapSize.Addr(),
		uintptr(hHeap),
		uintptr(dwFlags),
		uintptr(lpMem),
	)
	return ret
}

// HeapValidate validates the specified heap.
// The function scans all the memory blocks in the heap and verifies that the heap control structures maintained by the heap manager are in a consistent state.
//
// Parameters:
//   - hHeap: A handle to the heap to be validated.
//   - dwFlags: The heap validation options.
//   - lpMem: A pointer to a memory block within the specified heap. If this parameter is nil, the function validates the entire heap.
//
// Returns:
//   - true if the specified heap is valid, false otherwise.
func HeapValidate(hHeap handle.HANDLE, dwFlags uint32, lpMem unsafe.Pointer) bool {
	ret, _, _ := syscall.SyscallN(
		procHeapValidate.Addr(),
		uintptr(hHeap),
		uintptr(dwFlags),
		uintptr(lpMem),
	)
	return ret != 0
}

// GetProcessHeap retrieves a handle to the default heap of the calling process.
// This handle can be used in subsequent calls to heap functions.
//
// Returns:
//   - A handle to the calling process's heap.
func GetProcessHeap() handle.HANDLE {
	ret, _, _ := syscall.SyscallN(
		procGetProcessHeap.Addr(),
	)
	return handle.HANDLE(ret)
}

// GetProcessHeaps returns the number of active heaps and retrieves handles to all of the active heaps for the calling process.
//
// Parameters:
//   - numberOfHeaps: The maximum number of heap handles that can be stored into the buffer.
//   - processHeaps: A slice to receive the heap handles. The slice must be at least numberOfHeaps in length.
//
// Returns:
//   - The number of handles to heaps that are active for the calling process.
func GetProcessHeaps(numberOfHeaps uint32, processHeaps []handle.HANDLE) uint32 {
	var heapArrayPtr uintptr
	if len(processHeaps) > 0 {
		heapArrayPtr = uintptr(unsafe.Pointer(&processHeaps[0]))
	}

	ret, _, _ := syscall.SyscallN(
		procGetProcessHeaps.Addr(),
		uintptr(numberOfHeaps),
		heapArrayPtr,
	)
	return uint32(ret)
}

// HeapWalk walks through all the entries in a heap.
// This function is useful for debugging and analyzing heap usage.
//
// Parameters:
//   - hHeap: A handle to the heap to walk.
//   - entry: A pointer to a PROCESS_HEAP_ENTRY structure.
//
// Returns:
//   - true if the function succeeds, false otherwise.
func HeapWalk(hHeap handle.HANDLE, entry *PROCESS_HEAP_ENTRY) bool {
	var procHeapWalk = kernel32.NewProc("HeapWalk")
	ret, _, _ := syscall.SyscallN(
		procHeapWalk.Addr(),
		uintptr(hHeap),
		uintptr(unsafe.Pointer(entry)),
	)
	return ret != 0
}

// HeapLock locks a heap to prevent other threads from accessing it.
// Use HeapUnlock to unlock the heap.
//
// Parameters:
//   - hHeap: A handle to the heap to lock.
//
// Returns:
//   - true if successful, false otherwise.
func HeapLock(hHeap handle.HANDLE) bool {
	var procHeapLock = kernel32.NewProc("HeapLock")
	ret, _, _ := syscall.SyscallN(
		procHeapLock.Addr(),
		uintptr(hHeap),
	)
	return ret != 0
}

// HeapUnlock unlocks a heap that was locked by HeapLock.
//
// Parameters:
//   - hHeap: A handle to the heap to unlock.
//
// Returns:
//   - true if successful, false otherwise.
func HeapUnlock(hHeap handle.HANDLE) bool {
	var procHeapUnlock = kernel32.NewProc("HeapUnlock")
	ret, _, _ := syscall.SyscallN(
		procHeapUnlock.Addr(),
		uintptr(hHeap),
	)
	return ret != 0
}

// HeapCompact compacts a heap by coalescing adjacent free blocks.
//
// Parameters:
//   - hHeap: A handle to the heap to compact.
//   - dwFlags: Heap compaction options (usually 0).
//
// Returns:
//   - The size of the largest committed free block in the heap, or 0 on failure.
func HeapCompact(hHeap handle.HANDLE, dwFlags uint32) uintptr {
	var procHeapCompact = kernel32.NewProc("HeapCompact")
	ret, _, _ := syscall.SyscallN(
		procHeapCompact.Addr(),
		uintptr(hHeap),
		uintptr(dwFlags),
	)
	return ret
}

// PROCESS_HEAP_ENTRY represents an entry in a heap, used by HeapWalk.
type PROCESS_HEAP_ENTRY struct {
	Data          unsafe.Pointer
	Size          uint32
	Overhead      byte
	RegionIndex   byte
	Flags         uint16
	BlockOrRegion [24]byte // Union of BLOCK and REGION structures
}

// Heap entry flags
const (
	PROCESS_HEAP_REGION             = 0x0001
	PROCESS_HEAP_UNCOMMITTED_RANGE  = 0x0002
	PROCESS_HEAP_ENTRY_BUSY         = 0x0004
	PROCESS_HEAP_ENTRY_MOVEABLE     = 0x0010
	PROCESS_HEAP_ENTRY_DDESHARE     = 0x0020
)
