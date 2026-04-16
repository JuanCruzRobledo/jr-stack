package model

type AgentID string

const (
	AgentClaudeCode    AgentID = "claude-code"
	AgentOpenCode      AgentID = "opencode"
	AgentGeminiCLI     AgentID = "gemini-cli"
	AgentCursor        AgentID = "cursor"
	AgentVSCodeCopilot AgentID = "vscode-copilot"
	AgentCodex         AgentID = "codex"
	AgentAntigravity   AgentID = "antigravity"
	AgentWindsurf      AgentID = "windsurf"
)

// SupportTier indicates how fully an agent supports the JR Stack ecosystem.
// All current agents receive the full OPSX orchestrator, skill files, MCP config,
// and system prompt injection. The tier is kept as metadata for display purposes.
type SupportTier string

const (
	// TierFull — the agent receives all ecosystem features: OPSX orchestrator,
	// skill files, MCP servers, system prompt, and sub-agent delegation.
	TierFull SupportTier = "full"
)

type ComponentID string

const (
	ComponentEngram     ComponentID = "engram"
	ComponentSDD        ComponentID = "sdd"
	ComponentSkills     ComponentID = "skills"
	ComponentContext7   ComponentID = "context7"
	ComponentPersona    ComponentID = "persona"
	ComponentPermission ComponentID = "permissions"
	ComponentGGA        ComponentID = "gga"
	ComponentTheme      ComponentID = "theme"
)

type SkillID string

const (
	SkillOpenSpecInit    SkillID = "openspec-init"
	SkillOpenSpecApply   SkillID = "openspec-apply-change"
	SkillOpenSpecVerify  SkillID = "openspec-verify"
	SkillOpenSpecExplore SkillID = "openspec-explore"
	SkillOpenSpecPropose SkillID = "openspec-propose"
	SkillOpenSpecSpec    SkillID = "openspec-spec"
	SkillOpenSpecDesign  SkillID = "openspec-design"
	SkillOpenSpecTasks   SkillID = "openspec-tasks"
	SkillOpenSpecArchive SkillID = "openspec-archive-change"
	SkillOpenSpecOnboard SkillID = "openspec-onboard"
	SkillGoTesting     SkillID = "go-testing"
	SkillCreator       SkillID = "skill-creator"
	SkillJudgmentDay   SkillID = "judgment-day"
	SkillBranchPR      SkillID = "branch-pr"
	SkillIssueCreation SkillID = "issue-creation"
)

type PersonaID string

const (
	PersonaGentleman PersonaID = "gentleman"
	PersonaNeutral   PersonaID = "neutral"
	PersonaCustom    PersonaID = "custom"
)

// SystemPromptStrategy defines how an agent's system prompt file is managed.
type SystemPromptStrategy int

const (
	// StrategyMarkdownSections uses <!-- jr-stack:ID --> markers to inject sections
	// into an existing file without clobbering user content (Claude Code CLAUDE.md).
	StrategyMarkdownSections SystemPromptStrategy = iota
	// StrategyFileReplace replaces the entire system prompt file (OpenCode AGENTS.md).
	StrategyFileReplace
	// StrategyAppendToFile appends content to an existing system prompt file.
	StrategyAppendToFile
)

// MCPStrategy defines how MCP server configs are written for an agent.
type MCPStrategy int

const (
	// StrategySeparateMCPFiles writes one JSON file per server in a dedicated directory
	// (e.g., ~/.claude/mcp/context7.json).
	StrategySeparateMCPFiles MCPStrategy = iota
	// StrategyMergeIntoSettings merges mcpServers into a settings.json file
	// (e.g., OpenCode, Gemini CLI).
	StrategyMergeIntoSettings
	// StrategyMCPConfigFile writes to a dedicated mcp.json config file (e.g., Cursor ~/.cursor/mcp.json).
	StrategyMCPConfigFile
	// StrategyTOMLFile writes MCP config to a TOML file (e.g., Codex ~/.codex/config.toml).
	StrategyTOMLFile
)

type PresetID string

const (
	PresetFullGentleman PresetID = "full-gentleman"
	PresetEcosystemOnly PresetID = "ecosystem-only"
	PresetMinimal       PresetID = "minimal"
	PresetCustom        PresetID = "custom"
)

type SDDModeID string

const (
	SDDModeSingle SDDModeID = "single"
	SDDModeMulti  SDDModeID = "multi"
)
