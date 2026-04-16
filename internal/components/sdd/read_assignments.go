package sdd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/juancruzrobledo/jr-stack/internal/model"
	"github.com/juancruzrobledo/jr-stack/internal/opencode"
)

// sddPhaseSet is the set of valid SDD phase agent names that may appear in
// opencode.json. It includes the current OPSX actions, the orchestrator, and
// legacy phase names (sdd-init, sdd-spec, sdd-design, sdd-tasks, sdd-verify)
// for backwards compatibility with existing user configs.
var sddPhaseSet = buildSDDPhaseSet()

func buildSDDPhaseSet() map[string]bool {
	phases := opencode.SDDPhases()
	set := make(map[string]bool, len(phases)+6) // current phases + orchestrator + 5 legacy
	for _, p := range phases {
		set[p] = true
	}
	set["sdd-orchestrator"] = true
	// Legacy phases — still accepted from existing configs but no longer offered in TUI.
	set["sdd-init"] = true
	set["sdd-spec"] = true
	set["sdd-design"] = true
	set["sdd-tasks"] = true
	set["sdd-verify"] = true
	return set
}

// ReadCurrentModelAssignments reads the agent definitions from opencode.json
// at settingsPath and extracts the "model" field for each SDD phase agent.
//
// Only agents whose names match an SDD phase (from opencode.SDDPhases()) or
// "sdd-orchestrator" are included. Agents without a "model" field, or with a
// malformed model value (not in "provider:model-id" format), are silently
// skipped.
//
// Returns an empty map (no error) when the file does not exist, contains no
// "agent" key, or has no matching phase agents with a valid model field.
func ReadCurrentModelAssignments(settingsPath string) (map[string]model.ModelAssignment, error) {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]model.ModelAssignment{}, nil
		}
		return nil, err
	}

	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		// Unparseable JSON — return empty map, no error.
		return map[string]model.ModelAssignment{}, nil
	}

	agentRaw, ok := root["agent"]
	if !ok {
		return map[string]model.ModelAssignment{}, nil
	}
	agentMap, ok := agentRaw.(map[string]any)
	if !ok {
		return map[string]model.ModelAssignment{}, nil
	}

	result := make(map[string]model.ModelAssignment)
	for name, defRaw := range agentMap {
		if !sddPhaseSet[name] {
			continue
		}
		defMap, ok := defRaw.(map[string]any)
		if !ok {
			continue
		}
		modelStr, ok := defMap["model"].(string)
		if !ok || modelStr == "" {
			continue
		}
		// Try colon first (standard: "anthropic:claude-sonnet-4"), then slash
		// ("zai-coding-plan/glm-5-turbo") for custom providers (issue #152).
		idx := strings.Index(modelStr, ":")
		if idx <= 0 {
			idx = strings.Index(modelStr, "/")
		}
		if idx <= 0 {
			// No separator or separator is the first character — skip malformed value.
			continue
		}
		providerID := modelStr[:idx]
		modelID := modelStr[idx+1:]
		if modelID == "" {
			continue
		}
		result[name] = model.ModelAssignment{
			ProviderID: providerID,
			ModelID:    modelID,
		}
	}

	return result, nil
}
