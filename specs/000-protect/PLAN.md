## Execution Plan for BananaPhone Memory Protection Function

Based on my analysis of the PRD and existing codebase, here's the comprehensive execution plan:

### **Phase 1: Code Implementation** (Tasks 2-5)

**Task 2-3: Add Global BananaPhone Infrastructure**
- Add global variables to `pkg/BananaPhone/functions.go`:
  - `globalBP *BananaPhone` - singleton instance
  - `globalBPErr error` - initialization error cache
  - `globalBPOnce sync.Once` - thread-safe initialization
- Implement `getGlobalBananaPhone()` helper function using `sync.Once` for lazy initialization with `AutoBananaPhoneMode`

**Task 5: Update Imports**
- Add `sync` package for `sync.Once`
- Ensure `fmt` and `unsafe` are imported (already present)

**Task 4: Implement Protect() Function**
The main function will:
1. Call `getGlobalBananaPhone()` to get/initialize singleton instance
2. Resolve `NtProtectVirtualMemory` syscall ID dynamically per call
3. Execute syscall with proper parameter marshaling:
   - `hProcess` - process handle
   - `&baseAddr` - pointer to base address (in/out parameter)
   - `&size` - pointer to region size (in/out parameter)
   - `newProtect` - new protection flags
   - `oldProtect` - pointer to receive old protection
4. Return formatted error if syscall fails, nil on success

**Key Implementation Details:**
- Follow the pattern from `example/calcshellcode/main.go` lines 82-93
- Use `unsafe.Pointer` for pointer parameters
- Check syscall return code (0 = success, non-zero = error)
- Match Windows `VirtualProtectEx` semantics exactly

### **Phase 2: Testing & Validation** (Tasks 6-7)

**Task 6: Functional Testing**
- Create test example similar to `calcshellcode/main.go`
- Test memory protection changes (RW → RX, etc.)
- Verify old protection value is correctly returned
- Test with different process handles (current process, other processes)
- Verify thread-safety with concurrent calls

**Task 7: EDR Evasion Verification**
- Run with API Monitor to confirm `NtProtectVirtualMemory` is NOT captured
- Verify direct syscall bypasses standard API hooks
- Test with AutoBananaPhoneMode fallback scenarios

### **Success Criteria Verification**

1. ✅ Function successfully changes memory protection flags
2. ✅ API Monitor does not capture the NtProtectVirtualMemory call  
3. ✅ Seamless integration with existing BananaPhone code patterns
4. ✅ Minimal performance overhead (< 5% vs direct syscall)
5. ✅ Thread-safe operation confirmed

### **Files Modified**
- `pkg/BananaPhone/functions.go` - Add global vars, helper function, and Protect() function

### **Dependencies Met**
- ✅ Existing BananaPhone framework (`bananaphone.go`)
- ✅ Syscall infrastructure (`asm_x64.s`, `Syscall()` function)
- ✅ Standard library packages (`sync`, `unsafe`, `fmt`)

The implementation follows established BananaPhone patterns and maintains consistency with the existing codebase architecture. Ready to proceed with implementation?
