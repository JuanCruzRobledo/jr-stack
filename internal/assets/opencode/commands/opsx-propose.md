---
description: Propose a new OPSX change — create it and generate all artifacts in one step
agent: sdd-orchestrator
---

Load the `openspec-propose` skill and follow it exactly.

CONTEXT:
- Working directory: !`echo -n "$(pwd)"`
- Current project: !`echo -n "$(basename $(pwd))"`
- Change name or description: $ARGUMENTS
