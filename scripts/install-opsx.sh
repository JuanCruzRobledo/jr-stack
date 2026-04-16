#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# jr-stack OPSX — Install Script
# Compiles the OPSX fork and syncs it, creating the full config from scratch.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/JuanCruzRobledo/jr-stack/main/scripts/install-opsx.sh | bash
#
# Requires: Go 1.24+, git
# ============================================================================

GITHUB_OWNER="JuanCruzRobledo"
GITHUB_REPO="jr-stack"
BINARY_NAME="jr-stack"

# ============================================================================
# Color support
# ============================================================================

setup_colors() {
    if [ -t 1 ] && [ "${TERM:-}" != "dumb" ]; then
        RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
        BLUE='\033[0;34m'; CYAN='\033[0;36m'; BOLD='\033[1m'
        DIM='\033[2m'; NC='\033[0m'
    else
        RED='' GREEN='' YELLOW='' BLUE='' CYAN='' BOLD='' DIM='' NC=''
    fi
}

info()    { echo -e "${BLUE}[info]${NC}    $*"; }
success() { echo -e "${GREEN}[ok]${NC}      $*"; }
warn()    { echo -e "${YELLOW}[warn]${NC}    $*"; }
fatal()   { echo -e "${RED}[error]${NC}   $*" >&2; exit 1; }
step()    { echo -e "\n${CYAN}${BOLD}==>${NC} ${BOLD}$*${NC}"; }

print_banner() {
    echo ""
    echo -e "${CYAN}${BOLD}"
    echo "       _ ____    ____  _             _    "
    echo "      | |  _ \  / ___|| |_ __ _  ___| | __"
    echo "   _  | | |_) | \___ \| __/ _\` |/ __| |/ /"
    echo "  | |_| |  _ <   ___) | || (_| | (__|   < "
    echo "   \___/|_| \_\ |____/ \__\__,_|\___|_|\_\\"
    echo -e "${NC}"
    echo -e "  ${DIM}OPSX Edition — Fluid workflow powered by OpenSpec CLI${NC}"
    echo ""
}

# ============================================================================
# Prerequisites
# ============================================================================

check_prerequisites() {
    step "Checking prerequisites"

    local missing=()
    if ! command -v git &>/dev/null; then missing+=("git"); fi
    if ! command -v go &>/dev/null; then missing+=("go (https://go.dev/dl/)"); fi

    if [ ${#missing[@]} -gt 0 ]; then
        fatal "Missing required tools: ${missing[*]}"
    fi

    local go_version
    go_version="$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')"
    local go_major go_minor
    go_major="$(echo "$go_version" | cut -d. -f1)"
    go_minor="$(echo "$go_version" | cut -d. -f2)"

    if [ "$go_major" -lt 1 ] || { [ "$go_major" -eq 1 ] && [ "$go_minor" -lt 24 ]; }; then
        fatal "Go 1.24+ required, found go${go_version}. Update from https://go.dev/dl/"
    fi

    success "git and Go ${go_version} available"
}

# ============================================================================
# Clone, build, install
# ============================================================================

install_from_source() {
    step "Cloning OPSX fork"

    local tmpdir
    tmpdir="$(mktemp -d)"
    trap '[ -n "${tmpdir:-}" ] && rm -rf "$tmpdir"' EXIT

    git clone --depth 1 "https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}.git" "$tmpdir/repo" 2>&1 | tail -1
    success "Cloned"

    step "Building ${BINARY_NAME}"

    cd "$tmpdir/repo"
    go build -o "${BINARY_NAME}" ./cmd/jr-stack/
    success "Built"

    step "Installing binary"

    # Store install dir globally so run_sync can use the exact path
    if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
        OPSX_INSTALL_DIR="/usr/local/bin"
    else
        OPSX_INSTALL_DIR="${HOME}/.local/bin"
        mkdir -p "$OPSX_INSTALL_DIR"
    fi
    local install_dir="$OPSX_INSTALL_DIR"

    if cp "${BINARY_NAME}" "${install_dir}/${BINARY_NAME}" 2>/dev/null; then
        chmod +x "${install_dir}/${BINARY_NAME}"
    elif command -v sudo &>/dev/null; then
        warn "Permission denied. Trying with sudo..."
        sudo cp "${BINARY_NAME}" "${install_dir}/${BINARY_NAME}"
        sudo chmod +x "${install_dir}/${BINARY_NAME}"
    else
        fatal "Cannot write to ${install_dir}."
    fi

    success "Installed to ${install_dir}/${BINARY_NAME}"

    if [[ ":$PATH:" != *":${install_dir}:"* ]]; then
        warn "${install_dir} is not in your PATH"
        echo -e "  ${DIM}export PATH=\"\$PATH:${install_dir}\"${NC}"
        export PATH="$PATH:${install_dir}"
    fi
}

remove_conflicting_binaries() {
    step "Checking for conflicting jr-stack binaries"

    local opsx_dir="$OPSX_INSTALL_DIR"
    local found_conflict=false

    # Check go/bin (from original `go install`)
    if command -v go &>/dev/null; then
        local gobin
        gobin="$(go env GOBIN 2>/dev/null)"
        if [ -z "$gobin" ]; then
            gobin="$(go env GOPATH 2>/dev/null)/bin"
        fi
        local go_binary="${gobin}/${BINARY_NAME}"
        if [ -f "$go_binary" ] && [ "$gobin" != "$opsx_dir" ]; then
            warn "Removing conflicting binary: $go_binary"
            rm -f "$go_binary" 2>/dev/null || sudo rm -f "$go_binary" 2>/dev/null || warn "Could not remove $go_binary — delete it manually"
            found_conflict=true
        fi
    fi

    # Check all instances in PATH that aren't ours
    local all_matches
    all_matches="$(which -a "$BINARY_NAME" 2>/dev/null || true)"
    if [ -n "$all_matches" ]; then
        while IFS= read -r match; do
            local match_dir
            match_dir="$(dirname "$match")"
            if [ "$match_dir" != "$opsx_dir" ] && [ -f "$match" ]; then
                warn "Removing conflicting binary: $match"
                rm -f "$match" 2>/dev/null || sudo rm -f "$match" 2>/dev/null || warn "Could not remove $match — delete it manually"
                found_conflict=true
            fi
        done <<< "$all_matches"
    fi

    if [ "$found_conflict" = false ]; then
        success "No conflicting binaries found"
    fi
}

set_persistent_no_self_update() {
    step "Disabling self-update permanently"

    export JR_STACK_NO_SELF_UPDATE=1

    # Persist in shell profile
    local shell_rc=""
    if [ -n "${ZSH_VERSION:-}" ] || [ -f "$HOME/.zshrc" ]; then
        shell_rc="$HOME/.zshrc"
    elif [ -f "$HOME/.bashrc" ]; then
        shell_rc="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
        shell_rc="$HOME/.bash_profile"
    fi

    if [ -n "$shell_rc" ]; then
        local export_line='export JR_STACK_NO_SELF_UPDATE=1'
        if ! grep -qF "$export_line" "$shell_rc" 2>/dev/null; then
            echo "" >> "$shell_rc"
            echo "# Prevent jr-stack from self-updating (OPSX fork)" >> "$shell_rc"
            echo "$export_line" >> "$shell_rc"
            success "Added JR_STACK_NO_SELF_UPDATE=1 to $shell_rc"
        else
            success "JR_STACK_NO_SELF_UPDATE already in $shell_rc"
        fi
    else
        warn "Could not detect shell profile. Add manually: export JR_STACK_NO_SELF_UPDATE=1"
    fi
}

# ============================================================================
# Patch state.json to include all detected agents
# ============================================================================

update_agent_state() {
    step "Detecting installed AI agents"

    local state_dir="$HOME/.jr-stack"
    local state_path="$state_dir/state.json"
    local detected=()

    # Claude Code
    if command -v claude &>/dev/null; then
        detected+=("claude-code")
    fi

    # OpenCode
    if command -v opencode &>/dev/null; then
        detected+=("opencode")
    fi

    # Cursor
    if [ -d "$HOME/.cursor" ]; then
        detected+=("cursor")
    fi

    # Windsurf
    if [ -d "$HOME/.windsurf" ]; then
        detected+=("windsurf")
    fi

    # VS Code (Copilot)
    if [ -d "$HOME/.vscode" ]; then
        detected+=("vscode")
    fi

    if [ ${#detected[@]} -eq 0 ]; then
        warn "No AI agents detected - sync may be a no-op"
        return
    fi

    success "Detected agents: ${detected[*]}"

    mkdir -p "$state_dir"

    local json_array=""
    local first=true
    for agent in "${detected[@]}"; do
        if [ "$first" = true ]; then
            json_array="    \"$agent\""
            first=false
        else
            json_array="$json_array,\n    \"$agent\""
        fi
    done

    printf '{\n  "installed_agents": [\n%b\n  ]\n}\n' "$json_array" > "$state_path"
    success "Updated state.json with all detected agents"
}

# ============================================================================
# Clean previous config (so sync creates fresh OPSX)
# ============================================================================

clean_previous_config() {
    step "Cleaning previous config (so sync creates fresh OPSX)"

    local cleaned=false

    # OpenCode: remove opencode.json entirely (sync recreates it with OPSX)
    if [ -f "$HOME/.config/opencode/opencode.json" ]; then
        rm -f "$HOME/.config/opencode/opencode.json"
        cleaned=true
    fi

    # OpenCode: remove old sdd-* commands
    if ls "$HOME/.config/opencode/commands/sdd-"*.md &>/dev/null 2>&1; then
        rm -f "$HOME/.config/opencode/commands/sdd-"*.md
        cleaned=true
    fi

    # Claude Code: remove ALL legacy gentle-ai marker sections so sync injects fresh
    local claude_md="$HOME/.claude/CLAUDE.md"
    if [ -f "$claude_md" ] && grep -q "<!-- gentle-ai:" "$claude_md"; then
        local tmpfile
        tmpfile="$(mktemp)"
        awk '
        /<!-- gentle-ai:/ { skip=1; next }
        /<!-- \/gentle-ai:/ { skip=0; next }
        skip==0 { print }
        ' "$claude_md" > "$tmpfile"
        cp "$tmpfile" "$claude_md"
        rm -f "$tmpfile"
        cleaned=true
    fi

    # Cursor: remove old sdd-* agent files
    if ls "$HOME/.cursor/agents/sdd-"*.md &>/dev/null 2>&1; then
        rm -f "$HOME/.cursor/agents/sdd-"*.md
        cleaned=true
    fi

    if [ "$cleaned" = true ]; then
        success "Previous config cleaned"
    else
        success "No previous config found — clean install"
    fi
}

# ============================================================================
# Installation mode selection
# ============================================================================

select_install_mode() {
    step "Select installation mode"
    echo ""
    echo -e "  ${GREEN}${BOLD}1)${NC} ${BOLD}Lite${NC}    — OPSX essentials ${DIM}(orchestrator + engram + context7)${NC} ${CYAN}← recommended${NC}"
    echo -e "  ${BOLD}2)${NC} ${BOLD}Full${NC}    — Complete ecosystem ${DIM}(+ skills, GGA, persona)${NC}"
    echo -e "  ${BOLD}3)${NC} ${BOLD}Custom${NC}  — Launch TUI to pick components ${DIM}(run 'jr-stack' after install)${NC}"
    echo ""

    # Default to Lite if non-interactive (piped input)
    if [ ! -t 0 ]; then
        INSTALL_MODE="lite"
        info "Non-interactive mode detected — defaulting to Lite"
        return
    fi

    local choice=""
    while true; do
        printf "  ${BOLD}Choose [1/2/3]${NC} (default: 1): "
        read -r choice
        case "${choice:-1}" in
            1) INSTALL_MODE="lite";   success "Lite mode selected"; break ;;
            2) INSTALL_MODE="full";   success "Full mode selected"; break ;;
            3) INSTALL_MODE="custom"; success "Custom mode — run 'jr-stack' after install to configure"; break ;;
            *) warn "Invalid choice. Enter 1, 2, or 3." ;;
        esac
    done
}

# ============================================================================
# Sync (with self-update DISABLED)
# ============================================================================

run_sync() {
    # Custom mode: skip sync entirely — user will run jr-stack TUI.
    if [ "$INSTALL_MODE" = "custom" ]; then
        info "Skipping sync — launch 'jr-stack' to configure via TUI"
        return
    fi

    step "Running OPSX binary sync"

    # Use the EXACT path of our installed binary, NOT whatever is in PATH
    local opsx_binary="${OPSX_INSTALL_DIR}/${BINARY_NAME}"

    if [ ! -x "$opsx_binary" ]; then
        fatal "OPSX binary not found at $opsx_binary"
    fi

    info "Using: $opsx_binary"

    if [ "$INSTALL_MODE" = "lite" ]; then
        "$opsx_binary" sync --lite
        success "Sync complete — OPSX essentials configured (orchestrator + engram + context7)"
    else
        "$opsx_binary" sync
        success "Sync complete — full OPSX ecosystem configured"
    fi
}

# ============================================================================
# Summary
# ============================================================================

print_summary() {
    echo ""
    echo -e "${GREEN}${BOLD}Installation complete!${NC}"
    echo ""

    case "$INSTALL_MODE" in
        lite)
            echo -e "${BOLD}Lite mode:${NC} OPSX orchestrator, Engram memory, and Context7 docs are configured."
            echo -e "${DIM}Run 'jr-stack sync' for full ecosystem, or 'jr-stack' for the TUI.${NC}"
            ;;
        full)
            echo -e "${BOLD}Full mode:${NC} Complete OPSX ecosystem configured for all detected agents."
            ;;
        custom)
            echo -e "${BOLD}Custom mode:${NC} Binary installed. Run ${CYAN}jr-stack${NC} to open the TUI and configure."
            ;;
    esac

    echo ""
    echo -e "${BOLD}OPSX Commands:${NC}"
    echo -e "  ${CYAN}/opsx:explore${NC}  — Think through ideas before committing"
    echo -e "  ${CYAN}/opsx:propose${NC}  — Create a change with all artifacts"
    echo -e "  ${CYAN}/opsx:apply${NC}    — Implement tasks from the change"
    echo -e "  ${CYAN}/opsx:archive${NC}  — Sync specs and close the change"
    echo ""
    echo -e "${DIM}Docs: https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}${NC}"
    echo ""
}

# ============================================================================
# Main
# ============================================================================

main() {
    setup_colors
    print_banner
    check_prerequisites
    install_from_source
    remove_conflicting_binaries
    set_persistent_no_self_update
    update_agent_state
    select_install_mode
    clean_previous_config
    run_sync
    print_summary
}

main "$@"
