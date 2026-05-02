package di

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/V1merX/artifacts-bot/configs"
	artifactsGW "github.com/V1merX/artifacts-bot/internal/gateway/artifacts"
	"github.com/V1merX/artifacts-bot/internal/usecase/levelup"
	"github.com/V1merX/artifacts-bot/pkg/api"
	"go.uber.org/zap"
)

type Container struct {
	config *configs.Config
	logger *zap.Logger
	auth   *authHolder

	artifactClient *api.Client
	gateway        *artifactsGW.Gateway
	levelUpUseCase *levelup.UseCase
}

func NewContainer() Container {
	return Container{}
}

func (c *Container) Logger() *zap.Logger {
	if c.logger == nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		c.logger = logger
	}
	return c.logger
}

func (c *Container) Config() *configs.Config {
	if c.config == nil {
		config, err := configs.Load()
		if err != nil {
			c.Logger().Error("Failed to parse config", zap.Error(err))
		}
		c.config = config
	}
	return c.config
}

func (c *Container) Authenticate(ctx context.Context) error {
	resp, err := c.artifactAPIClient().GenerateTokenTokenPost(ctx)
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp api.TokenResponseSchema
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("decode token response: %w", err)
	}

	c.authHolder().SetBearerToken(tokenResp.Token)
	c.Logger().Info("Authenticated", zap.String("token", fmt.Sprintf("%s...", tokenResp.Token[:55])))
	return nil
}

func (c *Container) ArtifactsGateway() levelup.Gateway {
	if c.gateway == nil {
		c.gateway = artifactsGW.New(c.artifactAPIClient())
	}
	return c.gateway
}

func (c *Container) LevelUpUseCase() *levelup.UseCase {
	if c.levelUpUseCase == nil {
		c.levelUpUseCase = levelup.New(c.ArtifactsGateway(), c.Logger())
	}
	return c.levelUpUseCase
}

func (c *Container) authHolder() *authHolder {
	if c.auth == nil {
		cfg := c.Config()
		c.auth = newAuthHolder(cfg.HTTPBasic.Username, cfg.HTTPBasic.Password)
	}
	return c.auth
}

func (c *Container) artifactAPIClient() *api.Client {
	if c.artifactClient == nil {
		client, err := api.NewClient(c.Config().ServerAddr,
			api.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
				req.Header.Set("Accept", "application/json")
				return nil
			}),
			api.WithRequestEditorFn(c.authHolder().Editor),
		)
		if err != nil {
			c.Logger().Error("Failed to create artifact client", zap.Error(err))
		}
		c.artifactClient = client
	}
	return c.artifactClient
}
