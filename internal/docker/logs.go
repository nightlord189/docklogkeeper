package docker

import (
	"bufio"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
				Name: trimContainerName(cont.Names),
			}
			result = append(result, resultCont)
		}
	}

	// fmt.Printf("got containers: %v\n", result)

	log.Ctx(ctx).Info().Interface("containers", result).Msgf("start reading %d already existing containers", len(result))

	for _, cont := range result {
		go a.ensureReadContainerLogs(ctx, cont.ID, cont.Name)
	}
}

func (a *Adapter) ensureReadContainerLogs(ctx context.Context, containerID, containerName string) {
	a.readingMutex.RLock()
	_, ok := a.containersReading[containerID]
	a.readingMutex.RUnlock()
	if ok {
		return
	}
	a.readContainerLogs(ctx, containerID, containerName)
}

// running in a separate goroutine
func (a *Adapter) readContainerLogs(ctx context.Context, containerID, containerName string) {
	// fmt.Println("reading container logs", containerName)
	log.Ctx(ctx).Info().Msgf("start reading logs of container [%s, %s]", containerID, containerName)

	a.readingMutex.Lock()
	a.containersReading[containerID] = true
	a.readingMutex.Unlock()

	defer func() {
		a.readingMutex.Lock()
		delete(a.containersReading, containerID)
		a.readingMutex.Unlock()
	}()

	logger := log.Ctx(ctx).With().Str("container_id", containerID).Str("container_name", containerName).Logger()
	since := a.LogAdapter.GetSinceTimestamp(ctx, containerName)

	reader, err := a.cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      since,
		Until:      "",
		Timestamps: true,
		Follow:     true,
		Tail:       "",
		Details:    false,
	})
	if err != nil {
		logger.Err(err).Msgf("read container logs error [container %s]", containerName)
		return
	}

	defer func() {
		if err := reader.Close(); err != nil {
			logger.Err(err).Msg("close logs reader error")
		}
	}()

	buf := bufio.NewScanner(reader)

	for buf.Scan() {
		a.LogAdapter.WriteLine(ctx, containerName, buf.Bytes())
	}
	logger.Debug().Msg("stop reading logs")
}
