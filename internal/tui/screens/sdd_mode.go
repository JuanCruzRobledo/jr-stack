package screens

import (
	"strings"

	"github.com/juancruzrobledo/jr-stack/internal/model"
	"github.com/juancruzrobledo/jr-stack/internal/tui/styles"
)

func SDDModeOptions() []model.SDDModeID {
	return []model.SDDModeID{model.SDDModeSingle, model.SDDModeMulti}
}

var sddModeDescriptions = map[model.SDDModeID]string{
	model.SDDModeSingle: "Single orchestrator — one agent coordinates all OPSX actions",
	model.SDDModeMulti:  "Multi-agent — dedicated sub-agent overlay per OPSX action",
}

func RenderSDDMode(selected model.SDDModeID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Select OPSX Mode"))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("How should the OPSX orchestrator be configured for OpenCode?"))
	b.WriteString("\n\n")

	for idx, mode := range SDDModeOptions() {
		isSelected := mode == selected
		focused := idx == cursor
		b.WriteString(renderRadio(string(mode), isSelected, focused))
		b.WriteString(styles.SubtextStyle.Render("    "+sddModeDescriptions[mode]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(renderOptions([]string{"Back"}, cursor-len(SDDModeOptions())))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navigate • enter: select • esc: back"))

	return b.String()
}
