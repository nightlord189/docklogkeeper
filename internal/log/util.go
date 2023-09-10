package log

import (
	"context"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

func reverseLines(lines []string) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i] //reverse the slice
	}
}

// trim Docker log message header (8 bytes)
// see https://ahmet.im/blog/docker-logs-api-binary-format-explained/
func trimMessageHeader(line []byte) string {
	if line[0] <= 2 { // has prefix
		return string(line[8:])
	}
	// doesn't have prefix
	return string(line)
}

func timeGreaterOrEqualNil(t1, t2 *time.Time) bool {
	return t1 != nil && t2 != nil && t1.UnixMicro() >= t2.UnixMicro()
}

func getTimestampFromLog(ctx context.Context, log string) *time.Time {
	splitted := strings.Split(log, " ")
	if len(splitted) < 2 {
		zerolog.Ctx(ctx).Error().Msgf("parse timestamp error: splitted count less than required, log:", log)
		return nil
	}
	timestamp, err := time.Parse(time.RFC3339, splitted[0])
	if err != nil {
		zerolog.Ctx(ctx).Error().Msgf("parse timestamp error: %v log: %s", err, log)
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
