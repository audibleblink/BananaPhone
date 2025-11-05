package bananaphone

import (
	"fmt"
	"sync"
	"unsafe"
)

// Global BananaPhone instance for standalone functions
var (
	globalBP     *BananaPhone
	globalBPErr  error
	globalBPOnce sync.Once
)

// getGlobalBananaPhone returns the global BananaPhone instance, initializing it if necessary.
// Uses sync.Once to ensure thread-safe lazy initialization.
func getGlobalBananaPhone() (*BananaPhone, error) {
	globalBPOnce.Do(func() {
		globalBP, globalBPErr = NewBananaPhone(AutoBananaPhoneMode)
	})
	return globalBP, globalBPErr
}

//Syscall calls the system function specified by callid with n arguments. Works much the same as syscall.Syscall - return value is the call error code and optional error text. All args are uintptrs to make it easy.
func Syscall(callid uint16, argh ...uintptr) (errcode uint32, err error) {
	errcode = bpSyscall(callid, argh...)

	if errcode != 0 {
		err = fmt.Errorf("non-zero return from syscall")
	}
	return errcode, err
}

//Syscall calls the system function specified by callid with n arguments. Works much the same as syscall.Syscall - return value is the call error code and optional error text. All args are uintptrs to make it easy.
func bpSyscall(callid uint16, argh ...uintptr) (errcode uint32)

//GetPEB returns the in-memory address of the start of PEB while making no api calls
func GetPEB() uintptr

//GetNtdllStart returns the start address of ntdll in memory
func GetNtdllStart() (start uintptr, size uintptr)

//getModuleLoadedOrder returns the start address of module located at i in the load order. This might be useful if there is a function you need that isn't in ntdll, or if some rude individual has loaded themselves before ntdll.
func getModuleLoadedOrder(i int) (start uintptr, size uintptr, modulepath *stupidstring)

//GetModuleLoadedOrderPtr returns a pointer to the ldr data table entry in full, incase there is something interesting in there you want to see.
func GetModuleLoadedOrderPtr(i int) *LdrDataTableEntry

//GetModuleLoadedOrder returns the start address of module located at i in the load order. This might be useful if there is a function you need that isn't in ntdll, or if some rude individual has loaded themselves before ntdll.
func GetModuleLoadedOrder(i int) (start uintptr, size uintptr, modulepath string) {
	var badstring *stupidstring
	start, size, badstring = getModuleLoadedOrder(i)
	modulepath = badstring.String()
	return
}

//Image contains info about a loaded image. Literally just a Base Addr and a Size - it should allow someone with a handy PE parser to pull the image out of memory...
type Image struct {
	BaseAddr uint64
	Size     uint64
}

//InMemLoads returns a map of loaded dll paths to current process offsets (aka images) in the current process. No syscalls are made.
func InMemLoads() (map[string]Image, error) {
	ret := make(map[string]Image)
	s, si, p := GetModuleLoadedOrder(0)
	start := p
	i := 1
	ret[p] = Image{uint64(s), uint64(si)}
	for {
		s, si, p = GetModuleLoadedOrder(i)
		if p != "" {
			ret[p] = Image{uint64(s), uint64(si)}
		}
		if p == start {
			break
		}
		i++
	}

	return ret, nil
}

//GetSysIDFromMemory takes the exported syscall name or ordinal and gets the ID it refers to (try not to supply both, it might not work how you expect). This function will not use a clean version of the dll, if AV has hooked the in-memory ntdll module, the results of this call may be bad.
func GetSysIDFromMemory(funcname string) (uint16, error) {
	return getSysIDFromMemory(funcname, 0, false)
}

//GetSysIDFromDiskOrd takes the exported ordinal and gets the ID it refers to. This function will access the ntdll file _on disk_, and relevant events/logs will be generated for those actions.
func GetSysIDFromDiskOrd(ordinal uint32) (uint16, error) {
	return getSysIDFromDisk("", ordinal, true)
}

//GetSysIDFromDisk takes the exported syscall name and gets the ID it refers to. This function will access the ntdll file _on disk_, and relevant events/logs will be generated for those actions.
func GetSysIDFromDisk(funcname string) (uint16, error) {
	return getSysIDFromDisk(funcname, 0, false)
}

//WriteMemory writes the provided memory to the specified memory address. Does **not** check permissions, may cause panic if memory is not writable etc.
func WriteMemory(inbuf []byte, destination uintptr) {
	for index := uint32(0); index < uint32(len(inbuf)); index++ {
		writePtr := unsafe.Pointer(destination + uintptr(index))
		v := (*byte)(writePtr)
		*v = inbuf[index]
	}
}

// Protect changes the memory protection flags for a region of memory using NtProtectVirtualMemory.
// This function bypasses standard Windows API calls to avoid detection by AV/EDR systems.
//
// Parameters:
//   - hProcess: Handle to the process whose memory protection is to be changed
//   - baseAddr: Base address of the memory region to protect
//   - size: Size of the memory region in bytes
//   - newProtect: New protection flags (e.g., PAGE_READWRITE, PAGE_EXECUTE_READ)
//   - oldProtect: Pointer to receive the previous protection flags
//
// Returns:
//   - error: nil on success, error with details on failure
//
// Example:
//   var oldProtect uint32
//   err := bananaphone.Protect(hProcess, baseAddr, size, windows.PAGE_EXECUTE_READ, &oldProtect)
func Protect(hProcess uintptr, baseAddr uintptr, size uintptr, newProtect uint32, oldProtect *uint32) error {
	// Get or initialize the global BananaPhone instance
	bp, err := getGlobalBananaPhone()
	if err != nil {
		return fmt.Errorf("failed to initialize BananaPhone: %w", err)
	}

	// Resolve NtProtectVirtualMemory syscall ID
	protectSysID, err := bp.GetSysID("NtProtectVirtualMemory")
	if err != nil {
		return fmt.Errorf("failed to resolve NtProtectVirtualMemory: %w", err)
	}

	// Call NtProtectVirtualMemory
	// NTSTATUS NtProtectVirtualMemory(
	//   HANDLE ProcessHandle,
	//   PVOID *BaseAddress,
	//   PSIZE_T RegionSize,
	//   ULONG NewProtect,
	//   PULONG OldProtect
	// )
	errcode, err := Syscall(
		protectSysID,
		hProcess,
		uintptr(unsafe.Pointer(&baseAddr)),
		uintptr(unsafe.Pointer(&size)),
		uintptr(newProtect),
		uintptr(unsafe.Pointer(oldProtect)),
	)

	if errcode != 0 {
		return fmt.Errorf("NtProtectVirtualMemory failed with status: 0x%x", errcode)
	}

	return nil
}
