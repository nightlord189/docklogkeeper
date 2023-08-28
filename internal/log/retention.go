package log

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

func (a *Adapter) ClearOldFiles(ctx context.Context) {
	beforeThan := time.Now().Add(-1 * time.Duration(a.Config.Retention) * time.Second).Unix()
	if err := a.deleteOldLogs(beforeThan); err != nil {
		log.Ctx(ctx).Err(err).Msg("delete old logs error")
	}
}
