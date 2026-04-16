package skills

import "github.com/juancruzrobledo/jr-stack/internal/model"

// openSpecSkills are the OpenSpec orchestrator skills — always included.
var openSpecSkills = []model.SkillID{
	model.SkillOpenSpecInit,
	model.SkillOpenSpecExplore,
	model.SkillOpenSpecPropose,
	model.SkillOpenSpecSpec,
	model.SkillOpenSpecDesign,
	model.SkillOpenSpecTasks,
	model.SkillOpenSpecApply,
	model.SkillOpenSpecVerify,
	model.SkillOpenSpecArchive,
	model.SkillOpenSpecOnboard,
	model.SkillJudgmentDay,
}

// foundationSkills are baseline learning skills for the "recommended" tier.
var foundationSkills = []model.SkillID{
	model.SkillGoTesting,
	model.SkillCreator,
	model.SkillBranchPR,
	model.SkillIssueCreation,
}

// SkillsForPreset returns which skills should be installed for a given preset.
//
//   - "minimal" / PresetMinimal:       OpenSpec skills only
//   - "lite" / PresetLite: OpenSpec + common framework skills
//   - "full" / PresetFull: all available skills
//   - "custom" / PresetCustom:         empty (caller should provide explicit list)
func SkillsForPreset(preset model.PresetID) []model.SkillID {
	switch preset {
	case model.PresetMinimal:
		return copySkills(openSpecSkills)
	case model.PresetLite:
		return copySkills(append(openSpecSkills, foundationSkills...))
	case model.PresetFull:
		all := make([]model.SkillID, 0, len(openSpecSkills)+len(foundationSkills))
		all = append(all, openSpecSkills...)
		all = append(all, foundationSkills...)
		return all
	case model.PresetCustom:
		return nil
	default:
		// Unknown preset — default to full.
		all := make([]model.SkillID, 0, len(openSpecSkills)+len(foundationSkills))
		all = append(all, openSpecSkills...)
		all = append(all, foundationSkills...)
		return all
	}
}

// AllSkillIDs returns every known skill ID.
func AllSkillIDs() []model.SkillID {
	all := make([]model.SkillID, 0, len(openSpecSkills)+len(foundationSkills))
	all = append(all, openSpecSkills...)
	all = append(all, foundationSkills...)
	return all
}

func copySkills(src []model.SkillID) []model.SkillID {
	dst := make([]model.SkillID, len(src))
	copy(dst, src)
	return dst
}
