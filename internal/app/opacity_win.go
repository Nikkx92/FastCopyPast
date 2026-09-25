package app

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	procFindWindowW          = user32.NewProc("FindWindowW")
	procGetWindowLongPtrW    = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtrW    = user32.NewProc("SetWindowLongPtrW")
	procSetLayeredWindowAttr = user32.NewProc("SetLayeredWindowAttributes")
)

var globalHWND uintptr

const (
	WsExLayered = 0x00080000
	LwaAlpha    = 0x00000002
)

func SetOpacity(appTitle string) {
	gwlExStyle := -20
	titlePtr, _ := syscall.UTF16PtrFromString(appTitle)
	var hwnd uintptr

	for i := 0; i < 50; i++ {
		hwnd, _, _ = procFindWindowW.Call(0, uintptr(unsafe.Pointer(titlePtr)))
		if hwnd != 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if hwnd == 0 {
		println("Ошибка: Окно не найдено за отведенное время")
		return
	}

	globalHWND = hwnd

	currentStyle, _, _ := procGetWindowLongPtrW.Call(hwnd, uintptr(gwlExStyle))
	if currentStyle == 0 {
		procGetWindowLongW := user32.NewProc("GetWindowLongW")
		currentStyle, _, _ = procGetWindowLongW.Call(hwnd, uintptr(gwlExStyle))
	}

	newStyle := currentStyle | WsExLayered

	_, _, err := procSetWindowLongPtrW.Call(hwnd, uintptr(gwlExStyle), newStyle)
	if err != nil && err.(syscall.Errno) != 0 {
		procSetWindowLongW := user32.NewProc("SetWindowLongW")
		procSetWindowLongW.Call(hwnd, uintptr(gwlExStyle), newStyle)
	}

	UpdateOpacity(255)
}

func UpdateOpacity(alpha uint8) {
	if globalHWND == 0 {
		return
	}

	procSetLayeredWindowAttr.Call(globalHWND, 0, uintptr(alpha), LwaAlpha)
}
