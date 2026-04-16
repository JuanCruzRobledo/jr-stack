---
name: openspec-init
description: >
  Initialize OPSX in a project. Runs `openspec init` to scaffold the openspec/ directory and config.
  Trigger: When user wants to initialize OPSX in a project, or says "opsx init", "iniciar opsx", "openspec init".
license: MIT
compatibility: Requires openspec CLI.
metadata:
  author: gentleman-programming
  version: "2.0"
---

## Purpose

Initialize OPSX in the current project by running the openspec CLI scaffolding.

## Steps

### Step 1: Check if already initialized

```bash
openspec status --json
```

If this succeeds without error, OPSX is already initialized. Show current status and stop.

If it fails (no `openspec/` directory), proceed to initialization.

### Step 2: Run initialization

```bash
openspec init
```

This creates:
```
openspec/
├── config.yaml     <- Project-specific config (schema, context, rules)
├── specs/          <- Source of truth for specifications
└── changes/        <- Active changes
    └── archive/    <- Completed changes
```

### Step 3: Customize config

Open `openspec/config.yaml` and help the user fill in the `context` field with their project's tech stack, conventions, and domain knowledge. This context is injected into every artifact the AI creates.

Example:
```yaml
schema: spec-driven

context: |
  Tech stack: TypeScript, React, Node.js
  We use conventional commits
  Domain: e-commerce platform
  Testing: Vitest + Testing Library
```

### Step 4: Confirm

Run:
```bash
openspec list --json
```

Show the user the initialized structure and explain next steps:
- "Run `/opsx:propose` to start your first change"
- "Or `/opsx:explore` to think through what to build first"

## Rules

- NEVER create the openspec/ structure manually — always use the CLI
- If already initialized, do NOT re-initialize — just show current status
- Help the user fill in config.yaml context — this is the most important customization
