package engram

import (
	"github.com/juancruzrobledo/jr-stack/internal/installcmd"
	"github.com/juancruzrobledo/jr-stack/internal/model"
	"github.com/juancruzrobledo/jr-stack/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentEngram)
}
