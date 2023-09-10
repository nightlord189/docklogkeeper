package log

import (
	"context"
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/rs/zerolog/log"
	"time"
)

func (a *Adapter) WriteLine(ctx context.Context, containerName string, line []byte) {
	shortName := a.GetShortContainerName(containerName)

	//fmt.Println("WriteLine", shortName, string(line))

	ctx = log.Ctx(ctx).With().Str("short_name", shortName).Logger().WithContext(ctx)

	lastTimestamp := a.lastTimestamps[shortName]

	if len(line) < 8 {
		return
	}

	lineStr := trimMessageHeader(line)

	timestampFromLog := getTimestampFromLog(ctx, lineStr)
	if timeGreaterOrEqualNil(lastTimestamp, timestampFromLog) { // check also for equal
		//log.Ctx(ctx).Debug().Msgf("skipping line because timestamp, last_tt: %v, current_tt: %v", lastTimestamp, ttFromLog)
		return
	}

	var createdAt time.Time
	if timestampFromLog != nil {
		createdAt = *timestampFromLog
	} else {
		createdAt = time.Now()
	}

	logToInsert := entity.LogDataDB{
		ContainerName: shortName,
		LogText:       lineStr,
		CreatedAt:     createdAt.Unix(),
	}

	if err := a.Repo.EnsureContainer(shortName); err != nil {
		log.Ctx(ctx).Err(err).Msg("ensure container in db error")
	}
	if err := a.Repo.InsertLog(&logToInsert); err != nil {
		log.Ctx(ctx).Err(err).Msg("insert log error")
	}

	if timestampFromLog != nil {
		a.lastTimestamps[shortName] = timestampFromLog
		//fmt.Println(containerName, "last timestamp", timestampFromLog)
	}
}

func (a *Adapter) GetShortContainerName(containerName string) string {
	result := a.names[containerName]
	if result == "" {
		result = a.Repo.GetMappedName(containerName)
		if result == "" {
			result = calcShortContainerName(containerName)
			a.Repo.SetMappedName(containerName, result)
		}
		a.names[containerName] = result
	}
	return result
}
