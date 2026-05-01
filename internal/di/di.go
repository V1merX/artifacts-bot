package di

import (
	"context"
	"net/http"

	"github.com/V1merX/artifacts-bot/configs"
	"github.com/V1merX/artifacts-bot/pkg/api"
	"go.uber.org/zap"
)

type Container struct {
	config *configs.Config
	logger *zap.Logger

	artifactClient *api.Client
}

func NewContainer() Container {
	return Container{}
}

func (di *Container) Logger() *zap.Logger {
	if di.logger == nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		di.logger = logger
	}

	return di.logger
}

func (di *Container) Config() *configs.Config {
	if di.config == nil {
		config, err := configs.Load()
		if err != nil {
			di.Logger().Error("Failed to parse config", zap.Error(err))
		}

		di.config = config
	}

	return di.config
}

func (di *Container) ArtifactClient() *api.Client {
	if di.artifactClient == nil {
		client, err := api.NewClient(di.Config().ServerAddr, api.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set("Accept", "application/json")
			req.SetBasicAuth(di.Config().HTTPBasic.Username, di.Config().HTTPBasic.Password)
			return nil
		}))
		if err != nil {
			di.Logger().Error("Failed to create new artifact client", zap.Error(err))
		}

		di.artifactClient = client
	}

	return di.artifactClient
}
