//go:build windows

package collector

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/AppeiYA/x-parity/internal/domain"
)

func collectPlatformTelemetry(ctx context.Context, rt *domain.Runtime) error {
	// 1. Interrogate Windows Registry for exact OS Name, Version, Build and UBR
	productName, displayVersion, buildNum, ubr := readWindowsRegistryVersion()

	if productName != "" {
		kernelStr := productName
		if displayVersion != "" {
			kernelStr += " " + displayVersion
		}
		if buildNum != "" {
			if ubr > 0 {
				kernelStr += fmt.Sprintf(" (Build %s.%d)", buildNum, ubr)
			} else {
				kernelStr += fmt.Sprintf(" (Build %s)", buildNum)
			}
		}
		rt.Kernel = kernelStr
	} else {
		// Fallback to "cmd /c ver"
		if out, err := exec.CommandContext(ctx, "cmd", "/c", "ver").Output(); err == nil {
			rt.Kernel = strings.TrimSpace(string(out))
		} else {
			rt.Kernel = "Windows NT"
		}
	}

	// 2. Check for Windows Container
	if isWindowsContainer() {
		rt.Virtualization = "windows-container"
	}

	return nil
}

func readWindowsRegistryVersion() (productName, displayVersion, buildNum string, ubr uint32) {
	subkey, err := syscall.UTF16PtrFromString(`SOFTWARE\Microsoft\Windows NT\CurrentVersion`)
	if err != nil {
		return
	}

	var hKey syscall.Handle
	if err := syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, subkey, 0, syscall.KEY_READ, &hKey); err != nil {
		return
	}
	defer syscall.RegCloseKey(hKey)

	productName = queryRegString(hKey, "ProductName")
	displayVersion = queryRegString(hKey, "DisplayVersion")
	if displayVersion == "" {
		displayVersion = queryRegString(hKey, "ReleaseId")
	}
	buildNum = queryRegString(hKey, "CurrentBuildNumber")
	if buildNum == "" {
		buildNum = queryRegString(hKey, "CurrentBuild")
	}
	ubr = queryRegDword(hKey, "UBR")

	return
}

func queryRegString(hKey syscall.Handle, name string) string {
	namePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return ""
	}

	var valType uint32
	var bufLen uint32
	// First query size
	if err := syscall.RegQueryValueEx(hKey, namePtr, nil, &valType, nil, &bufLen); err != nil || bufLen == 0 {
		return ""
	}

	buf := make([]uint16, bufLen/2+1)
	if err := syscall.RegQueryValueEx(hKey, namePtr, nil, &valType, (*byte)(unsafe.Pointer(&buf[0])), &bufLen); err != nil {
		return ""
	}

	return syscall.UTF16ToString(buf)
}

func queryRegDword(hKey syscall.Handle, name string) uint32 {
	namePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0
	}

	var valType uint32
	var val uint32
	bufLen := uint32(4)
	if err := syscall.RegQueryValueEx(hKey, namePtr, nil, &valType, (*byte)(unsafe.Pointer(&val)), &bufLen); err != nil {
		return 0
	}
	return val
}

func isWindowsContainer() bool {
	if os.Getenv("CONTAINER_SANDBOX_MOUNT_POINT") != "" {
		return true
	}
	if _, err := os.Stat(`C:\ContainerMappedDirectories`); err == nil {
		return true
	}
	return false
}

