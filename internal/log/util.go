package log

import (
	"strconv"
	"strings"
	"time"
)

func getChunkNumberFromFileName(fileName string) int {
	parsedChunkNumber, err := strconv.Atoi(strings.TrimSuffix(fileName, ".txt"))
	if err != nil {
		return -1
	}
	return parsedChunkNumber
}

func reverseLines(lines []string) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i] //reverse the slice
	}
}

func timeGreaterOrEqualNil(t1, t2 *time.Time) bool {
	return t1 != nil && t2 != nil && t1.UnixMicro() >= t2.UnixMicro()
}

func getTimestampFromLog(log string) *time.Time {
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
