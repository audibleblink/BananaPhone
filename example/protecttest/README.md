# Protect() Function Test Example

This example demonstrates the new `Protect()` function in BananaPhone, which provides a "bananified" version of `VirtualProtectEx` that bypasses standard Windows API calls to avoid detection by AV/EDR systems.

## What it does

The test program:
1. Allocates memory with `PAGE_READWRITE` permissions using `NtAllocateVirtualMemory`
2. Writes test data (NOP instructions) to the allocated memory
3. Changes protection to `PAGE_EXECUTE_READ` using the new `Protect()` function
4. Changes protection back to `PAGE_READWRITE`
5. Changes protection to `PAGE_EXECUTE_READWRITE`
6. Verifies that old protection values are correctly returned

## Building

```bash
go build .
```

## Running

```bash
./protecttest.exe
```

Expected output:
```
Testing BananaPhone Protect() function...
✓ Allocated memory at: 0x... (size: 5 bytes)
✓ Wrote 5 bytes to memory
✓ Changed protection to PAGE_EXECUTE_READ (old: 0x4)
✓ Changed protection back to PAGE_READWRITE (old: 0x20)
✓ Changed protection to PAGE_EXECUTE_READWRITE (old: 0x4)

✅ All tests passed! Protect() function is working correctly.

Note: Run with API Monitor to verify NtProtectVirtualMemory is NOT captured.
```

## EDR Evasion Verification

To verify that the function successfully bypasses API monitoring:

1. Install [API Monitor](http://www.rohitab.com/apimonitor)
2. Configure it to monitor `VirtualProtectEx` and `NtProtectVirtualMemory`
3. Run this test program under API Monitor
4. Verify that `NtProtectVirtualMemory` calls are **NOT** captured (because we're using direct syscalls)

## Usage in Your Code

Before (using standard Windows API):
```go
err := windows.VirtualProtectEx(windows.Handle(hProcess), baseAddr, 1, windows.PAGE_READWRITE, &oldProtect)
if err != nil {
    return fmt.Errorf("VirtualProtectEx failed: %w", err)
}
```

After (using BananaPhone):
```go
err := bananaphone.Protect(hProcess, baseAddr, 1, windows.PAGE_READWRITE, &oldProtect)
if err != nil {
    return fmt.Errorf("Protect failed: %w", err)
}
```

## Features

- **Thread-safe**: Uses `sync.Once` for lazy initialization of global BananaPhone instance
- **Auto mode**: Uses `AutoBananaPhoneMode` for maximum compatibility (falls back to disk/HalosGate if memory is hooked)
- **Error handling**: Provides detailed error messages for debugging
- **API compatible**: Matches Windows `VirtualProtectEx` behavior exactly
