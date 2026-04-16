---
name: openspec-verify
description: >
  Validate that implementation matches specs, design, and tasks.
  Trigger: When the user wants to verify a completed or partially completed change.
license: MIT
compatibility: Requires openspec CLI.
metadata:
  author: gentleman-programming
  version: "2.0"
---

## Purpose

Validate that the implementation matches what was planned. Check completeness (all tasks done?), correctness (code matches specs?), and coherence (design decisions followed?).

## Steps

1. **Select the change**

   If not provided, run:
   ```bash
   openspec list --json
   ```
   And ask the user to select.

2. **Read context**

   ```bash
   openspec status --change "<name>" --json
   ```

   Then read:
   - `openspec/changes/<name>/specs/` — acceptance criteria
   - `openspec/changes/<name>/tasks.md` — what was planned
   - `openspec/changes/<name>/design.md` — how it should be built
   - The actual implementation files

3. **Run tests if available**

   Check if the project has a test runner and run it:
   ```bash
   # Examples:
   npm test
   go test ./...
   pytest
   ```

4. **Check completeness**

   Count `- [x]` vs `- [ ]` in tasks.md.
   Report: "N/M tasks complete"

5. **Check correctness**

   For each spec requirement:
   - Find the corresponding implementation
   - Verify it satisfies the scenario/acceptance criteria
   - Flag: PASS / FAIL / PARTIAL

6. **Check coherence**

   For key design decisions:
   - Verify the implementation follows the documented approach
   - Flag any deviations

7. **Build the compliance matrix**

   ```
   | Requirement | Status | Notes |
   |-------------|--------|-------|
   | REQ-001     | PASS   | -     |
   | REQ-002     | FAIL   | Missing error handling |
   ```

8. **Output report**

   Save to `openspec/changes/<name>/verify-report.md`.

## Output Format

```markdown
## Verification Report: <change-name>

**Date**: YYYY-MM-DD
**Tasks**: N/M complete

### Test Results
<test output or "No test runner detected">

### Spec Compliance
| Requirement | Status | Notes |
|-------------|--------|-------|
| REQ-001 | PASS | - |
| REQ-002 | FAIL | Missing validation |

### Design Coherence
- <decision>: FOLLOWED / DEVIATED (reason)

### Summary
- CRITICAL: <blockers that must be fixed before archive>
- WARNING: <issues worth noting but not blocking>
- SUGGESTION: <improvements for the future>

**Verdict**: READY FOR ARCHIVE / NEEDS FIXES
```

## Rules

- CRITICAL issues block archive — they must be fixed
- WARNING and SUGGESTION don't block archive
- Always run tests if a test runner exists
- If no specs exist, check against the proposal's success criteria instead
- Apply any `rules.verify` from `openspec/config.yaml`
