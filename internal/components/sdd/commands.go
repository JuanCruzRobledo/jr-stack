package sdd

type OpenCodeCommand struct {
	Name        string
	Description string
	Body        string
}

func OpenCodeCommands() []OpenCodeCommand {
	return []OpenCodeCommand{
		{Name: "opsx-explore", Description: "Explore mode — think through ideas", Body: "/opsx:explore ${topic}"},
		{Name: "opsx-propose", Description: "Propose a new change with all artifacts", Body: "/opsx:propose ${change-name}"},
		{Name: "opsx-apply", Description: "Implement tasks from a change", Body: "/opsx:apply ${change-name}"},
		{Name: "opsx-archive", Description: "Archive a completed change", Body: "/opsx:archive ${change-name}"},
	}
}
