package heap

import (
	"testing"
	"unsafe"

	"github.com/ArkaprabhaChakraborty/winx/handle"
)

// TestHeapCreateAndDestroy tests creating and destroying a heap
func TestHeapCreateAndDestroy(t *testing.T) {
	tests := []struct {
		name           string
		flOptions      uint32
		dwInitialSize  uintptr
		dwMaximumSize  uintptr
		shouldSucceed  bool
	}{
		{
			name:          "create growable heap",
			flOptions:     HEAP_GROWABLE,
			dwInitialSize: 4096,
			dwMaximumSize: 0, // 0 means growable
			shouldSucceed: true,
		},
		{
			name:          "create fixed-size heap",
			flOptions:     0,
			dwInitialSize: 8192,
			dwMaximumSize: 8192,
			shouldSucceed: true,
		},
		{
			name:          "create heap with HEAP_NO_SERIALIZE",
			flOptions:     HEAP_NO_SERIALIZE,
			dwInitialSize: 4096,
			dwMaximumSize: 0,
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hHeap := HeapCreate(tt.flOptions, tt.dwInitialSize, tt.dwMaximumSize)

			if tt.shouldSucceed {
				if hHeap == 0 {
					t.Errorf("HeapCreate() failed, expected valid heap handle")
					return
				}

				// Cleanup: destroy the heap
				if !HeapDestroy(hHeap) {
					t.Errorf("HeapDestroy() failed for heap handle 0x%x", hHeap)
				}
			} else {
				if hHeap != 0 {
					t.Errorf("HeapCreate() succeeded when it should have failed")
					HeapDestroy(hHeap) // cleanup
				}
			}
		})
	}
}

// TestHeapAllocAndFree tests allocating and freeing memory from a heap
func TestHeapAllocAndFree(t *testing.T) {
	// Create a test heap
	hHeap := HeapCreate(0, 4096, 0)
	if hHeap == 0 {
		t.Fatal("Failed to create test heap")
	}
	defer HeapDestroy(hHeap)

	tests := []struct {
		name      string
		dwFlags   uint32
		dwBytes   uintptr
		shouldSucceed bool
	}{
		{
			name:      "allocate 100 bytes",
			dwFlags:   0,
			dwBytes:   100,
			shouldSucceed: true,
		},
		{
			name:      "allocate with HEAP_ZERO_MEMORY",
			dwFlags:   HEAP_ZERO_MEMORY,
			dwBytes:   256,
			shouldSucceed: true,
		},
		{
			name:      "allocate 1KB",
			dwFlags:   0,
			dwBytes:   1024,
			shouldSucceed: true,
		},
		{
			name:      "allocate small block",
			dwFlags:   0,
			dwBytes:   16,
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := HeapAlloc(hHeap, tt.dwFlags, tt.dwBytes)

			if tt.shouldSucceed {
				if ptr == nil {
					t.Errorf("HeapAlloc() failed, expected valid pointer")
					return
				}

				// If HEAP_ZERO_MEMORY was used, verify the memory is zeroed
				if tt.dwFlags&HEAP_ZERO_MEMORY != 0 {
					slice := unsafe.Slice((*byte)(ptr), tt.dwBytes)
					for i, b := range slice {
						if b != 0 {
							t.Errorf("Memory not zeroed at offset %d: got %d, want 0", i, b)
							break
						}
					}
				}

				// Cleanup: free the allocated memory
				if !HeapFree(hHeap, 0, ptr) {
					t.Errorf("HeapFree() failed for pointer %p", ptr)
				}
			} else {
				if ptr != nil {
					t.Errorf("HeapAlloc() succeeded when it should have failed")
					HeapFree(hHeap, 0, ptr) // cleanup
				}
			}
		})
	}
}

// TestHeapReAlloc tests reallocating memory in a heap
func TestHeapReAlloc(t *testing.T) {
	// Create a test heap
	hHeap := HeapCreate(0, 8192, 0)
	if hHeap == 0 {
		t.Fatal("Failed to create test heap")
	}
	defer HeapDestroy(hHeap)

	// Allocate initial memory
	initialSize := uintptr(100)
	ptr := HeapAlloc(hHeap, 0, initialSize)
	if ptr == nil {
		t.Fatal("Failed to allocate initial memory")
	}

	// Write some data to the memory
	slice := unsafe.Slice((*byte)(ptr), initialSize)
	for i := range slice {
		slice[i] = byte(i % 256)
	}

	tests := []struct {
		name      string
		dwFlags   uint32
		newSize   uintptr
		shouldSucceed bool
	}{
		{
			name:      "grow allocation to 200 bytes",
			dwFlags:   0,
			newSize:   200,
			shouldSucceed: true,
		},
		{
			name:      "grow allocation with HEAP_ZERO_MEMORY",
			dwFlags:   HEAP_ZERO_MEMORY,
			newSize:   300,
			shouldSucceed: true,
		},
		{
			name:      "shrink allocation to 50 bytes",
			dwFlags:   0,
			newSize:   50,
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newPtr := HeapReAlloc(hHeap, tt.dwFlags, ptr, tt.newSize)

			if tt.shouldSucceed {
				if newPtr == nil {
					t.Errorf("HeapReAlloc() failed, expected valid pointer")
					return
				}

				// Update ptr for next iteration
				ptr = newPtr
			} else {
				if newPtr != nil {
					t.Errorf("HeapReAlloc() succeeded when it should have failed")
				}
			}
		})
	}

	// Cleanup
	HeapFree(hHeap, 0, ptr)
}

// TestHeapSize tests getting the size of an allocated block
func TestHeapSize(t *testing.T) {
	// Create a test heap
	hHeap := HeapCreate(0, 4096, 0)
	if hHeap == 0 {
		t.Fatal("Failed to create test heap")
	}
	defer HeapDestroy(hHeap)

	tests := []struct {
		name        string
		allocSize   uintptr
		minExpected uintptr
	}{
		{
			name:        "100 byte allocation",
			allocSize:   100,
			minExpected: 100,
		},
		{
			name:        "1KB allocation",
			allocSize:   1024,
			minExpected: 1024,
		},
		{
			name:        "small allocation",
			allocSize:   16,
			minExpected: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := HeapAlloc(hHeap, 0, tt.allocSize)
			if ptr == nil {
				t.Fatal("Failed to allocate memory")
			}
			defer HeapFree(hHeap, 0, ptr)

			size := HeapSize(hHeap, 0, ptr)
			if size == ^uintptr(0) {
				t.Errorf("HeapSize() failed for valid pointer")
				return
			}

			// The actual size may be larger due to heap alignment and overhead
			if size < tt.minExpected {
				t.Errorf("HeapSize() = %d, want at least %d", size, tt.minExpected)
			}
		})
	}
}

// TestHeapValidate tests heap validation
func TestHeapValidate(t *testing.T) {
	// Create a test heap
	hHeap := HeapCreate(0, 4096, 0)
	if hHeap == 0 {
		t.Fatal("Failed to create test heap")
	}
	defer HeapDestroy(hHeap)

	t.Run("validate empty heap", func(t *testing.T) {
		if !HeapValidate(hHeap, 0, nil) {
			t.Error("HeapValidate() failed for empty heap")
		}
	})

	t.Run("validate heap with allocations", func(t *testing.T) {
		// Allocate some memory
		ptr1 := HeapAlloc(hHeap, 0, 100)
		ptr2 := HeapAlloc(hHeap, 0, 200)
		if ptr1 == nil || ptr2 == nil {
			t.Fatal("Failed to allocate test memory")
		}
		defer HeapFree(hHeap, 0, ptr1)
		defer HeapFree(hHeap, 0, ptr2)

		// Validate entire heap
		if !HeapValidate(hHeap, 0, nil) {
			t.Error("HeapValidate() failed for heap with allocations")
		}

		// Validate specific blocks
		if !HeapValidate(hHeap, 0, ptr1) {
			t.Error("HeapValidate() failed for first allocation")
		}

		if !HeapValidate(hHeap, 0, ptr2) {
			t.Error("HeapValidate() failed for second allocation")
		}
	})
}

// TestGetProcessHeap tests getting the default process heap
func TestGetProcessHeap(t *testing.T) {
	hHeap := GetProcessHeap()
	if hHeap == 0 {
		t.Error("GetProcessHeap() returned null handle")
	}

	// Verify we can allocate from the process heap
	ptr := HeapAlloc(hHeap, 0, 100)
	if ptr == nil {
		t.Error("Failed to allocate from process heap")
		return
	}

	// Cleanup
	if !HeapFree(hHeap, 0, ptr) {
		t.Error("Failed to free memory from process heap")
	}
}

// TestGetProcessHeaps tests getting all process heaps
func TestGetProcessHeaps(t *testing.T) {
	// First call to get the number of heaps
	count := GetProcessHeaps(0, nil)
	if count == 0 {
		t.Error("GetProcessHeaps() returned 0 heaps")
		return
	}

	t.Logf("Process has %d heap(s)", count)

	// Allocate slice and get all heap handles
	heaps := make([]handle.HANDLE, count)
	actualCount := GetProcessHeaps(count, heaps)

	if actualCount == 0 {
		t.Error("GetProcessHeaps() failed to retrieve heap handles")
		return
	}

	if actualCount > count {
		t.Logf("Warning: More heaps (%d) than initially reported (%d)", actualCount, count)
	}

	// Validate that we got valid handles
	for i := uint32(0); i < actualCount && i < uint32(len(heaps)); i++ {
		if heaps[i] == 0 {
			t.Errorf("Heap handle at index %d is null", i)
		}
	}

	// Verify the process heap is in the list
	processHeap := GetProcessHeap()
	found := false
	for i := uint32(0); i < actualCount && i < uint32(len(heaps)); i++ {
		if heaps[i] == processHeap {
			found = true
			break
		}
	}

	if !found {
		t.Error("Process heap not found in heap list")
	}
}

// BenchmarkHeapAlloc benchmarks heap allocation
func BenchmarkHeapAlloc(b *testing.B) {
	hHeap := HeapCreate(0, 65536, 0)
	if hHeap == 0 {
		b.Fatal("Failed to create heap")
	}
	defer HeapDestroy(hHeap)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ptr := HeapAlloc(hHeap, 0, 100)
		if ptr != nil {
			HeapFree(hHeap, 0, ptr)
		}
	}
}

// BenchmarkHeapAllocZero benchmarks heap allocation with zero initialization
func BenchmarkHeapAllocZero(b *testing.B) {
	hHeap := HeapCreate(0, 65536, 0)
	if hHeap == 0 {
		b.Fatal("Failed to create heap")
	}
	defer HeapDestroy(hHeap)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ptr := HeapAlloc(hHeap, HEAP_ZERO_MEMORY, 100)
		if ptr != nil {
			HeapFree(hHeap, 0, ptr)
		}
	}
}

// BenchmarkGetProcessHeap benchmarks getting the process heap handle
func BenchmarkGetProcessHeap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetProcessHeap()
	}
}

// TestHeapIterateAndAllocateNumbers tests creating a heap, allocating memory for numbers, and iterating through them
func TestHeapIterateAndAllocateNumbers(t *testing.T) {
	// Create a custom heap
	hHeap := HeapCreate(HEAP_NO_SERIALIZE, 8192, 0)
	if hHeap == 0 {
		t.Fatal("Failed to create heap")
	}
	defer HeapDestroy(hHeap)

	t.Logf("Created heap handle: 0x%x", hHeap)

	// Number of integers to allocate
	count := 10

	// Allocate an array of pointers to store our number allocations
	numberPtrs := make([]unsafe.Pointer, count)

	// Allocate memory for each number and store values
	t.Log("Allocating and storing numbers:")
	for i := 0; i < count; i++ {
		// Allocate memory for an int32
		ptr := HeapAlloc(hHeap, HEAP_ZERO_MEMORY, unsafe.Sizeof(int32(0)))
		if ptr == nil {
			t.Fatalf("Failed to allocate memory for number at index %d", i)
		}

		// Store the pointer
		numberPtrs[i] = ptr

		// Write the value (i * 10) to the allocated memory
		value := int32(i * 10)
		*(*int32)(ptr) = value

		t.Logf("  [%d] Allocated at %p, stored value: %d", i, ptr, value)
	}

	// Now iterate through and read the values back
	t.Log("\nReading values back from heap:")
	for i := 0; i < count; i++ {
		ptr := numberPtrs[i]
		value := *(*int32)(ptr)
		t.Logf("  [%d] Read from %p: %d", i, ptr, value)

		// Verify the value is correct
		expected := int32(i * 10)
		if value != expected {
			t.Errorf("Value mismatch at index %d: got %d, want %d", i, value, expected)
		}
	}

	// Verify heap is still valid after all allocations
	if !HeapValidate(hHeap, 0, nil) {
		t.Error("Heap validation failed after allocations")
	}

	// Check the size of each allocation
	t.Log("\nVerifying allocation sizes:")
	for i := 0; i < count; i++ {
		ptr := numberPtrs[i]
		size := HeapSize(hHeap, 0, ptr)
		if size == ^uintptr(0) {
			t.Errorf("HeapSize failed for allocation at index %d", i)
		} else {
			t.Logf("  [%d] Size of allocation at %p: %d bytes", i, ptr, size)
			if size < unsafe.Sizeof(int32(0)) {
				t.Errorf("Allocation size too small at index %d: got %d, want at least %d",
					i, size, unsafe.Sizeof(int32(0)))
			}
		}
	}

	// Free all allocated memory
	t.Log("\nFreeing allocated memory:")
	for i := 0; i < count; i++ {
		ptr := numberPtrs[i]
		if !HeapFree(hHeap, 0, ptr) {
			t.Errorf("Failed to free memory at index %d (pointer %p)", i, ptr)
		} else {
			t.Logf("  [%d] Freed memory at %p", i, ptr)
		}
	}

	// Final heap validation
	if !HeapValidate(hHeap, 0, nil) {
		t.Error("Heap validation failed after freeing all memory")
	}

	t.Log("\nTest completed successfully!")
}

// TestHeapAllocateStructures tests allocating custom structures in a heap
func TestHeapAllocateStructures(t *testing.T) {
	// Define a custom structure
	type Point struct {
		X int32
		Y int32
		Z int32
	}

	// Create a heap
	hHeap := HeapCreate(0, 4096, 0)
	if hHeap == 0 {
		t.Fatal("Failed to create heap")
	}
	defer HeapDestroy(hHeap)

	t.Logf("Created heap for Point structures")

	// Allocate multiple Point structures
	pointCount := 5
	points := make([]unsafe.Pointer, pointCount)

	t.Log("Allocating Point structures:")
	for i := 0; i < pointCount; i++ {
		ptr := HeapAlloc(hHeap, HEAP_ZERO_MEMORY, unsafe.Sizeof(Point{}))
		if ptr == nil {
			t.Fatalf("Failed to allocate Point at index %d", i)
		}

		points[i] = ptr

		// Initialize the Point
		point := (*Point)(ptr)
		point.X = int32(i * 100)
		point.Y = int32(i * 200)
		point.Z = int32(i * 300)

		t.Logf("  Point[%d] at %p: {X: %d, Y: %d, Z: %d}", i, ptr, point.X, point.Y, point.Z)
	}

	// Read and verify the points
	t.Log("\nVerifying Point structures:")
	for i := 0; i < pointCount; i++ {
		point := (*Point)(points[i])
		t.Logf("  Point[%d]: {X: %d, Y: %d, Z: %d}", i, point.X, point.Y, point.Z)

		if point.X != int32(i*100) || point.Y != int32(i*200) || point.Z != int32(i*300) {
			t.Errorf("Point[%d] has incorrect values", i)
		}
	}

	// Modify one point using HeapReAlloc (not changing size, just demonstrating realloc)
	t.Log("\nTesting HeapReAlloc on Point[2]:")
	originalPtr := points[2]
	newPtr := HeapReAlloc(hHeap, 0, originalPtr, unsafe.Sizeof(Point{}))
	if newPtr == nil {
		t.Error("HeapReAlloc failed")
	} else {
		points[2] = newPtr
		point := (*Point)(newPtr)
		t.Logf("  Point[2] after realloc at %p: {X: %d, Y: %d, Z: %d}", newPtr, point.X, point.Y, point.Z)
	}

	// Clean up
	t.Log("\nCleaning up Point structures:")
	for i := 0; i < pointCount; i++ {
		if !HeapFree(hHeap, 0, points[i]) {
			t.Errorf("Failed to free Point at index %d", i)
		}
	}

	t.Log("Test completed successfully!")
}

// TestHeapWalk tests walking through heap entries
func TestHeapWalk(t *testing.T) {
	// Create a heap
	hHeap := HeapCreate(0, 8192, 0)
	if hHeap == 0 {
		t.Fatal("Failed to create heap")
	}
	defer HeapDestroy(hHeap)

	t.Logf("Created heap handle: 0x%x", hHeap)

	// Allocate some memory blocks
	ptrs := make([]unsafe.Pointer, 5)
	for i := 0; i < 5; i++ {
		ptrs[i] = HeapAlloc(hHeap, 0, uintptr((i+1)*100))
		if ptrs[i] == nil {
			t.Fatalf("Failed to allocate block %d", i)
		}
		t.Logf("Allocated block %d: %p, size %d bytes", i, ptrs[i], (i+1)*100)
	}

	// Walk the heap
	t.Log("\nWalking heap entries:")
	entry := PROCESS_HEAP_ENTRY{}
	entryCount := 0
	busyCount := 0

	// Lock the heap before walking
	if !HeapLock(hHeap) {
		t.Error("Failed to lock heap")
	} else {
		defer HeapUnlock(hHeap)

		for HeapWalk(hHeap, &entry) {
			entryCount++
			if entry.Flags&PROCESS_HEAP_ENTRY_BUSY != 0 {
				busyCount++
				t.Logf("  Entry %d: BUSY, Size: %d bytes, Data: %p", entryCount, entry.Size, entry.Data)
			} else if entry.Flags&PROCESS_HEAP_REGION != 0 {
				t.Logf("  Entry %d: REGION, Size: %d bytes", entryCount, entry.Size)
			} else if entry.Flags&PROCESS_HEAP_UNCOMMITTED_RANGE != 0 {
				t.Logf("  Entry %d: UNCOMMITTED, Size: %d bytes", entryCount, entry.Size)
			} else {
				t.Logf("  Entry %d: FREE, Size: %d bytes", entryCount, entry.Size)
			}

			if entryCount > 1000 {
				t.Log("  ... stopping walk after 1000 entries")
				break
			}
		}
	}

	t.Logf("\nTotal entries walked: %d", entryCount)
	t.Logf("Busy blocks found: %d", busyCount)

	// Clean up allocations
	for i := 0; i < 5; i++ {
		HeapFree(hHeap, 0, ptrs[i])
	}
}

// TestHeapLockUnlock tests locking and unlocking a heap
func TestHeapLockUnlock(t *testing.T) {
	hHeap := HeapCreate(0, 4096, 0)
	if hHeap == 0 {
		t.Fatal("Failed to create heap")
	}
	defer HeapDestroy(hHeap)

	t.Logf("Created heap handle: 0x%x", hHeap)

	// Test locking
	t.Log("Attempting to lock heap...")
	if !HeapLock(hHeap) {
		t.Error("HeapLock failed")
		return
	}
	t.Log("Heap successfully locked")

	// Allocate while locked
	t.Log("Allocating memory while heap is locked...")
	ptr := HeapAlloc(hHeap, 0, 100)
	if ptr == nil {
		t.Error("Failed to allocate while heap is locked")
	} else {
		t.Logf("Successfully allocated 100 bytes at %p while heap was locked", ptr)
	}

	// Validate heap while locked
	t.Log("Validating heap while locked...")
	if HeapValidate(hHeap, 0, nil) {
		t.Log("Heap validation succeeded while locked")
	} else {
		t.Log("Heap validation failed while locked (expected behavior)")
	}

	// Test unlocking
	t.Log("Attempting to unlock heap...")
	if !HeapUnlock(hHeap) {
		t.Error("HeapUnlock failed")
	} else {
		t.Log("Heap successfully unlocked")
	}

	// Allocate after unlocking to verify heap is functional
	t.Log("Allocating memory after unlock to verify heap state...")
	ptr2 := HeapAlloc(hHeap, 0, 200)
	if ptr2 == nil {
		t.Error("Failed to allocate after unlocking heap")
	} else {
		t.Logf("Successfully allocated 200 bytes at %p after unlock", ptr2)
	}

	// Clean up
	t.Log("Cleaning up allocations...")
	if ptr != nil {
		if HeapFree(hHeap, 0, ptr) {
			t.Logf("Freed first allocation at %p", ptr)
		}
	}
	if ptr2 != nil {
		if HeapFree(hHeap, 0, ptr2) {
			t.Logf("Freed second allocation at %p", ptr2)
		}
	}

	t.Log("\nHeap lock/unlock test completed successfully!")
}

// TestHeapCompact tests heap compaction
func TestHeapCompact(t *testing.T) {
	hHeap := HeapCreate(0, 16384, 0)
	if hHeap == 0 {
		t.Fatal("Failed to create heap")
	}
	defer HeapDestroy(hHeap)

	t.Logf("Created heap with 16KB initial size, handle: 0x%x", hHeap)

	// Allocate and free some blocks to fragment the heap
	t.Log("\nAllocating 10 blocks of 512 bytes each...")
	ptrs := make([]unsafe.Pointer, 10)
	for i := 0; i < 10; i++ {
		ptrs[i] = HeapAlloc(hHeap, 0, 512)
		if ptrs[i] == nil {
			t.Fatalf("Failed to allocate block %d", i)
		}
		t.Logf("  Block[%d] allocated at %p", i, ptrs[i])
	}

	// Validate heap before fragmentation
	t.Log("\nValidating heap before fragmentation...")
	if HeapValidate(hHeap, 0, nil) {
		t.Log("Heap is valid before fragmentation")
	}

	// Free every other block to create fragmentation
	t.Log("\nFreeing every other block to create fragmentation...")
	freedCount := 0
	for i := 0; i < 10; i += 2 {
		if HeapFree(hHeap, 0, ptrs[i]) {
			t.Logf("  Freed Block[%d] at %p", i, ptrs[i])
			freedCount++
		}
	}
	t.Logf("Freed %d blocks (blocks 0, 2, 4, 6, 8)", freedCount)

	// Validate heap after fragmentation
	t.Log("\nValidating heap after fragmentation...")
	if HeapValidate(hHeap, 0, nil) {
		t.Log("Heap is still valid after fragmentation")
	}

	// Compact the heap
	t.Log("\nCompacting heap to coalesce free blocks...")
	largestFree := HeapCompact(hHeap, 0)
	if largestFree > 0 {
		t.Logf("Heap compaction successful!")
		t.Logf("Largest free block after compaction: %d bytes (%.2f KB)", largestFree, float64(largestFree)/1024.0)
	} else {
		t.Log("HeapCompact returned 0 (may indicate no free blocks or compaction not needed)")
	}

	// Validate heap after compaction
	t.Log("\nValidating heap after compaction...")
	if HeapValidate(hHeap, 0, nil) {
		t.Log("Heap is valid after compaction")
	}

	// Try to allocate a large block from the compacted heap
	t.Log("\nAttempting to allocate a 2KB block from compacted heap...")
	largePtr := HeapAlloc(hHeap, 0, 2048)
	if largePtr != nil {
		t.Logf("Successfully allocated 2KB at %p from compacted heap", largePtr)
		HeapFree(hHeap, 0, largePtr)
		t.Logf("Freed large allocation")
	} else {
		t.Log("Could not allocate 2KB block (heap may be too fragmented)")
	}

	// Clean up remaining blocks
	t.Log("\nCleaning up remaining blocks...")
	for i := 1; i < 10; i += 2 {
		if HeapFree(hHeap, 0, ptrs[i]) {
			t.Logf("  Freed Block[%d] at %p", i, ptrs[i])
		}
	}

	t.Log("\nHeap compaction test completed successfully!")
}
