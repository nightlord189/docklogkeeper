package docker

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func (a *Adapter) Run(ctx context.Context) {
	regularTicker := time.NewTicker(10 * time.Second)

	go a.listenEvents(ctx)

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

func (a *Adapter) listenEvents(ctx context.Context) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	eventsChan, errChan := a.cli.Events(ctx2, types.EventsOptions{
		Since:   time.Now().Format(time.RFC3339),
		Until:   "",
		Filters: filters.Args{},
	})

	log.Ctx(ctx).Debug().Msg("listening events")
outerLoop:
	for {
		select {
		case event, ok := <-eventsChan:
			if !ok {
				return
			}
			if event.Type != "container" {
				continue
			}
			log.Ctx(ctx).Info().Interface("event", event).Msg("new docker event")
			a.update(ctx)
		case err, ok := <-errChan:
			if !ok {
				return
			}
			log.Ctx(ctx).Err(err).Msg("error from docker events channel")
			break outerLoop
		}
	}

	log.Ctx(ctx).Debug().Msg("restart listen events")
	a.listenEvents(ctx)
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
		ShowStderr: false,
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
