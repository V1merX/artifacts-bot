package character

import "errors"

var (
	ErrImpossibleLevel = errors.New("impossible level of character")
)

type Level struct {
	value int
}

func NewLevel(inputLevel int) (Level, error) {
	if err := ValidateLevel(inputLevel); err != nil {
		return Level{}, err
	}

	return Level{
		value: inputLevel,
	}, nil
}

func ValidateLevel(lvl int) error {
	if lvl < 0 {
		return ErrImpossibleLevel
	}

	return nil
}

func (s Level) Value() int {
	return s.value
}
