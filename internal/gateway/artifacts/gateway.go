package artifacts

import (
	"context"
	"fmt"

	"github.com/V1merX/artifacts-bot/internal/domain/action"
	"github.com/V1merX/artifacts-bot/internal/domain/character"
	"github.com/V1merX/artifacts-bot/pkg/api"
)

type Gateway struct {
	client *api.Client
}

func New(client *api.Client) *Gateway {
	return &Gateway{client: client}
}

func (g *Gateway) GetCharacter(ctx context.Context, name string) (*character.Character, error) {
	rsp, err := g.client.GetCharacterCharactersNameGet(ctx, name)
	if err != nil {
		return nil, err
	}

	resp, err := api.ParseGetCharacterCharactersNameGetResponse(rsp)
	if err != nil {
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}

	status, err := character.NewStatus(character.MoveStatusType)
	if err != nil {
		return nil, err
	}

	return character.NewCharacter(
		resp.JSON200.Data.Name,
		resp.JSON200.Data.Level,
		resp.JSON200.Data.Xp,
		resp.JSON200.Data.MaxXp,
		resp.JSON200.Data.Hp,
		resp.JSON200.Data.MaxHp,
		status,
	)
}

func (g *Gateway) Fight(ctx context.Context, characterName string) (*action.FightResult, error) {
	rsp, err := g.client.ActionFightMyNameActionFightPost(ctx, characterName, api.ActionFightMyNameActionFightPostJSONBody{})
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

	data := resp.JSON200.Data
	characters := make([]action.CharacterSnapshot, len(data.Characters))
	for i, c := range data.Characters {
		characters[i] = action.CharacterSnapshot{
			HP:    c.Hp,
			MaxHP: c.MaxHp,
			Level: c.Level,
		}
	}

	return &action.FightResult{
		Cooldown:   action.Cooldown{TotalSeconds: data.Cooldown.TotalSeconds},
		Characters: characters,
	}, nil
}

func (g *Gateway) Move(ctx context.Context, characterName string, x, y int) (*action.MoveResult, error) {
	rsp, err := g.client.ActionMoveMyNameActionMovePost(ctx, characterName, api.ActionMoveMyNameActionMovePostJSONRequestBody{
		X: &x,
		Y: &y,
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

	data := resp.JSON200.Data
	return &action.MoveResult{
		Cooldown: action.Cooldown{TotalSeconds: data.Cooldown.TotalSeconds},
		Character: action.CharacterSnapshot{
			HP:    data.Character.Hp,
			MaxHP: data.Character.MaxHp,
			Level: data.Character.Level,
		},
	}, nil
}

func (g *Gateway) Rest(ctx context.Context, characterName string) (*action.RestResult, error) {
	rsp, err := g.client.ActionRestMyNameActionRestPost(ctx, characterName)
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

	return &action.RestResult{
		Cooldown: action.Cooldown{TotalSeconds: resp.JSON200.Data.Cooldown.TotalSeconds},
	}, nil
}
