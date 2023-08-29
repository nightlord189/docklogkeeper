package log

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

func (a *Adapter) GetSinceTimestamp(ctx context.Context, containerName string) string {
	shortName := a.GetShortContainerName(containerName)
	since := a.lastTimestamps[shortName]
	if since == nil {
		ctx = log.Ctx(ctx).With().Str("short_name", shortName).Logger().WithContext(ctx)
		since = a.Repo.GetLastTimestamp(shortName)
		if since == nil {
			newSince := time.Now().Add(-5 * time.Minute)
			since = &newSince
		} else {
			log.Ctx(ctx).Info().Str("short_name", shortName).Time("timestamp", *since).Msg("last timestamp loaded from file")
		}
		a.lastTimestamps[shortName] = since
	}
	return since.Format(time.RFC3339)
}
