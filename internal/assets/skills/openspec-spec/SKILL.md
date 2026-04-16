---
name: openspec-spec
description: >
  Write or update specifications for a change using the openspec CLI.
  Trigger: When the orchestrator or user needs specs created or updated for a change.
license: MIT
compatibility: Requires openspec CLI.
metadata:
  author: gentleman-programming
  version: "2.0"
---

## Purpose

Write delta specs for a change. Specs define WHAT the system must do — they are your acceptance criteria for implementation.

## Steps

1. **Get instructions from the CLI**

   ```bash
   openspec instructions specs --change "<name>" --json
   ```

   Parse `outputPath`, `template`, `instruction`, `context`, `rules`, and `dependencies`.

2. **Read dependencies**

   Read any files listed in `dependencies` (typically the proposal).

3. **Write the spec**

   Use `template` as the structure. Follow `instruction` for content guidance.
   Apply `context` and `rules` as constraints — do NOT copy them into the file.

   Specs should use:
   - Given/When/Then format for scenarios
   - RFC 2119 keywords (MUST, SHALL, SHOULD, MAY) for requirements
   - Concrete, testable acceptance criteria

4. **Verify**

   ```bash
   openspec status --change "<name>" --json
   ```

   Confirm the spec artifact shows `status: "done"`.

## Output Format

```markdown
# Spec: <capability-name>

## Overview
<Brief description of what this capability covers>

## Requirements

### REQ-001: <Requirement Name>
The system MUST <specific behavior>.

**Scenarios:**

**Scenario: <scenario name>**
- Given: <precondition>
- When: <action>
- Then: <expected outcome>
```

## Rules

- Specs define WHAT, not HOW — leave implementation details to design
- Every requirement must have at least one testable scenario
- Use concrete examples, not abstract descriptions
- Apply any `rules.specs` from `openspec/config.yaml`
