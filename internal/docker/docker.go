package docker

import (
	"context"
	"sync"

	"github.com/docker/docker/client"
	"github.com/nightlord189/docklogkeeper/internal/config"
	"github.com/nightlord189/docklogkeeper/internal/log"
	"github.com/rs/zerolog"
)

type Adapter struct {
	Config            config.Config
	LogAdapter        *log.Adapter
	updateMutex       *sync.Mutex
	cli               *client.Client
	containersReading map[string]bool // container_id -> true
	readingMutex      *sync.RWMutex
}

func New(ctx context.Context, cfg config.Config, lg *log.Adapter) (*Adapter, error) {
	cli, err := getCli(ctx)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		Config:            cfg,
		LogAdapter:        lg,
		updateMutex:       &sync.Mutex{},
		containersReading: make(map[string]bool, 10),
		readingMutex:      &sync.RWMutex{},
		cli:               cli,
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

func (a *Adapter) Close(ctx context.Context) {
	if err := a.cli.Close(); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("close docker adapter error")
	}
}
