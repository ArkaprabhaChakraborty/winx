package exitcodes

import (
	"strings"
	"testing"
)

func TestGetErrorMessage(t *testing.T) {
	tests := []struct {
		code        uint32
		wantMessage string
		wantErr     bool
	}{
		{0, "The operation completed successfully.", false},
		{2, "The system cannot find the file specified.", false},
		{5, "Access is denied.", false},
		{87, "The parameter is incorrect.", false},
		{9999, "", true}, // Non-existent error code
	}

	for _, tt := range tests {
		got, err := GetErrorMessage(tt.code)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetErrorMessage(%d) error = %v, wantErr %v", tt.code, err, tt.wantErr)
			continue
		}
		if got != tt.wantMessage {
			t.Errorf("GetErrorMessage(%d) = %v, want %v", tt.code, got, tt.wantMessage)
		}
	}
}

func TestGetErrorName(t *testing.T) {
	tests := []struct {
		code     uint32
		wantName string
		wantErr  bool
	}{
		{0, "SUCCESS", false},
		{2, "ERROR_FILE_NOT_FOUND", false},
		{5, "ERROR_ACCESS_DENIED", false},
		{87, "ERROR_INVALID_PARAMETER", false},
		{9999, "", true}, // Non-existent error code
	}

	for _, tt := range tests {
		got, err := GetErrorName(tt.code)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetErrorName(%d) error = %v, wantErr %v", tt.code, err, tt.wantErr)
			continue
		}
		if got != tt.wantName {
			t.Errorf("GetErrorName(%d) = %v, want %v", tt.code, got, tt.wantName)
		}
	}
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		code    uint32
		want    WindowsErrorCode
		wantErr bool
	}{
		{
			5,
			WindowsErrorCode{5, "ERROR_ACCESS_DENIED", "Access is denied."},
			false,
		},
		{
			9999,
			WindowsErrorCode{},
			true,
		},
	}

	for _, tt := range tests {
		got, err := GetErrorCode(tt.code)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetErrorCode(%d) error = %v, wantErr %v", tt.code, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && (got.Code != tt.want.Code || got.Name != tt.want.Name || got.Message != tt.want.Message) {
			t.Errorf("GetErrorCode(%d) = %v, want %v", tt.code, got, tt.want)
		}
	}
}

func TestIsSuccess(t *testing.T) {
	tests := []struct {
		code uint32
		want bool
	}{
		{0, true},
		{1, false},
		{5, false},
		{87, false},
	}

	for _, tt := range tests {
		if got := IsSuccess(tt.code); got != tt.want {
			t.Errorf("IsSuccess(%d) = %v, want %v", tt.code, got, tt.want)
		}
	}
}

func TestFormatError(t *testing.T) {
	tests := []struct {
		code uint32
		want string
	}{
		{0, "[Return Value: 0] SUCCESS: The operation completed successfully."},
		{2, "[Return Value: 2] ERROR_FILE_NOT_FOUND: The system cannot find the file specified."},
		{9999, "Unknown error code: 9999"},
	}

	for _, tt := range tests {
		if got := FormatError(tt.code); got != tt.want {
			t.Errorf("FormatError(%d) = %v, want %v", tt.code, got, tt.want)
		}
	}
}

func TestErrorCodeMapCompleteness(t *testing.T) {
	// Test that all documented error codes are present
	requiredCodes := []uint32{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		50, 87, 109, 122, 183, 258, 487,
	}

	for _, code := range requiredCodes {
		if _, exists := ErrorCodeMap[code]; !exists {
			t.Errorf("Required error code %d is missing from ErrorCodeMap", code)
		}
	}
}

func TestErrorCodeStructure(t *testing.T) {
	// Test that all entries in the map have valid structure
	for code, errCode := range ErrorCodeMap {
		if errCode.Code != code {
			t.Errorf("ErrorCodeMap[%d].Code = %d, want %d", code, errCode.Code, code)
		}
		if errCode.Name == "" {
			t.Errorf("ErrorCodeMap[%d].Name is empty", code)
		}
		if errCode.Message == "" {
			t.Errorf("ErrorCodeMap[%d].Message is empty", code)
		}
		if code != 0 && !strings.HasPrefix(errCode.Name, "ERROR_") && !strings.HasPrefix(errCode.Name, "WAIT_") {
			t.Errorf("ErrorCodeMap[%d].Name = %s, should start with ERROR_ or WAIT_", code, errCode.Name)
		}
	}
}
