package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

func (a *Adapter) GetSinceTimestamp(containerName string) string {
	shortName := a.getMergedContainerName(containerName)
	since := a.lastTimestamps[shortName]
	if since == nil {
		// TODO: try get from files
		newSince := time.Now().Add(-5 * time.Minute)
		since = &newSince
		a.lastTimestamps[shortName] = since
	}
	return since.Format(time.RFC3339)
}

func (a *Adapter) WriteMessage(ctx context.Context, containerName string, buf *bytes.Buffer) {
	if buf.Len() == 0 {
		return
	}

	shortName := a.getMergedContainerName(containerName)

	fmt.Println("writeMessage", shortName, buf.Len())

	fileWriter := a.getFileWriter(ctx, shortName)
	if fileWriter == nil {
		return
	}

	lastLine := ""
	linesCount := 0
	for {
		readBytes, err := buf.ReadBytes('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Ctx(ctx).Err(err).Msg("read bytes error")
			}
			break
		}
		if len(readBytes) > 12 {
			lastLine = string(readBytes[8:])
		}
		fileWriter.Write(readBytes[8:]) // strip header from docker
		linesCount++
	}
	fmt.Printf("lines count: %d, lastLine: %s\n", linesCount, lastLine)

	timestampFromLog := getLastTimestampFromLog(lastLine)
	if timestampFromLog != nil {
		newTimestamp := timestampFromLog.Add(1 * time.Second)
		a.lastTimestamps[shortName] = &newTimestamp
		fmt.Println("new timestamp", timestampFromLog, newTimestamp)
	}
}

func (a *Adapter) getFileWriter(ctx context.Context, shortName string) *os.File {
	fileWriter, exists := a.currentFiles[shortName]
	if !exists {
		ensureDir(path.Join(a.Config.Dir, shortName))
		fileName := fmt.Sprintf("%s.txt", shortName)
		fullFileName := path.Join(a.Config.Dir, shortName, fileName)
		newFile, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Ctx(ctx).Err(err).Str("file_name", fullFileName).Msg("open new file error")
			return nil
		}
		fileWriter = newFile
		a.currentFiles[shortName] = fileWriter
	}
	return fileWriter
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

func (a *Adapter) getMergedContainerName(containerName string) string {
	result := a.names[containerName]
	if result == "" {
		result = calcMergedContainerName(containerName)
		a.names[containerName] = result
	}
	return result
}

func calcMergedContainerName(containerName string) string {
	if strings.Contains(containerName, ".") && !strings.HasPrefix(containerName, ".") {
		splitted := strings.Split(containerName, ".")
		return splitted[0]
	}
	return containerName
}
