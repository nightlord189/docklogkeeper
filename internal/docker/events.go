package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/rs/zerolog/log"
	"time"
)

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
			if event.Action == "die" {
				log.Ctx(ctx).Info().Msgf("new docker event: %s %s", event.Action, event.Actor.Attributes["name"])
				a.readLogsOfDyingContainer(ctx, &event)
			}
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
