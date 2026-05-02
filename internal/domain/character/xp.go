package character

import "errors"

var (
	ErrImpossibleXP = errors.New("impossible xp of character")
)

type XP struct {
	value int
}

func NewXP(inputXP int) (XP, error) {
	if err := ValidateXP(inputXP); err != nil {
		return XP{}, err
	}

	return XP{
		value: inputXP,
	}, nil
}

func ValidateXP(xp int) error {
	if xp < 0 {
		return ErrImpossibleXP
	}

	return nil
}

func (s XP) Value() int {
	return s.value
}
