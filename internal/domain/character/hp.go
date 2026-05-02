package character

import "errors"

var (
	ErrImpossibleHP = errors.New("impossible hp of character")
)

type HP struct {
	value int
}

func NewHP(inputHP int) (HP, error) {
	if err := ValidateHP(inputHP); err != nil {
		return HP{}, err
	}

	return HP{
		value: inputHP,
	}, nil
}

func ValidateHP(hp int) error {
	if hp < 0 {
		return ErrImpossibleHP
	}

	return nil
}

func (s HP) Value() int {
	return s.value
}
