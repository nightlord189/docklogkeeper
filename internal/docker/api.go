package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/rs/zerolog"
	"strings"
)

func (a *Adapter) Run(ctx context.Context) {
	a.update(ctx)

	a.listenEvents(ctx)
}

func (a *Adapter) GetAliveContainers(ctx context.Context) ([]string, error) {
	containers, err := a.cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("get containers error")
		return nil, fmt.Errorf("get containers error: %w", err)
	}
	result := make([]string, 0, len(containers))
	for _, cont := range containers {
		if cont.State == "running" {
			result = append(result, strings.TrimPrefix(cont.Names[0], "/"))
		}
	}
	return result, nil
}
