---
name: openspec-tasks
description: >
  Create an implementation task checklist for a change using the openspec CLI.
  Trigger: When the orchestrator or user needs tasks created or updated for a change.
license: MIT
compatibility: Requires openspec CLI.
metadata:
  author: gentleman-programming
  version: "2.0"
---

## Purpose

Break down a change into concrete, actionable implementation tasks. Tasks are what the apply phase works through — each task should be completable in one focused session.

## Steps

1. **Get instructions from the CLI**

   ```bash
   openspec instructions tasks --change "<name>" --json
   ```

   Parse `outputPath`, `template`, `instruction`, `context`, `rules`, and `dependencies`.

2. **Read dependencies**

   Read any files listed in `dependencies` (typically the proposal, specs, and design).

3. **Write the tasks**

   Use `template` as the structure. Follow `instruction` for content guidance.
   Apply `context` and `rules` as constraints — do NOT copy them into the file.

4. **Verify**

   ```bash
   openspec status --change "<name>" --json
   ```

   Confirm the tasks artifact shows `status: "done"`.

## Output Format

```markdown
# Tasks: <Change Name>

## Phase 1: <Phase Name>

- [ ] 1.1 <Specific, actionable task>
- [ ] 1.2 <Specific, actionable task>
- [ ] 1.3 <Specific, actionable task>

## Phase 2: <Phase Name>

- [ ] 2.1 <Specific, actionable task>
- [ ] 2.2 <Specific, actionable task>
```

## Rules

- Each task must be concrete and completable (not "implement auth" — too broad)
- Group by logical phase (infrastructure, core logic, UI, tests)
- Use hierarchical numbering (1.1, 1.2, 2.1, 2.2)
- Tasks are checked off during apply: `- [ ]` becomes `- [x]`
- Apply any `rules.tasks` from `openspec/config.yaml`
