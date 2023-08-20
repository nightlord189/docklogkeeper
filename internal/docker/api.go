package docker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func (a *Adapter) update(ctx context.Context) {
	fmt.Println("update")

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

	since := a.lastTimestamps[containerName]
	if since == "" {
		since = time.Now().Add(-5 * time.Minute).Format(time.RFC3339)
	}
	reader, err := a.cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: false,
		Since:      since,
		Until:      "",
		Timestamps: false,
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
	a.lastTimestamps[containerName] = time.Now().Format(time.RFC3339)
}
