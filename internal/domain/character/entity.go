package character

type Character struct {
	name   string
	level  Level
	xp     XP
	maxXP  XP
	hp     HP
	maxHP  HP
	status Status
}

func NewCharacter(
	inputName string,
	inputLevel,
	inputXP,
	inputMaxXP,
	inputHP,
	inputMaxHP int,
	inputStatus Status,
) (*Character, error) {
	level, err := NewLevel(inputLevel)
	if err != nil {
		return nil, err
	}

	xp, err := NewXP(inputXP)
	if err != nil {
		return nil, err
	}

	maxXP, err := NewXP(inputMaxXP)
	if err != nil {
		return nil, err
	}

	hp, err := NewHP(inputHP)
	if err != nil {
		return nil, err
	}

	maxHP, err := NewHP(inputMaxHP)
	if err != nil {
		return nil, err
	}

	return &Character{
		name:   inputName,
		level:  level,
		xp:     xp,
		maxXP:  maxXP,
		hp:     hp,
		maxHP:  maxHP,
		status: inputStatus,
	}, nil
}

func (c *Character) Name() string       { return c.name }
func (c *Character) Level() int         { return c.level.Value() }
func (c *Character) XP() int            { return c.xp.Value() }
func (c *Character) MaxXP() int         { return c.maxXP.Value() }
func (c *Character) HP() int            { return c.hp.Value() }
func (c *Character) MaxHP() int         { return c.maxHP.Value() }
func (c *Character) Status() statusType { return c.status.Value() }

const restHPThresholdPercent = 50.0

func (c *Character) NeedsRest() bool {
	if c.maxHP.Value() == 0 {
		return false
	}
	return float64(c.hp.Value())/float64(c.maxHP.Value())*100 <= restHPThresholdPercent
}

func (c *Character) HasReachedLevel(target int) bool {
	return c.level.Value() >= target
}
