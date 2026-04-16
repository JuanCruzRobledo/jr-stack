---
description: Archive a completed OPSX change — sync delta specs and close the cycle
agent: sdd-orchestrator
---

Load the `openspec-archive-change` skill and follow it exactly.

CONTEXT:
- Working directory: !`echo -n "$(pwd)"`
- Current project: !`echo -n "$(basename $(pwd))"`
- Change name: $ARGUMENTS
