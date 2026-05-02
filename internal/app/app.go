package app

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/V1merX/artifacts-bot/internal/di"
	"go.uber.org/zap"
)

type app struct {
	di       di.Container
	initOnce *sync.Once
}

func New() *app {
	return &app{
		di:       di.NewContainer(),
		initOnce: &sync.Once{},
	}
}

func (a *app) Run(ctx context.Context) error {
	gracefulCtx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := a.di.Authenticate(gracefulCtx); err != nil {
		return err
	}

	botName := "Bot_2"
	targetLevel := 2

	character, err := a.di.ArtifactsGateway().GetCharacter(gracefulCtx, botName)
	if err != nil {
		a.di.Logger().Error("Failed to get character", zap.String("name", botName), zap.Error(err))
		return err
	}

	a.di.Logger().Info("Character loaded", zap.String("name", character.Name()), zap.Int("level", character.Level()))

	if err := a.di.LevelUpUseCase().Run(gracefulCtx, character.Name(), targetLevel); err != nil {
		a.di.Logger().Error("Failed to level up", zap.Error(err))
		return err
	}

	<-gracefulCtx.Done()
	cancel()

	a.di.Logger().Info("App stopped")
	return nil
}
