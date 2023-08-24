package log

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func getSize(ctx context.Context, fileWriter *os.File) int64 {
	fileInfo, err := fileWriter.Stat()
	if err != nil {
		log.Ctx(ctx).Err(err).Str("filename", fileWriter.Name()).Msg("get file info from fileWriter error")
		return 0
	}
	return fileInfo.Size()
}

func getNextFileName(ctx context.Context, lastFileName string) string {
	baseFileName := path.Base(lastFileName)
	splitted := strings.Split(baseFileName, ".")
	if len(splitted) != 2 {
		return defaultFileName
	}
	rawNumber := splitted[0]

	number, err := strconv.ParseInt(rawNumber, 10, 64)
	if err != nil {
		log.Ctx(ctx).Err(err).Str("file_name", lastFileName).Msg("parsing file name to int error")
		return defaultFileName
	}

	number++

	return fmt.Sprintf("%d.txt", number)
}

func getLastTimestampFromLog(log string) *time.Time {
	splitted := strings.Split(log, " ")
	if len(splitted) < 2 {
		return nil
	}
	timestamp, err := time.Parse(time.RFC3339, splitted[0])
	if err != nil {
		return nil
	}
	return &timestamp
}

func calcMergedContainerName(containerName string) string {
	if strings.Contains(containerName, ".") && !strings.HasPrefix(containerName, ".") {
		splitted := strings.Split(containerName, ".")
		return splitted[0]
	}
	return containerName
}
