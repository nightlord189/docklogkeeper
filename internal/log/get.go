package log

import (
	"context"
	"fmt"
	"github.com/icza/backscanner"
	"github.com/rs/zerolog/log"
	"path"
	"strings"
)

// GetLogs - returns log of container with cursor pagination
// cursor -> file_name-bytes_position
// 1-3983 means 1.txt, 3983 byte
// so we will read from 3983 position back to 0
func (a *Adapter) GetLogs(ctx context.Context, req GetLinesRequest) (GetLinesResponse, error) {
	fileEntries := a.getSortedFilesByDir(req.ShortName)
	if len(fileEntries) == 0 {
		return GetLinesResponse{}, fmt.Errorf("there is no log files for this container")
	}

	lines := make([]string, 0, req.Limit)

	currentFileNumber := req.FileNumber
	currentOffset := req.Offset

	for i := len(fileEntries) - 1; i >= 0; i-- {
		currentEntry := fileEntries[i]
		if strings.TrimSuffix(currentEntry.Name(), ".txt") != fmt.Sprintf("%d", currentFileNumber) {
			continue
		}
		currentFile := openFile(ctx, path.Join(a.Config.Dir, req.ShortName, currentEntry.Name()), true)
		if currentFile == nil {
			return GetLinesResponse{}, fmt.Errorf("open file error")
		}

		if currentOffset == 0 { //new file
			currentFileInfo, err := currentEntry.Info()
			if err != nil {
				log.Ctx(ctx).Err(err).Str("filename", currentEntry.Name()).Msg("get file info error")
				return GetLinesResponse{}, fmt.Errorf("get file info error")
			}
			currentOffset = int(currentFileInfo.Size())
		}

		scanner := backscanner.New(currentFile, currentOffset)
		for {
			line, pos, err := scanner.Line()
			if err != nil {
				fmt.Println("error read line", err)
				break
			}
			currentOffset = pos
			lines = append(lines, line)
			if len(lines) >= req.Limit {
				break
			}
			//fmt.Printf("Line position: %2d, line: %q\n", pos, line)
		}
		if err := currentFile.Close(); err != nil {
			log.Ctx(ctx).Err(err).Msg("close file error")
		}
		if len(lines) >= req.Limit {
			break
		}
		currentFileNumber--
		if currentFileNumber < 0 {
			break
		}
		currentOffset = 0
	}

	return GetLinesResponse{
		Lines:      lines,
		FileNumber: currentFileNumber,
		Offset:     currentOffset,
	}, nil
}