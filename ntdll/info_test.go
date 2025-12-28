package ntdll

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/ArkaprabhaChakraborty/winx/exitcodes"
)

func TestNtQuerySystemInformation(t *testing.T) {
	buf, ret := NtQuerySystemInformation(0x10, 0, true)
	fmt.Printf("NtQuerySystemInformation returned: 0x%08X (%s), Length: %d\n", ret, exitcodes.FormatError(ret), len(buf))
	if ret != 0 {
		t.Errorf("NtQuerySystemInformation failed with code: 0x%08X (%s)", ret, exitcodes.FormatError(ret))
	}

	for index := 0; index < len(buf); index += int(unsafe.Sizeof(uintptr(0))) {
		// You may want to process the buffer here
	}
}

func TestNtQuerySystemInformationEx(t *testing.T) {
	// Call with correct parameters: class, processorGroup, initialSize, debug
	buf, ret := NtQuerySystemInformationEx(0x08, 0, 0, true)
	fmt.Printf("NtQuerySystemInformationEx returned: 0x%08X (%s), Length: %d\n", ret, exitcodes.FormatError(ret), len(buf))
	if ret != 0 {
		t.Errorf("NtQuerySystemInformationEx failed with code: 0x%08X (%s)", ret, exitcodes.FormatError(ret))
	}
}
