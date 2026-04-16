#Requires -Version 5.1
<#
.SYNOPSIS
    jr-stack OPSX - Install Script for Windows
    Compiles the OPSX fork and syncs it, creating the full config from scratch.

.DESCRIPTION
    Clones the OPSX fork, builds from source, and runs sync with self-update
    disabled so the official binary does not overwrite OPSX changes.
    Requires Go 1.24+ and git.

.EXAMPLE
    irm https://raw.githubusercontent.com/JuanCruzRobledo/jr-stack/main/scripts/install-opsx.ps1 | iex
#>

$ErrorActionPreference = "Stop"

$GITHUB_OWNER = "JuanCruzRobledo"
$GITHUB_REPO = "jr-stack"
$BINARY_NAME = "jr-stack"

# ============================================================================
# Logging helpers
# ============================================================================

function Write-Info    { param([string]$Message) Write-Host "[info]    $Message" -ForegroundColor Blue }
function Write-Success { param([string]$Message) Write-Host "[ok]      $Message" -ForegroundColor Green }
function Write-Warn    { param([string]$Message) Write-Host "[warn]    $Message" -ForegroundColor Yellow }
function Write-Err     { param([string]$Message) Write-Host "[error]   $Message" -ForegroundColor Red }
function Write-Step    { param([string]$Message) Write-Host "`n==> $Message" -ForegroundColor Cyan }

function Stop-WithError {
    param([string]$Message)
    Write-Err $Message
    exit 1
}

function Show-Banner {
    Write-Host ""
    Write-Host "   ____            _   _              _    ___ " -ForegroundColor Cyan
    Write-Host "  / ___| ___ _ __ | |_| | ___        / \  |_ _|" -ForegroundColor Cyan
    Write-Host " | |  _ / _ \ '_ \| __| |/ _ \_____ / _ \  | | " -ForegroundColor Cyan
    Write-Host " | |_| |  __/ | | | |_| |  __/_____/ ___ \ | | " -ForegroundColor Cyan
    Write-Host "  \____|\___|_| |_|\__|_|\___|    /_/   \_\___|" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  OPSX Edition - Fluid workflow powered by OpenSpec CLI" -ForegroundColor DarkGray
    Write-Host ""
}

# ============================================================================
# Prerequisites
# ============================================================================

function Test-Prerequisites {
    Write-Step "Checking prerequisites"

    $missing = @()
    if (-not (Get-Command "git" -ErrorAction SilentlyContinue)) { $missing += "git" }
    if (-not (Get-Command "go" -ErrorAction SilentlyContinue))  { $missing += "go (https://go.dev/dl/)" }

    if ($missing.Count -gt 0) {
        Stop-WithError "Missing required tools: $($missing -join ', ')"
    }

    $goVersionOutput = & go version 2>&1
    if ($goVersionOutput -match "go(\d+)\.(\d+)") {
        $goMajor = [int]$Matches[1]
        $goMinor = [int]$Matches[2]
        if ($goMajor -lt 1 -or ($goMajor -eq 1 -and $goMinor -lt 24)) {
            Stop-WithError "Go 1.24+ required, found go${goMajor}.${goMinor}. Update from https://go.dev/dl/"
        }
    }

    Write-Success "git and Go available"
}

# ============================================================================
# Remove conflicting binaries
# ============================================================================

function Remove-ConflictingBinaries {
    Write-Step "Checking for conflicting jr-stack binaries"

    $opsxInstallDir = Join-Path $env:LOCALAPPDATA "jr-stack\bin"
    $conflicts = @()

    # Check go/bin (from original go install)
    if (Get-Command "go" -ErrorAction SilentlyContinue) {
        $gobin = & go env GOBIN 2>$null
        if (-not $gobin) {
            $gopath = & go env GOPATH 2>$null
            $gobin = Join-Path $gopath "bin"
        }
        $goBinary = Join-Path $gobin "$BINARY_NAME.exe"
        if ((Test-Path $goBinary) -and ($gobin -ne $opsxInstallDir)) {
            $conflicts += $goBinary
        }
    }

    # Check any other jr-stack.exe in PATH that is not ours
    $allMatches = where.exe $BINARY_NAME 2>$null
    if ($allMatches) {
        foreach ($match in $allMatches) {
            $matchDir = Split-Path $match -Parent
            if ($matchDir -ne $opsxInstallDir -and (Test-Path $match)) {
                if ($conflicts -notcontains $match) { $conflicts += $match }
            }
        }
    }

    if ($conflicts.Count -eq 0) {
        Write-Success "No conflicting binaries found"
        return
    }

    foreach ($conflict in $conflicts) {
        Write-Warn "Removing conflicting binary: $conflict"
        Remove-Item $conflict -Force -ErrorAction SilentlyContinue
        if (-not (Test-Path $conflict)) {
            Write-Success "Removed $conflict"
        } else {
            Write-Warn "Could not remove $conflict - delete it manually to avoid PATH conflicts"
        }
    }
}

# ============================================================================
# Clone, build, install
# ============================================================================

function Install-FromSource {
    Write-Step "Cloning OPSX fork"

    $script:TmpDir = Join-Path $env:TEMP "jr-stack-opsx-$(Get-Random)"
    New-Item -ItemType Directory -Path $script:TmpDir -Force | Out-Null

    $ErrorActionPreference = "Continue"
    & git clone --depth 1 "https://github.com/$GITHUB_OWNER/$GITHUB_REPO.git" "$script:TmpDir\repo" 2>$null
    $ErrorActionPreference = "Stop"
    if ($LASTEXITCODE -ne 0) { Stop-WithError "Failed to clone repository" }
    Write-Success "Cloned"

    Write-Step "Building $BINARY_NAME"

    Push-Location "$script:TmpDir\repo"
    & go build -o "$BINARY_NAME.exe" ./cmd/jr-stack/
    if ($LASTEXITCODE -ne 0) { Stop-WithError "Build failed" }
    Pop-Location
    Write-Success "Built"

    Write-Step "Installing binary"

    $installDir = Join-Path $env:LOCALAPPDATA "jr-stack\bin"
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }

    $destPath = Join-Path $installDir "$BINARY_NAME.exe"
    Copy-Item -Path "$script:TmpDir\repo\$BINARY_NAME.exe" -Destination $destPath -Force
    Write-Success "Installed to $destPath"

    if ($env:PATH -notlike "*$installDir*") {
        Write-Info "Adding to user PATH..."
        $currentPath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::User)
        if ($currentPath -notlike "*$installDir*") {
            [Environment]::SetEnvironmentVariable("PATH", ($currentPath + ";" + $installDir), [EnvironmentVariableTarget]::User)
        }
        $env:PATH = $env:PATH + ";" + $installDir
        Write-Success "Added to PATH"
    }
}

# ============================================================================
# Disable self-update permanently
# ============================================================================

function Set-PersistentNoSelfUpdate {
    Write-Step "Disabling self-update permanently (prevents official binary from overwriting OPSX)"

    $currentValue = [Environment]::GetEnvironmentVariable("JR_STACK_NO_SELF_UPDATE", [EnvironmentVariableTarget]::User)
    if ($currentValue -ne "1") {
        [Environment]::SetEnvironmentVariable("JR_STACK_NO_SELF_UPDATE", "1", [EnvironmentVariableTarget]::User)
        Write-Success "JR_STACK_NO_SELF_UPDATE=1 set as persistent user variable"
    } else {
        Write-Success "JR_STACK_NO_SELF_UPDATE already set"
    }

    $env:JR_STACK_NO_SELF_UPDATE = "1"
}

# ============================================================================
# Patch state.json to include all detected agents
# ============================================================================

function Update-AgentState {
    Write-Step "Detecting installed AI agents"

    $stateDir = Join-Path $HOME ".jr-stack"
    $statePath = Join-Path $stateDir "state.json"

    $detectedAgents = @()

    # Claude Code
    if (Get-Command "claude" -ErrorAction SilentlyContinue) {
        $detectedAgents += "claude-code"
    }

    # OpenCode
    if (Get-Command "opencode" -ErrorAction SilentlyContinue) {
        $detectedAgents += "opencode"
    }

    # Cursor
    $cursorDir = Join-Path $HOME ".cursor"
    if (Test-Path $cursorDir) {
        $detectedAgents += "cursor"
    }

    # Windsurf
    $windsurfDir = Join-Path $HOME ".windsurf"
    if (Test-Path $windsurfDir) {
        $detectedAgents += "windsurf"
    }

    # VS Code (Copilot)
    $vscodeDir = Join-Path $HOME ".vscode"
    if (Test-Path $vscodeDir) {
        $detectedAgents += "vscode"
    }

    if ($detectedAgents.Count -eq 0) {
        Write-Warn "No AI agents detected - sync may be a no-op"
        return
    }

    Write-Success ("Detected agents: " + ($detectedAgents -join ", "))

    if (-not (Test-Path $stateDir)) {
        New-Item -ItemType Directory -Path $stateDir -Force | Out-Null
    }

    $agentsJson = ($detectedAgents | ForEach-Object { "    `"$_`"" }) -join ",`n"
    $stateContent = "{`n  `"installed_agents`": [`n$agentsJson`n  ]`n}"
    Set-Content -Path $statePath -Value $stateContent -Encoding UTF8
    Write-Success "Updated state.json with all detected agents"
}

# ============================================================================
# Clean previous config
# ============================================================================

function Clear-LegacyConfig {
    Write-Step "Cleaning previous config (so sync creates fresh OPSX)"

    $cleaned = $false

    # OpenCode: remove opencode.json entirely (sync recreates it with OPSX)
    $ocJson = Join-Path $HOME ".config\opencode\opencode.json"
    if (Test-Path $ocJson) {
        Remove-Item $ocJson -Force
        $cleaned = $true
    }

    # OpenCode: remove old sdd-* commands
    $ocCmds = Join-Path $HOME ".config\opencode\commands"
    if (Test-Path $ocCmds) {
        $sddFiles = Get-ChildItem -Path $ocCmds -Filter "sdd-*.md" -ErrorAction SilentlyContinue
        if ($sddFiles) { $sddFiles | Remove-Item -Force; $cleaned = $true }
    }

    # Claude Code: remove sdd-orchestrator section marker so sync injects fresh
    $claudeMd = Join-Path $HOME ".claude\CLAUDE.md"
    if (Test-Path $claudeMd) {
        $content = Get-Content -Path $claudeMd -Raw
        $openTag = "<!-- jr-stack:sdd-orchestrator -->"
        $closeTag = "<!-- /jr-stack:sdd-orchestrator -->"
        if ($content -match [regex]::Escape($openTag)) {
            $escapedOpen = [regex]::Escape($openTag)
            $escapedClose = [regex]::Escape($closeTag)
            $regexPattern = "(?s)" + $escapedOpen + ".*?" + $escapedClose + "\r?\n?"
            $content = [regex]::Replace($content, $regexPattern, "")
            Set-Content -Path $claudeMd -Value $content -NoNewline
            $cleaned = $true
        }
    }

    # Cursor: remove old sdd-* agent files
    $cursorAgents = Join-Path $HOME ".cursor\agents"
    if (Test-Path $cursorAgents) {
        $sddFiles = Get-ChildItem -Path $cursorAgents -Filter "sdd-*.md" -ErrorAction SilentlyContinue
        if ($sddFiles) { $sddFiles | Remove-Item -Force; $cleaned = $true }
    }

    # Claude Code: remove old sdd-* skill folders (replaced by openspec-*)
    $claudeSkills = Join-Path $HOME ".claude\skills"
    if (Test-Path $claudeSkills) {
        $sddSkillDirs = Get-ChildItem -Path $claudeSkills -Directory -Filter "sdd-*" -ErrorAction SilentlyContinue
        if ($sddSkillDirs) {
            $sddSkillDirs | Remove-Item -Recurse -Force
            $cleaned = $true
            Write-Info "Removed $($sddSkillDirs.Count) old sdd-* skill folders from Claude"
        }
        # Remove core OPSX skills that are now installed as /opsx:* commands
        # (they caused duplicates when both skill and command existed)
        foreach ($dupSkill in @("openspec-explore", "openspec-propose", "openspec-apply-change", "openspec-archive-change")) {
            $dupPath = Join-Path $claudeSkills $dupSkill
            if (Test-Path $dupPath) {
                Remove-Item $dupPath -Recurse -Force
                $cleaned = $true
            }
        }
    }

    # Claude Code: remove old flat opsx-*.md command files (replaced by opsx/*.md)
    $claudeCommands = Join-Path $HOME ".claude\commands"
    if (Test-Path $claudeCommands) {
        $flatCmds = Get-ChildItem -Path $claudeCommands -Filter "opsx-*.md" -ErrorAction SilentlyContinue
        if ($flatCmds) {
            $flatCmds | Remove-Item -Force
            $cleaned = $true
            Write-Info "Removed $($flatCmds.Count) old flat opsx-*.md command files from Claude"
        }
    }

    # All agents: remove old sdd-* skill folders from all known skill dirs
    $skillDirs = @(
        (Join-Path $HOME ".config\opencode\skills"),
        (Join-Path $HOME ".gemini\skills"),
        (Join-Path $HOME ".gemini\antigravity\skills"),
        (Join-Path $HOME ".copilot\skills"),
        (Join-Path $HOME ".codex\skills"),
        (Join-Path $HOME ".cursor\skills"),
        (Join-Path $HOME ".codeium\windsurf\skills")
    )
    foreach ($dir in $skillDirs) {
        if (Test-Path $dir) {
            $sddSkillDirs = Get-ChildItem -Path $dir -Directory -Filter "sdd-*" -ErrorAction SilentlyContinue
            if ($sddSkillDirs) {
                $sddSkillDirs | Remove-Item -Recurse -Force
                $cleaned = $true
            }
        }
    }

    if ($cleaned) {
        Write-Success "Previous config cleaned"
    } else {
        Write-Success "No previous config found - clean install"
    }
}

# ============================================================================
# Sync
# ============================================================================

function Invoke-Sync {
    Write-Step "Running OPSX binary sync"

    $opsxBinary = Join-Path $env:LOCALAPPDATA "jr-stack\bin\$BINARY_NAME.exe"

    if (Test-Path $opsxBinary) {
        Write-Info "Using: $opsxBinary"
        & $opsxBinary sync
        Write-Success "Sync complete - OPSX config created"
    } else {
        Stop-WithError "OPSX binary not found at $opsxBinary"
    }
}

# ============================================================================
# Summary
# ============================================================================

function Show-Summary {
    Write-Host ""
    Write-Host "Installation complete!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Your agents are configured with OPSX." -ForegroundColor White
    Write-Host ""
    Write-Host "OPSX Commands:" -ForegroundColor White
    Write-Host "  /opsx:explore  - Think through ideas before committing" -ForegroundColor Cyan
    Write-Host "  /opsx:propose  - Create a change with all artifacts" -ForegroundColor Cyan
    Write-Host "  /opsx:apply    - Implement tasks from the change" -ForegroundColor Cyan
    Write-Host "  /opsx:archive  - Sync specs and close the change" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "PROTECTED: Self-update is permanently disabled so jr-stack sync" -ForegroundColor Yellow
    Write-Host "  will always use your OPSX binary. To undo: remove the" -ForegroundColor DarkGray
    Write-Host "  JR_STACK_NO_SELF_UPDATE user environment variable." -ForegroundColor DarkGray
    Write-Host ""
    Write-Host ("Docs: https://github.com/" + $GITHUB_OWNER + "/" + $GITHUB_REPO) -ForegroundColor DarkGray
    Write-Host ""
}

# ============================================================================
# Main
# ============================================================================

function Main {
    Show-Banner
    Test-Prerequisites
    Install-FromSource
    Remove-ConflictingBinaries
    Set-PersistentNoSelfUpdate
    Update-AgentState
    Clear-LegacyConfig
    Invoke-Sync
    Show-Summary

    # Cleanup
    if ($script:TmpDir -and (Test-Path $script:TmpDir)) {
        Remove-Item -Path $script:TmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

Main
