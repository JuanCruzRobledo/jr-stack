---
description: Implement tasks from an OPSX change
agent: sdd-orchestrator
---

Load the `openspec-apply-change` skill and follow it exactly.

CONTEXT:
- Working directory: !`echo -n "$(pwd)"`
- Current project: !`echo -n "$(basename $(pwd))"`
- Change name: $ARGUMENTS
