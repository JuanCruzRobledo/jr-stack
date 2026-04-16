package agents

import (
	"fmt"

	"github.com/juancruzrobledo/jr-stack/internal/agents/antigravity"
	"github.com/juancruzrobledo/jr-stack/internal/agents/claude"
	"github.com/juancruzrobledo/jr-stack/internal/agents/codex"
	cursoradapter "github.com/juancruzrobledo/jr-stack/internal/agents/cursor"
	"github.com/juancruzrobledo/jr-stack/internal/agents/gemini"
	"github.com/juancruzrobledo/jr-stack/internal/agents/opencode"
	"github.com/juancruzrobledo/jr-stack/internal/agents/vscode"
	"github.com/juancruzrobledo/jr-stack/internal/agents/windsurf"
	"github.com/juancruzrobledo/jr-stack/internal/model"
)

func NewAdapter(agent model.AgentID) (Adapter, error) {
	switch agent {
	case model.AgentClaudeCode:
		return claude.NewAdapter(), nil
	case model.AgentOpenCode:
		return opencode.NewAdapter(), nil
	case model.AgentGeminiCLI:
		return gemini.NewAdapter(), nil
	case model.AgentCursor:
		return cursoradapter.NewAdapter(), nil
	case model.AgentVSCodeCopilot:
		return vscode.NewAdapter(), nil
	case model.AgentCodex:
		return codex.NewAdapter(), nil
	case model.AgentAntigravity:
		return antigravity.NewAdapter(), nil
	case model.AgentWindsurf:
		return windsurf.NewAdapter(), nil
	default:
		return nil, AgentNotSupportedError{Agent: agent}
	}
}

func NewDefaultRegistry() (*Registry, error) {
	adapters := make([]Adapter, 0, 8)

	for _, agent := range []model.AgentID{
		model.AgentClaudeCode,
		model.AgentOpenCode,
		model.AgentGeminiCLI,
		model.AgentCursor,
		model.AgentVSCodeCopilot,
		model.AgentCodex,
		model.AgentAntigravity,
		model.AgentWindsurf,
	} {
		adapter, err := NewAdapter(agent)
		if err != nil {
			return nil, fmt.Errorf("create %s adapter: %w", agent, err)
		}
		adapters = append(adapters, adapter)
	}

	registry, err := NewRegistry(adapters...)
	if err != nil {
		return nil, fmt.Errorf("create registry: %w", err)
	}

	return registry, nil
}

// NewMVPRegistry creates a registry with only the MVP agents (Claude Code, OpenCode).
// Kept for backward compatibility.
func NewMVPRegistry() (*Registry, error) {
	claudeAdapter, err := NewAdapter(model.AgentClaudeCode)
	if err != nil {
		return nil, fmt.Errorf("create claude adapter: %w", err)
	}

	opencodeAdapter, err := NewAdapter(model.AgentOpenCode)
	if err != nil {
		return nil, fmt.Errorf("create opencode adapter: %w", err)
	}

	registry, err := NewRegistry(claudeAdapter, opencodeAdapter)
	if err != nil {
		return nil, fmt.Errorf("create registry: %w", err)
	}

	return registry, nil
}
