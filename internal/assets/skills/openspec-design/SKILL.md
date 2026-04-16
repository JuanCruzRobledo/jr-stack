---
name: openspec-design
description: >
  Create a technical design document for a change using the openspec CLI.
  Trigger: When the orchestrator or user needs a design document created or updated for a change.
license: MIT
compatibility: Requires openspec CLI.
metadata:
  author: gentleman-programming
  version: "2.0"
---

## Purpose

Write the technical design for a change. Design defines HOW the system will be built — architecture decisions, data models, API contracts, component structure.

## Steps

1. **Get instructions from the CLI**

   ```bash
   openspec instructions design --change "<name>" --json
   ```

   Parse `outputPath`, `template`, `instruction`, `context`, `rules`, and `dependencies`.

2. **Read dependencies**

   Read any files listed in `dependencies` (typically the proposal and specs).

3. **Read existing code**

   Before designing, understand the current codebase:
   - Existing architecture patterns
   - Current data models
   - Integration points
   - Tech stack conventions

4. **Write the design**

   Use `template` as the structure. Follow `instruction` for content guidance.
   Apply `context` and `rules` as constraints — do NOT copy them into the file.

5. **Verify**

   ```bash
   openspec status --change "<name>" --json
   ```

   Confirm the design artifact shows `status: "done"`.

## Output Format

```markdown
# Design: <Change Name>

## Architecture Overview
<High-level description of the approach>

## Components

### <Component Name>
- **Responsibility**: <what it does>
- **Location**: `path/to/component`
- **Interface**: <key methods/props/API>

## Data Model
<Schemas, types, database changes>

## API Changes
<New endpoints, modified contracts>

## Implementation Notes
<Non-obvious decisions and their rationale>

## Risks & Mitigations
| Risk | Mitigation |
|------|------------|
| <risk> | <how we handle it> |
```

## Rules

- Design defines HOW, not WHAT — specs define the requirements
- Document WHY decisions were made, not just what they are
- Include sequence diagrams for complex async flows (ASCII art is fine)
- Keep it CONCISE — it's a reference, not a novel
- Apply any `rules.design` from `openspec/config.yaml`
