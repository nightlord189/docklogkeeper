package log

import (
	"bytes"
	"context"
	"errors"
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/rs/zerolog/log"
	"io"
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

func (a *Adapter) WriteMessage(ctx context.Context, containerName string, buf *bytes.Buffer) {
	if buf.Len() == 0 {
		return
	}

	shortName := a.GetShortContainerName(containerName)

	//fmt.Println("writeMessage", shortName, buf.Len())

	ctx = log.Ctx(ctx).With().Str("short_name", shortName).Logger().WithContext(ctx)

	lastTimestamp := a.lastTimestamps[shortName]
	foundNewLogs := false

	logs := make([]entity.LogDataDB, 0, 10)

	for {
		readBytes, err := buf.ReadBytes('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Ctx(ctx).Err(err).Msg("read bytes error")
			}
			break
		}
		if len(readBytes) < 8 {
			continue
		}
		ttFromLog := getTimestampFromLog(ctx, string(readBytes[8:]))
		if !foundNewLogs && timeGreaterOrEqualNil(lastTimestamp, ttFromLog) { // check also for equal
			//log.Ctx(ctx).Debug().Msgf("skipping line because timestamp, last_tt: %v, current_tt: %v", lastTimestamp, ttFromLog)
			continue
		} else {
			foundNewLogs = true
		}

		var createdAt time.Time
		if ttFromLog != nil {
			createdAt = *ttFromLog
		} else {
			createdAt = time.Now()
		}

		logs = append(logs, entity.LogDataDB{
			ContainerName: shortName,
			LogText:       string(readBytes[8:]),
			CreatedAt:     createdAt.Unix(),
		})
	}

	if len(logs) == 0 {
		return
	}

	if err := a.Repo.EnsureContainer(shortName); err != nil {
		log.Ctx(ctx).Err(err).Msg("ensure container error")
	}
	if err := a.Repo.InsertLogs(logs); err != nil {
		log.Ctx(ctx).Err(err).Msg("insert logs error")
	}

	timestampFromLog := getTimestampFromLog(ctx, logs[len(logs)-1].LogText)
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
