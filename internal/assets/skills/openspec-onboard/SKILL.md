---
name: openspec-onboard
description: >
  Guided end-to-end walkthrough of the OPSX workflow using the real codebase.
  Trigger: When the user wants to learn OPSX or do a guided first change.
license: MIT
compatibility: Requires openspec CLI.
metadata:
  author: gentleman-programming
  version: "2.0"
---

## Purpose

Guide the user through a complete OPSX cycle using their actual codebase. This is a real change with real artifacts — the goal is to teach by doing.

## Flow

```
1. Check: is openspec installed and initialized?
2. Explore: think about what to build
3. Propose: create a real change with artifacts
4. Apply: implement the tasks
5. Archive: close the cycle
```

## Steps

### Step 1: Check Prerequisites

```bash
openspec --version
openspec list --json
```

If `openspec` is not installed:
- Direct the user to install it: `npm install -g openspec` or the relevant method
- Stop until it's available

If not initialized:
- Run through the `openspec-init` skill first

### Step 2: Pick a Real Change

Ask the user:
> "What's something small but real you want to build or fix in this project? Could be a bug fix, a small feature, or a refactor. Real is better than made-up."

Help them pick something that:
- Is concrete and completable (not "redesign the auth system")
- Has a clear definition of done
- Won't take more than a few hours

### Step 3: Walk Through the Full Cycle

Now guide them through each OPSX action:

**Explore** (optional but recommended):
- Load `openspec-explore`
- Think through the change together

**Propose**:
- Load `openspec-propose`
- Create the change with all artifacts

**Apply**:
- Load `openspec-apply-change`
- Implement at least one task together
- Show them how task checkboxes work

**Archive** (when done):
- Load `openspec-archive-change`
- Show how the change gets closed out

### Step 4: Debrief

After the cycle, explain:
- What each artifact does (proposal = what/why, design = how, tasks = checklist)
- How to run any step again: `/opsx:apply` to continue implementing
- How to update artifacts mid-flight: just edit the file and re-run apply
- The key OPSX principle: **no phase locks** — you can update anything at any time

## Rules

- Use the user's real codebase, not toy examples
- Don't rush — this is learning time
- After each step, check if they have questions before continuing
- Emphasize the fluidity: OPSX doesn't have phase gates
