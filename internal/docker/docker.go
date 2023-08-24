package docker

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/nightlord189/docklogkeeper/internal/config"
	"github.com/nightlord189/docklogkeeper/internal/log"
	"github.com/rs/zerolog"
	"sync"
)

type Adapter struct {
	Config      config.Config
	LogAdapter  *log.Adapter
	updateMutex *sync.Mutex
	cli         *client.Client
}

func New(ctx context.Context, cfg config.Config, lg *log.Adapter) (*Adapter, error) {
	cli, err := getCli(ctx)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		Config:      cfg,
		LogAdapter:  lg,
		updateMutex: &sync.Mutex{},
		cli:         cli,
	}, nil
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
