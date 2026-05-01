package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/V1merX/artifacts-bot/internal/di"
	"github.com/V1merX/artifacts-bot/pkg/api"
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

	a.initDeps()

	resp, err := a.di.ArtifactClient().GenerateTokenTokenPost(gracefulCtx)
	if err != nil {
		a.di.Logger().Error("Failed to generate new token", zap.Error(err))
		return err
	}

	var genTokenResp api.TokenResponseSchema
	err = json.NewDecoder(resp.Body).Decode(&genTokenResp)
	if err != nil {
		a.di.Logger().Error("Failed to generate new token", zap.Error(err))
		return err
	}

	a.di.ArtifactClient().RequestEditors = append(a.di.ArtifactClient().RequestEditors, func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", genTokenResp.Token))
		return nil
	})

	a.di.Logger().Info("Successful authorization with token", zap.String("token", fmt.Sprintf("%s...", genTokenResp.Token[:55])))

	character := "Bot_3"
	lvl := 2
	if err := a.UPLevel(gracefulCtx, character, lvl); err != nil {
		a.di.Logger().Error("Failed to up character level", zap.Error(err), zap.String("character_name", character))
		return err
	}

	a.di.Logger().Info("Successful level-up", zap.String("character_name", character), zap.Int("level", lvl))

	<-gracefulCtx.Done()
	cancel()

	a.di.Logger().Info("Stopping the app...")

	a.di.Logger().Info("App has been stopped")

	return nil
}

func (a *app) UPLevel(ctx context.Context, character string, lvl int) error {
	// Move character to base monster
	moveResp, err := a.MoveCharacter(ctx, character, 0, 1)
	if err != nil {
		return err
	}

	if float64(moveResp.Character.Hp)/float64(moveResp.Character.MaxHp)*100 <= 50.0 {
		return errors.New("HP less than 50%")
	}

	time.Sleep(time.Duration(moveResp.Cooldown.TotalSeconds) * time.Second)

	a.di.Logger().Info("Move has been ended", zap.Int("cooldown_total_seconds", moveResp.Cooldown.TotalSeconds))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fightResponse, err := a.FightCharacter(ctx, character)
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(fightResponse.Cooldown.TotalSeconds) * time.Second)

		a.di.Logger().Info("Fight has been ended", zap.Int("cooldown_total_seconds", fightResponse.Cooldown.TotalSeconds))

		if float64(fightResponse.Characters[0].Hp)/float64(fightResponse.Characters[0].MaxHp)*100 <= 50.0 {
			restResponse, err := a.RestCharacter(ctx, character)
			if err != nil {
				return err
			}

			time.Sleep(time.Duration(restResponse.Cooldown.TotalSeconds) * time.Second)

			a.di.Logger().Info("Rest has been ended", zap.Int("cooldown_total_seconds", fightResponse.Cooldown.TotalSeconds))
		}

		if fightResponse.Characters[0].Level >= lvl {
			return nil
		}
	}
}

func (a *app) FightCharacter(ctx context.Context, character string) (*api.CharacterFightDataSchema, error) {
	// TODO: add participants
	rsp, err := a.di.ArtifactClient().ActionFightMyNameActionFightPost(ctx, character, api.ActionFightMyNameActionFightPostJSONBody{})
	if err != nil {
		return nil, err
	}

	resp, err := api.ParseActionFightMyNameActionFightPostResponse(rsp)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}

	return &resp.JSON200.Data, nil
}

func (a *app) MoveCharacter(ctx context.Context, character string, x, y int) (*api.CharacterMovementDataSchema, error) {
	rsp, err := a.di.ArtifactClient().ActionMoveMyNameActionMovePost(ctx, character, api.ActionMoveMyNameActionMovePostJSONRequestBody{
		X: new(x),
		Y: new(y),
	})
	if err != nil {
		return nil, err
	}

	resp, err := api.ParseActionMoveMyNameActionMovePostResponse(rsp)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}

	return &resp.JSON200.Data, nil
}

func (a *app) RestCharacter(ctx context.Context, character string) (*api.CharacterRestDataSchema, error) {
	rsp, err := a.di.ArtifactClient().ActionRestMyNameActionRestPost(ctx, character)
	if err != nil {
		return nil, err
	}

	resp, err := api.ParseActionRestMyNameActionRestPostResponse(rsp)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}

	return &resp.JSON200.Data, nil
}

func (a *app) initDeps() {
	a.initOnce.Do(func() {
		deps := []func() error{}

		for n, dep := range deps {
			if err := dep(); err != nil {
				a.di.Logger().Error("Failed to init dep", zap.Int("n", n))
				panic(err)
			}
		}
	})
}
