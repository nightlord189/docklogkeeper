package log

import (
	"context"
	"github.com/rs/zerolog/log"
)

func (a *Adapter) SearchLines(ctx context.Context, shortName string, req SearchRequest) []string {
	logs, err := a.Repo.SearchLogs(shortName, req.Contains)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("search logs error")
		return []string{}
	}
	lines := make([]string, len(logs))
	for i, logEntry := range logs {
		lines[i] = logEntry.LogText
	}
	return lines
}
