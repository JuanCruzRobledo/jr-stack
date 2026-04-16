package update

import (
	"github.com/juancruzrobledo/jr-stack/internal/system"
)

// updateHint returns a platform-specific instruction string for updating the given tool.
func updateHint(tool ToolInfo, profile system.PlatformProfile) string {
	switch tool.Name {
	case "jr-stack":
		return gentleAIHint(profile)
	case "engram":
		return engramHint(profile)
	case "gga":
		return ggaHint(profile)
	default:
		return ""
	}
}

func gentleAIHint(profile system.PlatformProfile) string {
	switch profile.OS {
	case "darwin":
		return "brew upgrade jr-stack"
	case "linux":
		return "curl -fsSL https://raw.githubusercontent.com/Gentleman-Programming/jr-stack/main/scripts/install.sh | bash"
	case "windows":
		return "irm https://raw.githubusercontent.com/Gentleman-Programming/jr-stack/main/scripts/install.ps1 | iex"
	default:
		return ""
	}
}

func engramHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade engram"
	default:
		return "jr-stack upgrade (downloads pre-built binary)"
	}
}

func ggaHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade gga"
	default:
		return "See https://github.com/Gentleman-Programming/gentleman-guardian-angel"
	}
}
