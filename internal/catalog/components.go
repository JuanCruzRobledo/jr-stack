package catalog

import "github.com/juancruzrobledo/jr-stack/internal/model"

type Component struct {
	ID          model.ComponentID
	Name        string
	Description string
}

var mvpComponents = []Component{
	{ID: model.ComponentEngram, Name: "Engram", Description: "Persistent cross-session memory"},
	{ID: model.ComponentSDD, Name: "OPSX", Description: "OPSX workflow orchestrator"},
	{ID: model.ComponentSkills, Name: "Skills", Description: "Curated coding skill library"},
	{ID: model.ComponentContext7, Name: "Context7", Description: "Latest framework and library docs"},
	{ID: model.ComponentPersona, Name: "Persona", Description: "Gentleman, neutral or custom behavior"},
	{ID: model.ComponentPermission, Name: "Permissions", Description: "Security-first defaults and guardrails"},
	{ID: model.ComponentGGA, Name: "GGA", Description: "Gentleman Guardian Angel — AI provider switcher"},
	{ID: model.ComponentTheme, Name: "Theme", Description: "Gentleman Kanagawa theme overlay (future)"},
}

func MVPComponents() []Component {
	components := make([]Component, len(mvpComponents))
	copy(components, mvpComponents)
	return components
}
