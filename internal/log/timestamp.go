package log

import (
	"context"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"strings"
	"time"
)

func (a *Adapter) GetSinceTimestamp(ctx context.Context, containerName string) string {
	shortName := a.GetShortContainerName(containerName)
	since := a.lastTimestamps[shortName]
	if since == nil {
		ctx = log.Ctx(ctx).With().Str("short_name", shortName).Logger().WithContext(ctx)
		since = a.getLastTimestampFromFiles(ctx, shortName)
		if since == nil {
			newSince := time.Now().Add(-5 * time.Minute)
			since = &newSince
		} else {
			log.Ctx(ctx).Info().Str("short_name", shortName).Time("timestamp", *since).Msg("last timestamp loaded from file")
			since2 := since.Add(1 * time.Second)
			since = &since2
		}
		a.lastTimestamps[shortName] = since
	}
	return since.Format(time.RFC3339)
}

func (a *Adapter) getLastTimestampFromFiles(ctx context.Context, shortName string) *time.Time {
	files := a.getSortedFilesByDir(shortName)
	if len(files) == 0 {
		return nil
	}
	// check last two files to get timestamp
	for i := len(files) - 1; i >= len(files)-2 && i >= 0; i-- {
		filePath := path.Join(a.Config.Dir, shortName, files[i].Name())
		timestamp := getTimestampFromFile(ctx, filePath)
		if timestamp != nil {
			return timestamp
		}
	}
	return nil
}

const readEndSize = 1000

func getTimestampFromFile(ctx context.Context, filePath string) *time.Time {
	stat, err := os.Stat(filePath)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("get stat from file error")
		return nil
	}

	stat.Size()

	file := openFile(ctx, filePath, true)
	if file == nil {
		return nil
	}
	defer file.Close()

	buf := make([]byte, readEndSize)
	start := stat.Size() - readEndSize
	if start < 0 {
		start = 0
	}
	_, err = file.ReadAt(buf, start)
	if err != nil {
		log.Ctx(ctx).Err(err).Str("filename", filePath).Msg("read end of file error")
		return nil
	}

	splittedLines := strings.Split(string(buf), "\n")
	if len(splittedLines) == 0 {
		return nil
	}
	for i := len(splittedLines) - 1; i >= 0; i-- {
		timestamp := getLastTimestampFromLog(splittedLines[i])
		if timestamp != nil {
			log.Ctx(ctx).Info().Str("filename", filePath).
				Str("line", splittedLines[i]).Time("timestamp", *timestamp).Msgf("found timestamp in line %d", i)
			return timestamp
		}
	}
	return nil
}
