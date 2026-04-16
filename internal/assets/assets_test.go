package assets

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestAllEmbeddedAssetsAreReadable verifies that every expected embedded file
// can be loaded via Read() without error. This catches missing/misnamed files
// at test time rather than at runtime.
func TestAllEmbeddedAssetsAreReadable(t *testing.T) {
	expectedFiles := []string{
		// Claude agent files
		"claude/engram-protocol.md",
		"claude/persona-gentleman.md",
		"claude/sdd-orchestrator.md",

		// OpenCode agent files
		"opencode/persona-gentleman.md",
		"opencode/sdd-overlay-single.json",
		"opencode/sdd-overlay-multi.json",
		"opencode/commands/opsx-explore.md",
		"opencode/commands/opsx-propose.md",
		"opencode/commands/opsx-apply.md",
		"opencode/commands/opsx-archive.md",
		"opencode/plugins/background-agents.ts",

		// Gemini agent files
		"gemini/sdd-orchestrator.md",

		// Codex agent files
		"codex/sdd-orchestrator.md",

		// Cursor agent files
		"cursor/sdd-orchestrator.md",

		// Claude OPSX commands
		"claude/commands/opsx/explore.md",
		"claude/commands/opsx/propose.md",
		"claude/commands/opsx/apply.md",
		"claude/commands/opsx/archive.md",

		// OPSX skills
		"skills/openspec-init/SKILL.md",
		"skills/openspec-apply-change/SKILL.md",
		"skills/openspec-archive-change/SKILL.md",
		"skills/openspec-design/SKILL.md",
		"skills/openspec-explore/SKILL.md",
		"skills/openspec-propose/SKILL.md",
		"skills/openspec-spec/SKILL.md",
		"skills/openspec-tasks/SKILL.md",
		"skills/openspec-verify/SKILL.md",
		"skills/skill-registry/SKILL.md",
		"skills/_shared/persistence-contract.md",
		"skills/_shared/engram-convention.md",
		"skills/_shared/openspec-convention.md",
		"skills/_shared/sdd-phase-common.md",

		// Foundation skills
		"skills/go-testing/SKILL.md",
		"skills/skill-creator/SKILL.md",
	}

	for _, path := range expectedFiles {
		t.Run(path, func(t *testing.T) {
			content, err := Read(path)
			if err != nil {
				t.Fatalf("Read(%q) error = %v", path, err)
			}

			if len(strings.TrimSpace(content)) == 0 {
				t.Fatalf("Read(%q) returned empty content", path)
			}

			// Real content should be substantial, not a one-line stub.
			if len(content) < 50 {
				t.Fatalf("Read(%q) content is suspiciously short (%d bytes) — possible stub", path, len(content))
			}
		})
	}
}

func TestOpenCodeEmbeddedAssetLayout(t *testing.T) {
	entries, err := FS.ReadDir("opencode")
	if err != nil {
		t.Fatalf("ReadDir(opencode) error = %v", err)
	}

	seen := map[string]bool{}
	for _, entry := range entries {
		seen[entry.Name()] = true
	}

	for _, name := range []string{"commands", "plugins", "persona-gentleman.md", "sdd-overlay-single.json", "sdd-overlay-multi.json"} {
		if !seen[name] {
			t.Fatalf("opencode embedded assets missing %q", name)
		}
	}

	commandEntries, err := FS.ReadDir("opencode/commands")
	if err != nil {
		t.Fatalf("ReadDir(opencode/commands) error = %v", err)
	}
	if len(commandEntries) != 4 {
		t.Fatalf("opencode commands count = %d, want 4", len(commandEntries))
	}

	pluginEntries, err := FS.ReadDir("opencode/plugins")
	if err != nil {
		t.Fatalf("ReadDir(opencode/plugins) error = %v", err)
	}
	if len(pluginEntries) != 1 {
		t.Fatalf("opencode plugins count = %d, want 1", len(pluginEntries))
	}
	if pluginEntries[0].Name() != "background-agents.ts" {
		t.Fatalf("plugin entry = %q, want background-agents.ts", pluginEntries[0].Name())
	}
}

// TestMustReadPanicsOnMissingFile verifies that MustRead panics for a
// nonexistent file, confirming the safety mechanism works.
func TestMustReadPanicsOnMissingFile(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("MustRead() did not panic for missing file")
		}
	}()

	MustRead("nonexistent/file.md")
}

// TestEmbeddedAssetCount verifies we have the expected number of embedded files.
// This catches accidental deletions of asset files.
func TestEmbeddedAssetCount(t *testing.T) {
	// Count skill files.
	entries, err := FS.ReadDir("skills")
	if err != nil {
		t.Fatalf("ReadDir(skills) error = %v", err)
	}

	skillDirs := 0
	for _, entry := range entries {
		if entry.IsDir() {
			skillDirs++
		}
	}

	// We expect 17 skill directories (10 SDD + judgment-day + 5 foundation + _shared).
	if skillDirs != 17 {
		t.Fatalf("expected 17 skill directories, got %d", skillDirs)
	}

	// Verify each skill directory has a SKILL.md.
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == "_shared" {
			for _, sharedFile := range []string{"persistence-contract.md", "engram-convention.md", "openspec-convention.md", "sdd-phase-common.md", "skill-resolver.md"} {
				sharedPath := "skills/_shared/" + sharedFile
				if _, err := Read(sharedPath); err != nil {
					t.Fatalf("shared directory missing %q: %v", sharedFile, err)
				}
			}
			continue
		}
		skillPath := "skills/" + entry.Name() + "/SKILL.md"
		if _, err := Read(skillPath); err != nil {
			t.Fatalf("skill directory %q missing SKILL.md: %v", entry.Name(), err)
		}
	}
}

func TestSDDPhaseCommonEnforcesExecutorBoundary(t *testing.T) {
	content := MustRead("skills/_shared/sdd-phase-common.md")

	// Must enforce executor boundary — no delegation allowed.
	for _, want := range []string{
		"EXECUTOR, not an orchestrator",
		"Do NOT launch sub-agents",
		"do NOT call `delegate`/`task`",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("sdd-phase-common missing executor boundary rule %q", want)
		}
	}

	// Must instruct phase agents to search the skill registry themselves
	// when no explicit skill path was provided — this is skill LOADING, not delegation.
	if !strings.Contains(content, `mem_search(query: "skill-registry"`) {
		t.Fatal("sdd-phase-common must instruct phase agents to search skill-registry themselves for skill loading")
	}

	// Must NOT tell agents to launch sub-agents or delegate tasks.
	for _, forbidden := range []string{
		"launch a sub-agent",
		"delegate this to",
	} {
		if strings.Contains(content, forbidden) {
			t.Fatalf("sdd-phase-common should not contain delegation instruction %q", forbidden)
		}
	}
}

func TestOpenCodeSDDOverlayHasOrchestrator(t *testing.T) {
	// OPSX uses a slim overlay: only the sdd-orchestrator agent, no sub-agents.
	for _, assetPath := range []string{"opencode/sdd-overlay-single.json", "opencode/sdd-overlay-multi.json"} {
		t.Run(assetPath, func(t *testing.T) {
			var root map[string]any
			if err := json.Unmarshal([]byte(MustRead(assetPath)), &root); err != nil {
				t.Fatalf("Unmarshal(%q) error = %v", assetPath, err)
			}

			agents, ok := root["agent"].(map[string]any)
			if !ok {
				t.Fatalf("%q missing agent map", assetPath)
			}

			if _, ok := agents["sdd-orchestrator"]; !ok {
				t.Fatalf("%q missing sdd-orchestrator agent", assetPath)
			}

			// OPSX: no sdd-* sub-agents in the overlay (they were removed in the OPSX migration).
			for _, phase := range []string{"sdd-init", "sdd-explore", "sdd-apply", "sdd-verify"} {
				if _, ok := agents[phase]; ok {
					t.Fatalf("%q should NOT have legacy sub-agent %q (OPSX uses CLI, not sub-agents)", assetPath, phase)
				}
			}
		})
	}
}

func TestSDDOrchestratorAssetsScopedToDedicatedAgent(t *testing.T) {
	for _, assetPath := range []string{
		"generic/sdd-orchestrator.md",
		"claude/sdd-orchestrator.md",
		"gemini/sdd-orchestrator.md",
		"codex/sdd-orchestrator.md",
		"cursor/sdd-orchestrator.md",
	} {
		t.Run(assetPath, func(t *testing.T) {
			content := MustRead(assetPath)
			// Accept both legacy SDD and new OPSX scoping language.
			hasScopingNote := strings.Contains(content, "dedicated `sdd-orchestrator` agent or rule only") ||
				strings.Contains(content, "Bind this to the dedicated `sdd-orchestrator` agent only")
			if !hasScopingNote {
				t.Fatalf("%q missing dedicated-agent scoping note", assetPath)
			}
		})
	}
}
