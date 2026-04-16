package screens

import (
	"strings"

	"github.com/juancruzrobledo/jr-stack/internal/tui/styles"
)

// RenderUninstallResult shows the outcome of the Gentleman (gentle-ai) uninstall.
func RenderUninstallResult(items []string) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Uninstall Gentleman"))
	b.WriteString("\n\n")

	if len(items) == 0 {
		b.WriteString(styles.SubtextStyle.Render("  No Gentleman (gentle-ai) installation found."))
		b.WriteString("\n")
	} else {
		for _, item := range items {
			b.WriteString(styles.SuccessStyle.Render("  ✓ "+item) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("press enter to return"))

	return b.String()
}
