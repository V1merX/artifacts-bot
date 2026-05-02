package action

const restHPThresholdPercent = 50.0

type CharacterSnapshot struct {
	HP    int
	MaxHP int
	Level int
}

func (s CharacterSnapshot) NeedsRest() bool {
	if s.MaxHP == 0 {
		return false
	}
	return float64(s.HP)/float64(s.MaxHP)*100 <= restHPThresholdPercent
}

func (s CharacterSnapshot) HasReachedLevel(target int) bool {
	return s.Level >= target
}

type FightResult struct {
	Cooldown   Cooldown
	Characters []CharacterSnapshot
}
