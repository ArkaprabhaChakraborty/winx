package handle

import (
	"testing"
)

// TestHANDLE_IsValidHandle tests the IsValidHandle method
func TestHANDLE_IsValidHandle(t *testing.T) {
	tests := []struct {
		name   string
		handle HANDLE
		want   bool
	}{
		{
			name:   "valid handle - arbitrary value",
			handle: HANDLE(0x1234),
			want:   true,
		},
		{
			name:   "valid handle - large value",
			handle: HANDLE(0xFFFFFFFF),
			want:   true,
		},
		{
			name:   "invalid handle - zero",
			handle: HANDLE(0),
			want:   false,
		},
		{
			name:   "invalid handle - INVALID_HANDLE_VALUE (all bits set)",
			handle: ^HANDLE(0),
			want:   false,
		},
		{
			name:   "valid handle - small value",
			handle: HANDLE(1),
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.handle.IsValidHandle(); got != tt.want {
				t.Errorf("HANDLE.IsValidHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHANDLE_IsValid tests the IsValid alias method
func TestHANDLE_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		handle HANDLE
		want   bool
	}{
		{
			name:   "valid handle",
			handle: HANDLE(0x1234),
			want:   true,
		},
		{
			name:   "invalid handle - zero",
			handle: HANDLE(0),
			want:   false,
		},
		{
			name:   "invalid handle - INVALID_HANDLE_VALUE",
			handle: InvalidHandleValue,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.handle.IsValid(); got != tt.want {
				t.Errorf("HANDLE.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// BenchmarkIsValidHandle benchmarks the IsValidHandle method
func BenchmarkIsValidHandle(b *testing.B) {
	handle := HANDLE(0x1234)
	for i := 0; i < b.N; i++ {
		_ = handle.IsValidHandle()
	}
}
