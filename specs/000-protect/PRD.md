# Product Requirements Document: BananaPhone Memory Protection Function

## Overview
Add a `Protect()` function to BananaPhone that allows setting memory protection flags using direct syscalls to `NtProtectVirtualMemory`, bypassing standard Windows API calls to avoid detection by AV/EDR systems.

## Background
Currently, users need to use `windows.VirtualProtectEx()` directly when changing memory protection flags, which defeats the purpose of using BananaPhone for EDR evasion. This function will complete the memory management API by providing a "bananified" version of VirtualProtectEx.

## Requirements

### Functional Requirements

#### FR1: Function Signature
```go
func Protect(hProcess uintptr, baseAddr uintptr, size uintptr, newProtect uint32, oldProtect *uint32) error
```

- **hProcess**: Handle to the process whose memory protection is to be changed
- **baseAddr**: Base address of the memory region to protect
- **size**: Size of the memory region in bytes
- **newProtect**: New protection flags (e.g., PAGE_READWRITE, PAGE_EXECUTE_READ)
- **oldProtect**: Pointer to receive the previous protection flags
- **Returns**: Error if the operation fails, nil on success

#### FR2: Implementation Pattern
- Standalone package function in `functions.go`
- Uses global BananaPhone instance with lazy initialization
- Thread-safe initialization using `sync.Once`
- Calls `NtProtectVirtualMemory` syscall directly

#### FR3: Syscall Resolution
- Use `AutoBananaPhoneMode` for maximum compatibility
- Resolve `NtProtectVirtualMemory` sysid dynamically per call
- Handle syscall errors with descriptive error messages

#### FR4: API Compatibility
- Match Windows `VirtualProtectEx` behavior
- Accept same parameter types and semantics
- Return old protection value via pointer parameter

### Non-Functional Requirements

#### NFR1: Performance
- BananaPhone instance created only once (lazy initialization)
- Minimal overhead compared to direct syscall
- Sysid resolution occurs per call (trade-off for anti-hooking)

#### NFR2: Security
- Must bypass API monitoring tools (API Monitor, EDR hooks)
- Use direct syscall via BananaPhone framework
- Support all BananaPhone evasion modes (Memory, Disk, Auto, HalosGate)

#### NFR3: Maintainability
- Follow existing BananaPhone code patterns
- Clear error messages for debugging
- Consistent with `WriteMemory()` standalone function pattern

#### NFR4: Compatibility
- Windows-only (like rest of BananaPhone)
- Support all memory protection constants from Windows API
- Handle process handles consistently with Windows API

## Implementation Details

### Files to Modify
1. **pkg/BananaPhone/functions.go**
   - Add global BananaPhone instance variables
   - Add `getGlobalBananaPhone()` helper function
   - Add `Protect()` function implementation

### Code Structure

```go
// Global variables (add to functions.go)
var (
    globalBP     *BananaPhone
    globalBPErr  error
    globalBPOnce sync.Once
)

// Helper function (add to functions.go)
func getGlobalBananaPhone() (*BananaPhone, error) {
    globalBPOnce.Do(func() {
        globalBP, globalBPErr = NewBananaPhone(AutoBananaPhoneMode)
    })
    return globalBP, globalBPErr
}

// Main function (add to functions.go)
func Protect(hProcess uintptr, baseAddr uintptr, size uintptr, newProtect uint32, oldProtect *uint32) error {
    // Implementation as per Option 2
}
```

### Error Handling
- Return initialization errors if BananaPhone creation fails
- Return sysid resolution errors if NtProtectVirtualMemory lookup fails
- Return formatted error with error code if syscall returns non-zero

### Usage Example

Before (using standard Windows API):
```go
err := windows.VirtualProtectEx(windows.Handle(hProcess), baseAddr, 1, windows.PAGE_READWRITE, &oldProtect)
if err != nil {
    return log.Add("VirtualProtectEx[rw]").Wrap(err)
}

bananaphone.WriteMemory(data, baseAddr)

err = windows.VirtualProtectEx(windows.Handle(hProcess), baseAddr, 1, oldProtect, &old)
if err != nil {
    return log.Add("VirtualProtectEx[rx]").Wrap(err)
}
```

After (using BananaPhone):
```go
err := bananaphone.Protect(hProcess, baseAddr, 1, windows.PAGE_READWRITE, &oldProtect)
if err != nil {
    return log.Add("Protect[rw]").Wrap(err)
}

bananaphone.WriteMemory(data, baseAddr)

err = bananaphone.Protect(hProcess, baseAddr, 1, oldProtect, &old)
if err != nil {
    return log.Add("Protect[rx]").Wrap(err)
}
```

## Success Criteria
1. Function successfully changes memory protection flags
2. API Monitor does not capture the NtProtectVirtualMemory call
3. Function integrates seamlessly with existing BananaPhone code
4. Performance overhead is minimal (< 5% compared to direct syscall)
5. Thread-safe operation in concurrent scenarios

## Out of Scope
- Additional memory management functions (VirtualAlloc, VirtualFree wrappers)
- Custom memory protection constant definitions (use Windows API constants)
- Helper functions for common protection patterns
- Extended documentation beyond code comments

## Dependencies
- Existing BananaPhone framework
- Windows syscall infrastructure
- sync package for Once initialization
- unsafe package for pointer manipulation

## Timeline
Single implementation phase - add function to existing codebase without breaking changes.
