package docker

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/nightlord189/docklogkeeper/internal/config"
	"github.com/nightlord189/docklogkeeper/internal/log"
	"github.com/rs/zerolog"
	"time"
)

type Adapter struct {
	Config         config.Config
	LogAdapter     *log.Adapter
	cli            *client.Client
	lastTimestamps map[string]string
}

func New(ctx context.Context, cfg config.Config, lg *log.Adapter) (*Adapter, error) {
	cli, err := getCli(ctx)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		Config:         cfg,
		LogAdapter:     lg,
		cli:            cli,
		lastTimestamps: make(map[string]string, 10),
	}, nil
}

func (a *Adapter) Run(ctx context.Context) {
	regularTicker := time.NewTicker(10 * time.Second)

	a.update(ctx)
outerLoop:
	for {
		select {
		case <-regularTicker.C:
			a.update(ctx)
		case <-ctx.Done():
			break outerLoop
		}
	}
	zerolog.Ctx(ctx).Info().Msg("stopping docker adapter run")
}

func getCli(ctx context.Context) (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("create docker client error")
		return nil, err
	}
	return cli, nil
}

func (a *Adapter) Close() {
	a.cli.Close()
}
