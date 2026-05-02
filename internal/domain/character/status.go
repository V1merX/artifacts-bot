package character

import "errors"

type statusType string

const (
	MoveStatusType   statusType = "MOVE"
	FightStatusType  statusType = "FIGHT"
	LoadedStatusType statusType = "LOADED"
)

var (
	ErrInvalidStatusType = errors.New("invalid status type")
)

type Status struct {
	value statusType
}

func NewStatus(inputStatus statusType) (Status, error) {
	if err := ValidateStatus(inputStatus); err != nil {
		return Status{}, err
	}

	return Status{
		value: inputStatus,
	}, nil
}

func ValidateStatus(status statusType) error {
	switch status {
	case MoveStatusType:
		return nil
	case FightStatusType:
		return nil
	case LoadedStatusType:
		return nil
	default:
		return ErrInvalidStatusType
	}
}

func (s Status) Value() statusType {
	return s.value
}
