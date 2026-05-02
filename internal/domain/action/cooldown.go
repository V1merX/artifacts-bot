package action

import (
	"context"
	"time"
)

type Cooldown struct {
	TotalSeconds int
}

func (c Cooldown) Wait(ctx context.Context) error {
	if c.TotalSeconds <= 0 {
		return nil
	}
	select {
	case <-time.After(time.Duration(c.TotalSeconds) * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
