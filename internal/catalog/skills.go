package catalog

import "github.com/juancruzrobledo/jr-stack/internal/model"

type Skill struct {
	ID       model.SkillID
	Name     string
	Category string
	Priority string
}

var mvpSkills = []Skill{
	// OpenSpec skills
	{ID: model.SkillOpenSpecInit, Name: "openspec-init", Category: "openspec", Priority: "p0"},

	{ID: model.SkillOpenSpecApply, Name: "openspec-apply-change", Category: "openspec", Priority: "p0"},
	{ID: model.SkillOpenSpecVerify, Name: "openspec-verify", Category: "openspec", Priority: "p0"},
	{ID: model.SkillOpenSpecExplore, Name: "openspec-explore", Category: "openspec", Priority: "p0"},
	{ID: model.SkillOpenSpecPropose, Name: "openspec-propose", Category: "openspec", Priority: "p0"},
	{ID: model.SkillOpenSpecSpec, Name: "openspec-spec", Category: "openspec", Priority: "p0"},
	{ID: model.SkillOpenSpecDesign, Name: "openspec-design", Category: "openspec", Priority: "p0"},
	{ID: model.SkillOpenSpecTasks, Name: "openspec-tasks", Category: "openspec", Priority: "p0"},
	{ID: model.SkillOpenSpecArchive, Name: "openspec-archive-change", Category: "openspec", Priority: "p0"},
	{ID: model.SkillOpenSpecOnboard, Name: "openspec-onboard", Category: "openspec", Priority: "p0"},
	// Foundation skills
	{ID: model.SkillGoTesting, Name: "go-testing", Category: "testing", Priority: "p0"},
	{ID: model.SkillCreator, Name: "skill-creator", Category: "workflow", Priority: "p0"},
	{ID: model.SkillJudgmentDay, Name: "judgment-day", Category: "workflow", Priority: "p0"},
	{ID: model.SkillBranchPR, Name: "branch-pr", Category: "workflow", Priority: "p0"},
	{ID: model.SkillIssueCreation, Name: "issue-creation", Category: "workflow", Priority: "p0"},
}

func MVPSkills() []Skill {
	skills := make([]Skill, len(mvpSkills))
	copy(skills, mvpSkills)
	return skills
}
