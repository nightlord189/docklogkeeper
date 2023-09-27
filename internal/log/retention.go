package log

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (a *Adapter) ClearOldLogs(ctx context.Context) {
	beforeThan := time.Now().Add(-1 * time.Duration(a.Config.Retention) * time.Second).Unix()
	if err := a.Repo.DeleteOldLogs(beforeThan); err != nil {
		log.Ctx(ctx).Err(err).Msg("delete old logs error")
	}
}
