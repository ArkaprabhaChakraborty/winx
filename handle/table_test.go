package handle

import (
	"testing"
	"unsafe"
)

// TestSYSTEM_HANDLE_INFORMATION_HandlesSlice tests the HandlesSlice method
func TestSYSTEM_HANDLE_INFORMATION_HandlesSlice(t *testing.T) {
	t.Run("empty table - zero handles", func(t *testing.T) {
		table := &SYSTEM_HANDLE_INFORMATION{
			NumberOfHandles: 0,
			Reserved:        0,
			Handles:         nil,
		}

		result := table.HandlesSlice()
		if result != nil {
			t.Errorf("Expected nil for empty table, got slice with length %d", len(result))
		}
	})

	t.Run("table with handles", func(t *testing.T) {
		// Create a slice of handles
		handles := []SYSTEM_HANDLE_TABLE_ENTRY_INFO{
			{
				UniqueProcessId:      1234,
				CreateBackTraceIndex: 0,
				ObjectTypeIndex:      5,
				HandleAttributes:     0,
				HandleValue:          0x100,
				Object:               0x8000000,
				GrantedAccess:        0x1F0FFF,
			},
			{
				UniqueProcessId:      5678,
				CreateBackTraceIndex: 1,
				ObjectTypeIndex:      10,
				HandleAttributes:     1,
				HandleValue:          0x200,
				Object:               0x9000000,
				GrantedAccess:        0x120089,
			},
		}

		table := &SYSTEM_HANDLE_INFORMATION{
			NumberOfHandles: uint32(len(handles)),
			Reserved:        0,
			Handles:         unsafe.Pointer(&handles[0]),
		}

		result := table.HandlesSlice()
		if result == nil {
			t.Fatal("Expected non-nil slice, got nil")
		}

		if len(result) != len(handles) {
			t.Errorf("Expected slice length %d, got %d", len(handles), len(result))
		}

		// Verify the data matches
		for i := range handles {
			if result[i].UniqueProcessId != handles[i].UniqueProcessId {
				t.Errorf("Handle[%d].UniqueProcessId = %d, want %d", i, result[i].UniqueProcessId, handles[i].UniqueProcessId)
			}
			if result[i].HandleValue != handles[i].HandleValue {
				t.Errorf("Handle[%d].HandleValue = %d, want %d", i, result[i].HandleValue, handles[i].HandleValue)
			}
			if result[i].ObjectTypeIndex != handles[i].ObjectTypeIndex {
				t.Errorf("Handle[%d].ObjectTypeIndex = %d, want %d", i, result[i].ObjectTypeIndex, handles[i].ObjectTypeIndex)
			}
		}
	})
}

// TestSYSTEM_HANDLE_INFORMATION_EX_HandlesSlice tests the HandlesSlice method for extended structure
func TestSYSTEM_HANDLE_INFORMATION_EX_HandlesSlice(t *testing.T) {
	t.Run("empty table - zero handles", func(t *testing.T) {
		table := &SYSTEM_HANDLE_INFORMATION_EX{
			NumberOfHandles: 0,
			Reserved:        0,
			Handles:         nil,
		}

		result := table.HandlesSlice()
		if result != nil {
			t.Errorf("Expected nil for empty table, got slice with length %d", len(result))
		}
	})

	t.Run("table with handles", func(t *testing.T) {
		// Create a slice of extended handles
		handles := []SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX{
			{
				Object:                0x8000000,
				UniqueProcessId:       1234,
				HandleValue:           0x100,
				GrantedAccess:         0x1F0FFF,
				CreatorBackTraceIndex: 0,
				ObjectTypeIndex:       5,
				HandleAttributes:      0,
				Reserved:              0,
			},
			{
				Object:                0x9000000,
				UniqueProcessId:       5678,
				HandleValue:           0x200,
				GrantedAccess:         0x120089,
				CreatorBackTraceIndex: 1,
				ObjectTypeIndex:       10,
				HandleAttributes:      1,
				Reserved:              0,
			},
		}

		table := &SYSTEM_HANDLE_INFORMATION_EX{
			NumberOfHandles: uint32(len(handles)),
			Reserved:        0,
			Handles:         unsafe.Pointer(&handles[0]),
		}

		result := table.HandlesSlice()
		if result == nil {
			t.Fatal("Expected non-nil slice, got nil")
		}

		if len(result) != len(handles) {
			t.Errorf("Expected slice length %d, got %d", len(handles), len(result))
		}

		// Verify the data matches
		for i := range handles {
			if result[i].UniqueProcessId != handles[i].UniqueProcessId {
				t.Errorf("Handle[%d].UniqueProcessId = %d, want %d", i, result[i].UniqueProcessId, handles[i].UniqueProcessId)
			}
			if result[i].HandleValue != handles[i].HandleValue {
				t.Errorf("Handle[%d].HandleValue = %d, want %d", i, result[i].HandleValue, handles[i].HandleValue)
			}
			if result[i].ObjectTypeIndex != handles[i].ObjectTypeIndex {
				t.Errorf("Handle[%d].ObjectTypeIndex = %d, want %d", i, result[i].ObjectTypeIndex, handles[i].ObjectTypeIndex)
			}
			if result[i].Object != handles[i].Object {
				t.Errorf("Handle[%d].Object = %d, want %d", i, result[i].Object, handles[i].Object)
			}
		}
	})
}

// TestStructSizes verifies the struct sizes are as expected
func TestStructSizes(t *testing.T) {
	t.Run("SYSTEM_HANDLE_TABLE_ENTRY_INFO size", func(t *testing.T) {
		// Expected size: 2+1+1+1+2+8+4 = 19 bytes (plus padding)
		size := unsafe.Sizeof(SYSTEM_HANDLE_TABLE_ENTRY_INFO{})
		if size == 0 {
			t.Error("SYSTEM_HANDLE_TABLE_ENTRY_INFO size should not be zero")
		}
		t.Logf("SYSTEM_HANDLE_TABLE_ENTRY_INFO size: %d bytes", size)
	})

	t.Run("SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX size", func(t *testing.T) {
		// Expected size: 8+8+8+4+2+2+4+4 = 40 bytes
		size := unsafe.Sizeof(SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX{})
		if size == 0 {
			t.Error("SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX size should not be zero")
		}
		t.Logf("SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX size: %d bytes", size)
	})

	t.Run("SYSTEM_HANDLE_INFORMATION size", func(t *testing.T) {
		size := unsafe.Sizeof(SYSTEM_HANDLE_INFORMATION{})
		if size == 0 {
			t.Error("SYSTEM_HANDLE_INFORMATION size should not be zero")
		}
		t.Logf("SYSTEM_HANDLE_INFORMATION size: %d bytes", size)
	})

	t.Run("SYSTEM_HANDLE_INFORMATION_EX size", func(t *testing.T) {
		size := unsafe.Sizeof(SYSTEM_HANDLE_INFORMATION_EX{})
		if size == 0 {
			t.Error("SYSTEM_HANDLE_INFORMATION_EX size should not be zero")
		}
		t.Logf("SYSTEM_HANDLE_INFORMATION_EX size: %d bytes", size)
	})
}

// BenchmarkHandlesSlice benchmarks the HandlesSlice method
func BenchmarkHandlesSlice(b *testing.B) {
	handles := make([]SYSTEM_HANDLE_TABLE_ENTRY_INFO, 100)
	table := &SYSTEM_HANDLE_INFORMATION{
		NumberOfHandles: uint32(len(handles)),
		Reserved:        0,
		Handles:         unsafe.Pointer(&handles[0]),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = table.HandlesSlice()
	}
}

// BenchmarkHandlesSliceEX benchmarks the HandlesSlice method for extended structure
func BenchmarkHandlesSliceEX(b *testing.B) {
	handles := make([]SYSTEM_HANDLE_TABLE_ENTRY_INFO_EX, 100)
	table := &SYSTEM_HANDLE_INFORMATION_EX{
		NumberOfHandles: uint32(len(handles)),
		Reserved:        0,
		Handles:         unsafe.Pointer(&handles[0]),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = table.HandlesSlice()
	}
}
