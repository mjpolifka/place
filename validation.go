//go:build windows

package main

import (
	"fmt"
	"strings"
	"syscall"
	"unicode"
	"unsafe"
)

type displayDimensions struct {
	Height        int
	Width         int
	DisplayNumber int
	IsMain        bool
	Left          int
	Top           int
	Right         int
	Bottom        int
}

type rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type monitorInfo struct {
	CbSize    uint32
	RcMonitor rect
	RcWork    rect
	DwFlags   uint32
}

func normalizeProcessName(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("window name cannot be empty")
	}

	for _, ch := range trimmed {
		if ch == '/' || ch == '\\' {
			return "", fmt.Errorf("window name cannot contain path separators")
		}
		if unicode.IsControl(ch) {
			return "", fmt.Errorf("window name cannot contain control characters")
		}
	}

	normalized := strings.ToLower(trimmed)
	normalized = strings.TrimSuffix(normalized, ".exe")
	normalized = normalized + ".exe"

	return normalized, nil
}

func validateIntOverflow(i int) error {
	minWindowInt32 := -2147483648
	maxWindowInt32 := 2147483647
	if i < minWindowInt32 || i > maxWindowInt32 {
		return fmt.Errorf("Int must be in signed 32-bit range [%d, %d]: %d", minWindowInt32, maxWindowInt32, i)
	}
	return nil
}

func validateWindowCoord(value int) error {
	maxWindowCoord := 10000 // Max 2147483647
	err := validateIntOverflow(value)
	if err != nil {
		return err
	}
	if value < -maxWindowCoord || value > maxWindowCoord {
		return fmt.Errorf("Window coord must be in range [0, %d]: %d", maxWindowCoord, value)
	}
	return nil
}

func validateWindowSize(value int) error {
	maxWindowSize := 10000
	err := validateIntOverflow(value)
	if err != nil {
		return err
	}
	if value < 0 || value > maxWindowSize {
		return fmt.Errorf("Window size must be in range [0, %d]: %d", maxWindowSize, value)
	}
	return nil
}

func getDisplayDimensions() ([]displayDimensions, error) {
	user32DLL := syscall.NewLazyDLL("user32.dll")
	enumDisplayMonitorsProc := user32DLL.NewProc("EnumDisplayMonitors")
	getMonitorInfoWProc := user32DLL.NewProc("GetMonitorInfoW")

	displays := []displayDimensions{}
	displayNumber := 1

	callback := syscall.NewCallback(func(hMonitor uintptr, hdcMonitor uintptr, lprcMonitor uintptr, dwData uintptr) uintptr {
		info := monitorInfo{CbSize: uint32(unsafe.Sizeof(monitorInfo{}))}
		ret, _, _ := getMonitorInfoWProc.Call(
			hMonitor,
			uintptr(unsafe.Pointer(&info)),
		)
		if ret == 0 {
			return 1 // continue
		}

		width := int(info.RcMonitor.Right - info.RcMonitor.Left)
		height := int(info.RcMonitor.Bottom - info.RcMonitor.Top)

		displays = append(displays, displayDimensions{
			Height:        height,
			Width:         width,
			DisplayNumber: displayNumber,
			IsMain:        info.DwFlags&0x1 != 0, // MONITORINFOF_PRIMARY
			Left:          int(info.RcMonitor.Left),
			Top:           int(info.RcMonitor.Top),
			Right:         int(info.RcMonitor.Right),
			Bottom:        int(info.RcMonitor.Bottom),
		})
		displayNumber++

		return 1 // continue
	})

	ret, _, err := enumDisplayMonitorsProc.Call(
		0,
		0,
		callback,
		0,
	)
	if ret == 0 {
		return nil, fmt.Errorf("failed to enumerate displays: %w", err)
	}
	if len(displays) == 0 {
		return nil, fmt.Errorf("no displays found")
	}

	return displays, nil
}

func validatePointWithinDisplays(x int, y int, displays []displayDimensions) error {
	for _, display := range displays {
		if x >= display.Left && x < display.Right && y >= display.Top && y < display.Bottom {
			return nil
		}
	}

	return fmt.Errorf("point (%d, %d) is outside all detected displays", x, y)
}
