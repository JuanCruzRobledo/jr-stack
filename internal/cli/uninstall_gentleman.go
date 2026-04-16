package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// UninstallGentleman removes remnants of the old gentle-ai installation:
//   - The gentle-ai binary from common PATH locations
//   - The ~/.gentle-ai/ config directory
//   - GENTLE_AI_NO_SELF_UPDATE lines from shell profiles
//
// It returns a list of human-readable strings describing what was cleaned up.
// A nil error is returned even when nothing was found to remove.
func UninstallGentleman(homeDir string) ([]string, error) {
	var cleaned []string

	// 1. Remove gentle-ai binary from common locations.
	binaryName := "gentle-ai"
	if runtime.GOOS == "windows" {
		binaryName = "gentle-ai.exe"
	}

	searchDirs := binarySearchDirs(homeDir)

	for _, dir := range searchDirs {
		binPath := filepath.Join(dir, binaryName)
		if _, err := os.Stat(binPath); err == nil {
			if err := os.Remove(binPath); err != nil {
				return cleaned, fmt.Errorf("failed to remove binary %s: %w", binPath, err)
			}
			cleaned = append(cleaned, "Removed binary: "+binPath)
		}
	}

	// Also check if gentle-ai is somewhere else on PATH via which/where.
	if extraPath, err := findOnPath(binaryName); err == nil && extraPath != "" {
		// Avoid double-removing if already handled above.
		alreadyRemoved := false
		for _, msg := range cleaned {
			if strings.Contains(msg, extraPath) {
				alreadyRemoved = true
				break
			}
		}
		if !alreadyRemoved {
			if err := os.Remove(extraPath); err == nil {
				cleaned = append(cleaned, "Removed binary: "+extraPath)
			}
		}
	}

	// 2. Remove ~/.gentle-ai/ config directory.
	configDir := filepath.Join(homeDir, ".gentle-ai")
	if info, err := os.Stat(configDir); err == nil && info.IsDir() {
		if err := os.RemoveAll(configDir); err != nil {
			return cleaned, fmt.Errorf("failed to remove config directory %s: %w", configDir, err)
		}
		cleaned = append(cleaned, "Removed config directory: "+configDir)
	}

	// 3. Clean legacy <!-- gentle-ai:* --> marker sections from agent config files.
	agentConfigs := []struct {
		path string
		desc string
	}{
		{filepath.Join(homeDir, ".claude", "CLAUDE.md"), "Claude Code"},
		{filepath.Join(homeDir, ".cursor", "rules", "gentle-ai.mdc"), "Cursor"},
		{filepath.Join(homeDir, ".copilot", "instructions", "gentle-ai.instructions.md"), "VS Code Copilot"},
	}

	for _, cfg := range agentConfigs {
		if content, err := os.ReadFile(cfg.path); err == nil {
			original := string(content)
			stripped := stripLegacyGentleAIMarkers(original)
			if stripped != original {
				if err := os.WriteFile(cfg.path, []byte(stripped), 0644); err == nil {
					cleaned = append(cleaned, "Cleaned legacy markers from: "+cfg.path+" ("+cfg.desc+")")
				}
			}
		}
	}

	// 3b. Remove old sdd-* command files from OpenCode.
	openCodeCmds := filepath.Join(homeDir, ".config", "opencode", "commands")
	if entries, err := os.ReadDir(openCodeCmds); err == nil {
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "sdd-") && strings.HasSuffix(e.Name(), ".md") {
				path := filepath.Join(openCodeCmds, e.Name())
				if err := os.Remove(path); err == nil {
					cleaned = append(cleaned, "Removed legacy command: "+path)
				}
			}
		}
	}

	// 3c. Remove old sdd-* sub-agent files from Cursor.
	cursorAgents := filepath.Join(homeDir, ".cursor", "agents")
	if entries, err := os.ReadDir(cursorAgents); err == nil {
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "sdd-") && strings.HasSuffix(e.Name(), ".md") {
				path := filepath.Join(cursorAgents, e.Name())
				if err := os.Remove(path); err == nil {
					cleaned = append(cleaned, "Removed legacy sub-agent: "+path)
				}
			}
		}
	}

	// 4. Remove GENTLE_AI_NO_SELF_UPDATE from shell profiles.
	shellProfiles := []string{
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".bash_profile"),
	}

	for _, profile := range shellProfiles {
		removed, err := removeEnvLineFromFile(profile, "GENTLE_AI_NO_SELF_UPDATE")
		if err != nil {
			// Non-fatal: log it but continue.
			continue
		}
		if removed {
			cleaned = append(cleaned, "Removed GENTLE_AI_NO_SELF_UPDATE from: "+profile)
		}
	}

	return cleaned, nil
}

// binarySearchDirs returns the common directories where gentle-ai may have been installed.
func binarySearchDirs(homeDir string) []string {
	dirs := []string{
		"/usr/local/bin",
		filepath.Join(homeDir, ".local", "bin"),
	}

	if gopath := os.Getenv("GOPATH"); gopath != "" {
		dirs = append(dirs, filepath.Join(gopath, "bin"))
	} else {
		dirs = append(dirs, filepath.Join(homeDir, "go", "bin"))
	}

	// Windows-specific locations.
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("LOCALAPPDATA"); appData != "" {
			dirs = append(dirs, filepath.Join(appData, "Programs", "gentle-ai"))
		}
		if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
			dirs = append(dirs, filepath.Join(userProfile, "go", "bin"))
		}
	}

	return dirs
}

// findOnPath uses the OS "which" (or "where" on Windows) to locate a binary.
func findOnPath(name string) (string, error) {
	cmd := "which"
	if runtime.GOOS == "windows" {
		cmd = "where"
	}
	out, err := exec.Command(cmd, name).Output()
	if err != nil {
		return "", err
	}
	// Take only the first result line.
	path := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	return path, nil
}

// stripLegacyGentleAIMarkers removes all <!-- gentle-ai:* --> sections from content.
func stripLegacyGentleAIMarkers(content string) string {
	const openPrefix = "<!-- gentle-ai:"
	const closePrefix = "<!-- /gentle-ai:"
	const suffix = " -->"

	for {
		openIdx := strings.Index(content, openPrefix)
		if openIdx < 0 {
			break
		}

		afterOpen := content[openIdx+len(openPrefix):]
		endOfID := strings.Index(afterOpen, suffix)
		if endOfID < 0 {
			break
		}
		sectionID := afterOpen[:endOfID]
		closeMarker := closePrefix + sectionID + suffix

		closeIdx := strings.Index(content[openIdx:], closeMarker)
		if closeIdx < 0 {
			// No matching close — remove the open marker line only.
			lineEnd := strings.Index(content[openIdx:], "\n")
			if lineEnd < 0 {
				content = content[:openIdx]
			} else {
				content = content[:openIdx] + content[openIdx+lineEnd+1:]
			}
			continue
		}

		absClose := openIdx + closeIdx
		after := content[absClose+len(closeMarker):]
		if len(after) > 0 && after[0] == '\n' {
			after = after[1:]
		}
		before := strings.TrimRight(content[:openIdx], "\n")
		if before != "" && after != "" {
			content = before + "\n\n" + after
		} else if after != "" {
			content = after
		} else if before != "" {
			content = before + "\n"
		} else {
			content = ""
		}
	}

	// Collapse triple+ newlines.
	for strings.Contains(content, "\n\n\n") {
		content = strings.ReplaceAll(content, "\n\n\n", "\n\n")
	}
	return content
}

// removeEnvLineFromFile removes lines containing the given env var name from a file.
// Returns true if any lines were removed.
func removeEnvLineFromFile(filePath, envVar string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	var kept []string
	removed := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, envVar) {
			removed = true
			continue
		}
		kept = append(kept, line)
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}

	if !removed {
		return false, nil
	}

	return true, os.WriteFile(filePath, []byte(strings.Join(kept, "\n")+"\n"), 0644)
}

// UninstallJRStack removes the jr-stack installation:
//   - The jr-stack binary from common PATH locations
//   - The ~/.jr-stack/ config directory
//   - JR_STACK_NO_SELF_UPDATE lines from shell profiles
//   - <!-- jr-stack:* --> marker sections from agent config files
//   - opsx-* command files from OpenCode
//
// It returns a list of human-readable strings describing what was cleaned up.
func UninstallJRStack(homeDir string) ([]string, error) {
	var cleaned []string

	// 1. Remove jr-stack binary from common locations.
	binaryName := "jr-stack"
	if runtime.GOOS == "windows" {
		binaryName = "jr-stack.exe"
	}

	searchDirs := jrStackBinarySearchDirs(homeDir)

	for _, dir := range searchDirs {
		binPath := filepath.Join(dir, binaryName)
		if _, err := os.Stat(binPath); err == nil {
			if err := os.Remove(binPath); err != nil {
				return cleaned, fmt.Errorf("failed to remove binary %s: %w", binPath, err)
			}
			cleaned = append(cleaned, "Removed binary: "+binPath)
		}
	}

	// Also check if jr-stack is somewhere else on PATH via which/where.
	if extraPath, err := findOnPath(binaryName); err == nil && extraPath != "" {
		alreadyRemoved := false
		for _, msg := range cleaned {
			if strings.Contains(msg, extraPath) {
				alreadyRemoved = true
				break
			}
		}
		if !alreadyRemoved {
			if err := os.Remove(extraPath); err == nil {
				cleaned = append(cleaned, "Removed binary: "+extraPath)
			}
		}
	}

	// 2. Remove ~/.jr-stack/ config directory.
	configDir := filepath.Join(homeDir, ".jr-stack")
	if info, err := os.Stat(configDir); err == nil && info.IsDir() {
		if err := os.RemoveAll(configDir); err != nil {
			return cleaned, fmt.Errorf("failed to remove config directory %s: %w", configDir, err)
		}
		cleaned = append(cleaned, "Removed config directory: "+configDir)
	}

	// 3. Clean <!-- jr-stack:* --> marker sections from agent config files.
	agentConfigs := []struct {
		path string
		desc string
	}{
		{filepath.Join(homeDir, ".claude", "CLAUDE.md"), "Claude Code"},
		{filepath.Join(homeDir, ".cursor", "rules", "jr-stack.mdc"), "Cursor"},
		{filepath.Join(homeDir, ".copilot", "instructions", "jr-stack.instructions.md"), "VS Code Copilot"},
	}

	for _, cfg := range agentConfigs {
		if content, err := os.ReadFile(cfg.path); err == nil {
			original := string(content)
			stripped := stripJRStackMarkers(original)
			if stripped != original {
				if err := os.WriteFile(cfg.path, []byte(stripped), 0644); err == nil {
					cleaned = append(cleaned, "Cleaned markers from: "+cfg.path+" ("+cfg.desc+")")
				}
			}
		}
	}

	// 3b. Remove opsx-* command files from OpenCode.
	openCodeCmds := filepath.Join(homeDir, ".config", "opencode", "commands")
	if entries, err := os.ReadDir(openCodeCmds); err == nil {
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "opsx-") && strings.HasSuffix(e.Name(), ".md") {
				path := filepath.Join(openCodeCmds, e.Name())
				if err := os.Remove(path); err == nil {
					cleaned = append(cleaned, "Removed command: "+path)
				}
			}
		}
	}

	// 4. Remove JR_STACK_NO_SELF_UPDATE from shell profiles.
	shellProfiles := []string{
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".bash_profile"),
	}

	for _, profile := range shellProfiles {
		removed, err := removeEnvLineFromFile(profile, "JR_STACK_NO_SELF_UPDATE")
		if err != nil {
			continue
		}
		if removed {
			cleaned = append(cleaned, "Removed JR_STACK_NO_SELF_UPDATE from: "+profile)
		}
	}

	return cleaned, nil
}

// jrStackBinarySearchDirs returns the common directories where jr-stack may have been installed.
func jrStackBinarySearchDirs(homeDir string) []string {
	dirs := []string{
		"/usr/local/bin",
		filepath.Join(homeDir, ".local", "bin"),
	}

	if gopath := os.Getenv("GOPATH"); gopath != "" {
		dirs = append(dirs, filepath.Join(gopath, "bin"))
	} else {
		dirs = append(dirs, filepath.Join(homeDir, "go", "bin"))
	}

	if runtime.GOOS == "windows" {
		if appData := os.Getenv("LOCALAPPDATA"); appData != "" {
			dirs = append(dirs, filepath.Join(appData, "Programs", "jr-stack"))
		}
		if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
			dirs = append(dirs, filepath.Join(userProfile, "go", "bin"))
		}
	}

	return dirs
}

// stripJRStackMarkers removes all <!-- jr-stack:* --> sections from content.
func stripJRStackMarkers(content string) string {
	const openPrefix = "<!-- jr-stack:"
	const closePrefix = "<!-- /jr-stack:"
	const suffix = " -->"

	for {
		openIdx := strings.Index(content, openPrefix)
		if openIdx < 0 {
			break
		}

		afterOpen := content[openIdx+len(openPrefix):]
		endOfID := strings.Index(afterOpen, suffix)
		if endOfID < 0 {
			break
		}
		sectionID := afterOpen[:endOfID]
		closeMarker := closePrefix + sectionID + suffix

		closeIdx := strings.Index(content[openIdx:], closeMarker)
		if closeIdx < 0 {
			lineEnd := strings.Index(content[openIdx:], "\n")
			if lineEnd < 0 {
				content = content[:openIdx]
			} else {
				content = content[:openIdx] + content[openIdx+lineEnd+1:]
			}
			continue
		}

		absClose := openIdx + closeIdx
		after := content[absClose+len(closeMarker):]
		if len(after) > 0 && after[0] == '\n' {
			after = after[1:]
		}
		before := strings.TrimRight(content[:openIdx], "\n")
		if before != "" && after != "" {
			content = before + "\n\n" + after
		} else if after != "" {
			content = after
		} else if before != "" {
			content = before + "\n"
		} else {
			content = ""
		}
	}

	for strings.Contains(content, "\n\n\n") {
		content = strings.ReplaceAll(content, "\n\n\n", "\n\n")
	}
	return content
}
