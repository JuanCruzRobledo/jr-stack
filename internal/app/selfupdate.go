package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/juancruzrobledo/jr-stack/internal/system"
	"github.com/juancruzrobledo/jr-stack/internal/update"
	"github.com/juancruzrobledo/jr-stack/internal/update/upgrade"
)

// lookPathFn is a package-level var for testability.
var lookPathFn = exec.LookPath

// Environment variable names for self-update control.
const (
	envNoSelfUpdate   = "JR_STACK_NO_SELF_UPDATE"
	envSelfUpdateDone = "JR_STACK_SELF_UPDATE_DONE"
)

// selfUpdateTimeout is the maximum time allowed for the update check + upgrade.
const selfUpdateTimeout = 7 * time.Second

// reExec is swappable for testing — prevents actual syscall.Exec in tests.
var reExec = func(argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}

// goOS returns the current operating system name. Package-level var for testing.
var goOS = func() string { return runtime.GOOS }

// selfUpdate checks for and applies a jr-stack update before normal dispatch.
// Returns nil on success or skip; errors are non-fatal (caller logs and continues).
//
// Guard evaluation order (per spec):
//  1. JR_STACK_SELF_UPDATE_DONE=1 → skip (loop guard)
//  2. JR_STACK_NO_SELF_UPDATE=1 → skip (opt-out)
//  3. version == "dev" → skip (dev build)
//  4. Proceed with update check
func selfUpdate(ctx context.Context, version string, profile system.PlatformProfile, stdout io.Writer) error {
	// Guard 1: loop prevention — already updated this invocation.
	if os.Getenv(envSelfUpdateDone) == "1" {
		return nil
	}

	// Guard 2: user opt-out.
	if os.Getenv(envNoSelfUpdate) == "1" {
		return nil
	}

	// Guard 3: dev build — no meaningful version to compare.
	if version == "dev" {
		return nil
	}

	// Apply timeout to the entire check+upgrade cycle.
	ctx, cancel := context.WithTimeout(ctx, selfUpdateTimeout)
	defer cancel()

	// Check for updates (only jr-stack).
	results := updateCheckFiltered(ctx, version, profile, []string{"jr-stack"})

	// Find the jr-stack result.
	var target *update.UpdateResult
	for i := range results {
		if results[i].Tool.Name == "jr-stack" {
			target = &results[i]
			break
		}
	}

	// No result or not an available update — nothing to do.
	if target == nil || target.Status != update.UpdateAvailable {
		return nil
	}

	// Run upgrade (backup + strategy execution).
	homeDir, err := os.UserHomeDir()
	if err != nil {
		_, _ = fmt.Fprintf(stdout, "self-update: cannot resolve home directory: %v\n", err)
		return nil // non-fatal
	}

	report := upgradeExecute(ctx, results, profile, homeDir, false, stdout)

	// Check if upgrade succeeded.
	var succeeded bool
	for _, r := range report.Results {
		if r.ToolName == "jr-stack" && r.Status == upgrade.UpgradeSucceeded {
			succeeded = true
			break
		}
	}

	if !succeeded {
		// Upgrade failed or was skipped — non-fatal, continue with current binary.
		return nil
	}

	// Re-exec on Unix; print message on Windows.
	if goOS() == "windows" {
		_, _ = fmt.Fprintf(stdout, "Updated to v%s — please restart.\n", target.LatestVersion)
		return nil
	}

	// Unix: re-exec with the updated binary.
	//
	// Use exec.LookPath("jr-stack") rather than os.Executable() because
	// on Homebrew, os.Executable() resolves to the versioned Cellar path
	// (e.g. /opt/homebrew/Cellar/jr-stack/1.8.5/bin/jr-stack) which
	// still points to the OLD binary after upgrade. The PATH symlink
	// (/opt/homebrew/bin/jr-stack) is updated by Homebrew to the new
	// version, so LookPath gives us the correct binary.
	executable, err := lookPathFn("jr-stack")
	if err != nil {
		// Fallback to os.Executable() if LookPath fails.
		executable, err = os.Executable()
		if err != nil {
			return nil // non-fatal
		}
	}

	// Set loop guard env var before re-exec.
	os.Setenv(envSelfUpdateDone, "1")

	_, _ = fmt.Fprintf(stdout, "Updated to v%s, restarting...\n", target.LatestVersion)

	return reExec(executable, os.Args, os.Environ())
}
