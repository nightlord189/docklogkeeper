package log

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

func openFile(ctx context.Context, filePath string, readOnly bool) *os.File {
	flag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	if readOnly {
		flag = os.O_RDONLY
	}
	fileWriter, err := os.OpenFile(filePath, flag, 0644)
	if err != nil {
		log.Ctx(ctx).Err(err).Str("file_name", filePath).Msg("open file error")
		return nil
	}
	return fileWriter
}

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

func getFilteredAndSortedFiles(files []os.DirEntry) []os.DirEntry {
	filtered := make([]os.DirEntry, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			filtered = append(filtered, file)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name() < filtered[j].Name()
	})

	return filtered
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

func calcShortContainerName(containerName string) string {
	if strings.Contains(containerName, ".") && !strings.HasPrefix(containerName, ".") {
		splitted := strings.Split(containerName, ".")
		return splitted[0]
	}
	return containerName
}
