package screens

import (
	"strings"

	"github.com/juancruzrobledo/jr-stack/internal/model"
	"github.com/juancruzrobledo/jr-stack/internal/tui/styles"
)

func PresetOptions() []model.PresetID {
	return []model.PresetID{
		model.PresetLite,
		model.PresetFull,
		model.PresetMinimal,
		model.PresetCustom,
	}
}

var presetDescriptions = map[model.PresetID]string{
	model.PresetLite:    "OPSX essentials: orchestrator + Engram memory + Context7 docs",
	model.PresetFull:    "Complete ecosystem: + skills, GGA, persona & security",
	model.PresetMinimal: "Just Engram persistent memory",
	model.PresetCustom:  "Pick individual components yourself",
}

func RenderPreset(selected model.PresetID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Select Ecosystem Preset"))
	b.WriteString("\n\n")

	for idx, preset := range PresetOptions() {
		isSelected := preset == selected
		focused := idx == cursor
		b.WriteString(renderRadio(string(preset), isSelected, focused))
		b.WriteString(styles.SubtextStyle.Render("    "+presetDescriptions[preset]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(renderOptions([]string{"Back"}, cursor-len(PresetOptions())))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navigate • enter: select • esc: back"))

	return b.String()
}
