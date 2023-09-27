package log

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

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
		zerolog.Ctx(ctx).Error().Msgf("parse timestamp error: splitted count less than required, log: %s", log)
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
	switch {
	case strings.Contains(containerName, "captain-") && strings.Contains(containerName, "."): // for CapRover
		splitted := strings.Split(containerName, ".")
		return splitted[0]
	default:
		return containerName
	}
}
