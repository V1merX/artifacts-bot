package levelup

import (
	"context"
	"errors"

	"github.com/V1merX/artifacts-bot/internal/domain/action"
	"github.com/V1merX/artifacts-bot/internal/domain/character"
	"go.uber.org/zap"
)

type Gateway interface {
	GetCharacter(ctx context.Context, name string) (*character.Character, error)
	Fight(ctx context.Context, characterName string) (*action.FightResult, error)
	Move(ctx context.Context, characterName string, x, y int) (*action.MoveResult, error)
	Rest(ctx context.Context, characterName string) (*action.RestResult, error)
}

type UseCase struct {
	gw     Gateway
	logger *zap.Logger
}

func New(gw Gateway, logger *zap.Logger) *UseCase {
	return &UseCase{gw: gw, logger: logger}
}

func (uc *UseCase) Run(ctx context.Context, characterName string, targetLevel int) error {
	// Move to spawn
	moveToSpawnResult, err := uc.gw.Move(ctx, characterName, 0, 0)
	if err != nil {
		uc.logger.Error("Failed to move to spawn", zap.Error(err))
	}

	if err := moveToSpawnResult.Cooldown.Wait(ctx); err != nil {
		return err
	}

	uc.logger.Info("Move to spawn ended", zap.Int("cooldown_seconds", moveToSpawnResult.Cooldown.TotalSeconds))

	moveResult, err := uc.gw.Move(ctx, characterName, 0, 1)
	if err != nil {
		return err
	}

	if moveResult.Character.NeedsRest() {
		return errors.New("HP below 50% before fight loop started")
	}

	if err := moveResult.Cooldown.Wait(ctx); err != nil {
		return err
	}

	uc.logger.Info("Move ended",
		zap.Int("cooldown_seconds", moveResult.Cooldown.TotalSeconds),
		zap.Int("x", 0),
		zap.Int("y", 1),
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fightResult, err := uc.gw.Fight(ctx, characterName)
		if err != nil {
			return err
		}

		if err := fightResult.Cooldown.Wait(ctx); err != nil {
			return err
		}

		uc.logger.Info("Fight ended", zap.Int("cooldown_seconds", fightResult.Cooldown.TotalSeconds))

		if len(fightResult.Characters) == 0 {
			return errors.New("fight result contained no characters")
		}

		character := fightResult.Characters[0]

		if character.NeedsRest() {
			restResult, err := uc.gw.Rest(ctx, characterName)
			if err != nil {
				return err
			}

			if err := restResult.Cooldown.Wait(ctx); err != nil {
				return err
			}

			uc.logger.Info("Rest ended", zap.Int("cooldown_seconds", restResult.Cooldown.TotalSeconds))
		}

		if character.HasReachedLevel(targetLevel) {
			uc.logger.Info("Character has reached target level", zap.Int("level", character.Level))
			return nil
		}
	}
}
