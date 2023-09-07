package docker

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
)

func (a *Adapter) readLogsOfDyingContainer(ctx context.Context, event *events.Message) {
	a.updateMutex.Lock()
	defer a.updateMutex.Unlock()

	containerID := event.Actor.ID
	containerName := event.Actor.Attributes["name"]
	if containerName == "" {
		log.Ctx(ctx).Error().Msgf("name of actor is empty, skipping reading logs of dying container %s", containerID)
		return
	}
	log.Ctx(ctx).Info().Msgf("reading logs of dying container [%s %s]", containerID, containerName)
	a.readContainerLogs(ctx, containerID, containerName)
}

func (a *Adapter) update(ctx context.Context) {
	a.updateMutex.Lock()
	defer a.updateMutex.Unlock()

	containers, err := a.cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("get containers error")
		return
	}

	result := make([]ContainerInfo, 0, len(containers))
	for _, cont := range containers {
		if cont.State == "running" {
			resultCont := ContainerInfo{
				ID:   cont.ID,
				Name: strings.TrimPrefix(cont.Names[0], "/"),
			}
			result = append(result, resultCont)
		}
	}

	//fmt.Printf("got containers: %v\n", result)

	for _, cont := range result {
		a.readContainerLogs(ctx, cont.ID, cont.Name)
	}
}

func (a *Adapter) readContainerLogs(ctx context.Context, containerID, containerName string) {
	//fmt.Println("reading container logs", containerName)
	since := a.LogAdapter.GetSinceTimestamp(ctx, containerName)

	reader, err := a.cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      since,
		Until:      "",
		Timestamps: true,
		Follow:     false,
		Tail:       "",
		Details:    false,
	})
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("container logs error")
		return
	}

	defer reader.Close()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		log.Ctx(ctx).Err(err).Msg("read logs buffer error")
		return
	}
	a.LogAdapter.WriteMessage(ctx, containerName, buf)
}
