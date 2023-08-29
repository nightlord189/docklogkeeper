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

const defaultFileName = "1.txt"

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
		ttFromLog := getTimestampFromLog(string(readBytes[8:]))
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

	//fmt.Printf("lines count: %d, lastLine: %s\n", len(logs), logs[len(logs)-1].LogText)

	timestampFromLog := getTimestampFromLog(logs[len(logs)-1].LogText)
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
