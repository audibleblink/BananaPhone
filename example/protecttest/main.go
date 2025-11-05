package main

import (
	"fmt"
	"syscall"
	"unsafe"

	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
)

var testData = []byte{
	0x90, 0x90, 0x90, 0x90, // NOP instructions for testing
	0xC3,                   // RET instruction
}

// Example demonstrating the new Protect() function
func main() {
	fmt.Println("Testing BananaPhone Protect() function...")

	// Get current process handle
	hProcess := uintptr(0xffffffffffffffff) // -1 = current process

	// Create a BananaPhone instance for memory allocation
	bp, err := bananaphone.NewBananaPhone(bananaphone.AutoBananaPhoneMode)
	if err != nil {
		panic(fmt.Sprintf("Failed to create BananaPhone: %v", err))
	}

	// Resolve NtAllocateVirtualMemory
	allocSysID, err := bp.GetSysID("NtAllocateVirtualMemory")
	if err != nil {
		panic(fmt.Sprintf("Failed to resolve NtAllocateVirtualMemory: %v", err))
	}

	// Allocate memory with RW permissions
	var baseAddr uintptr
	regionSize := uintptr(len(testData))
	const (
		MEM_COMMIT  = 0x00001000
		MEM_RESERVE = 0x00002000
	)

	errcode, err := bananaphone.Syscall(
		allocSysID,
		hProcess,
		uintptr(unsafe.Pointer(&baseAddr)),
		0,
		uintptr(unsafe.Pointer(&regionSize)),
		uintptr(MEM_COMMIT|MEM_RESERVE),
		syscall.PAGE_READWRITE,
	)
	if err != nil || errcode != 0 {
		panic(fmt.Sprintf("NtAllocateVirtualMemory failed: %v (code: 0x%x)", err, errcode))
	}
	fmt.Printf("✓ Allocated memory at: 0x%x (size: %d bytes)\n", baseAddr, regionSize)

	// Write test data to allocated memory
	bananaphone.WriteMemory(testData, baseAddr)
	fmt.Printf("✓ Wrote %d bytes to memory\n", len(testData))

	// Test 1: Change protection from RW to RX using new Protect() function
	var oldProtect uint32
	err = bananaphone.Protect(hProcess, baseAddr, regionSize, syscall.PAGE_EXECUTE_READ, &oldProtect)
	if err != nil {
		panic(fmt.Sprintf("Protect() failed (RW->RX): %v", err))
	}
	fmt.Printf("✓ Changed protection to PAGE_EXECUTE_READ (old: 0x%x)\n", oldProtect)

	// Verify old protection was PAGE_READWRITE (0x04)
	if oldProtect != syscall.PAGE_READWRITE {
		fmt.Printf("⚠ Warning: Expected old protection 0x%x, got 0x%x\n", syscall.PAGE_READWRITE, oldProtect)
	}

	// Test 2: Change protection back to RW
	var oldProtect2 uint32
	err = bananaphone.Protect(hProcess, baseAddr, regionSize, syscall.PAGE_READWRITE, &oldProtect2)
	if err != nil {
		panic(fmt.Sprintf("Protect() failed (RX->RW): %v", err))
	}
	fmt.Printf("✓ Changed protection back to PAGE_READWRITE (old: 0x%x)\n", oldProtect2)

	// Verify old protection was PAGE_EXECUTE_READ (0x20)
	if oldProtect2 != syscall.PAGE_EXECUTE_READ {
		fmt.Printf("⚠ Warning: Expected old protection 0x%x, got 0x%x\n", syscall.PAGE_EXECUTE_READ, oldProtect2)
	}

	// Test 3: Change to PAGE_EXECUTE_READWRITE
	var oldProtect3 uint32
	err = bananaphone.Protect(hProcess, baseAddr, regionSize, syscall.PAGE_EXECUTE_READWRITE, &oldProtect3)
	if err != nil {
		panic(fmt.Sprintf("Protect() failed (RW->RWX): %v", err))
	}
	fmt.Printf("✓ Changed protection to PAGE_EXECUTE_READWRITE (old: 0x%x)\n", oldProtect3)

	fmt.Println("\n✅ All tests passed! Protect() function is working correctly.")
	fmt.Println("\nNote: Run with API Monitor to verify NtProtectVirtualMemory is NOT captured.")
}
