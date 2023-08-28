package log

import (
	"bufio"
	"context"
	"fmt"
	"github.com/icza/backscanner"
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

func (a *Adapter) GetLogs(ctx context.Context, req entity.GetLogsRequest) (GetLogsResponse, error) {
	fileEntries := a.getSortedFilesByDir(req.ShortName)
	if len(fileEntries) == 0 {
		return GetLogsResponse{Records: []string{}}, nil
	}

	currentChunk := req.ChunkNumber

	if currentChunk == 0 {
		currentChunk = getChunkNumberFromFileName(fileEntries[len(fileEntries)-1].Name())
	}

	chunkIndex := getIndexOfChunk(currentChunk, fileEntries)
	if chunkIndex < 0 {
		return GetLogsResponse{}, fmt.Errorf("index of the first chunk is invalid")
	}

	pos := req.Position
	var (
		eof bool
		err error
	)

	lines := make([]string, 0, 100)

	switch {
	case req.Direction == entity.DirFuture && req.ChunkNumber == 0 && req.Position == 0:
		for i := chunkIndex; i >= 0; i-- {
			currentChunk = getChunkNumberFromFileName(fileEntries[i].Name())
			pos, eof, err = readLogFileBackward(ctx, path.Join(a.Config.Dir, req.ShortName, fileEntries[i].Name()), pos, lines, req.Limit)
			if eof || err != nil {
				break
			}
		}
	case req.Direction == entity.DirFuture && req.ChunkNumber > 0:
		panic("read from specific chunk from pos to end")
	case req.Direction == entity.DirPast && req.ChunkNumber > 0:
		panic("read from specific chunk from pos to start")
	default:
		return GetLogsResponse{}, fmt.Errorf("invalid request")
	}

	return GetLogsResponse{
		Records:     lines,
		ChunkNumber: currentChunk,
		Offset:      pos,
	}, err
}

func

// from start (or initialPos) to end
// returns pos, isEOF, err
func readLogFileForward(ctx context.Context, filePath string, initialPos int, lines []string, limit int) (int, bool, error) {
	currentFile := openFile(ctx, filePath, true)
	if currentFile == nil {
		return 0, false, fmt.Errorf("open file error")
	}

	defer func() {
		if err := currentFile.Close(); err != nil {
			log.Ctx(ctx).Err(err).Msg("close file error")
		}
	}()

	if _, err := currentFile.Seek(int64(initialPos), 0); err != nil {
		log.Ctx(ctx).Err(err).Str("filename", filePath).Msg("seek error")
		return 0, false, fmt.Errorf("seek error: %w", err)
	}

	r := bufio.NewReader(currentFile)
	pos := int64(initialPos)
	for {
		data, err := r.ReadString('\n')
		pos += int64(len(data))
		if err == nil || err == io.EOF {
			if len(data) > 0 && data[len(data)-1] == '\n' {
				data = data[:len(data)-1]
			}
			if len(data) > 0 && data[len(data)-1] == '\r' {
				data = data[:len(data)-1]
			}
			lines = append(lines, data)
			if len(lines) >= limit {
				break
			}
			//fmt.Printf("Pos: %d, Read: %s\n", pos, data)
		}
		if err != nil {
			if err != io.EOF {
				log.Ctx(ctx).Err(err).Str("filename", filePath).Msg("read error")
			}
			return int(pos), true, nil
		}
	}

	return int(pos), false, nil
}

// from end (or initialPos) to start
func readLogFileBackward(ctx context.Context, filePath string, initialPos int, lines []string, limit int) (int, bool, error) {
	currentFile := openFile(ctx, filePath, true)
	if currentFile == nil {
		return 0, false, fmt.Errorf("open file error")
	}

	defer func() {
		if err := currentFile.Close(); err != nil {
			log.Ctx(ctx).Err(err).Msg("close file error")
		}
	}()

	if initialPos == 0 { //new file
		currentFileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Ctx(ctx).Err(err).Str("filename", filePath).Msg("get file info error")
			return 0, false, fmt.Errorf("get file info error")
		}
		initialPos = int(currentFileInfo.Size())
	}

	scanner := backscanner.New(currentFile, initialPos)

	pos := initialPos
	for {
		line, currentPos, err := scanner.Line()
		if err != nil {
			fmt.Println("error read line", err)
			return pos, true, nil
		}
		pos = currentPos
		lines = append(lines, line)
		if len(lines) >= limit {
			break
		}
		//fmt.Printf("Line position: %2d, line: %q\n", pos, line)
	}
	return pos, false, nil
}

func getIndexOfChunk(chunkNumber int, chunks []os.DirEntry) int {
	chunkNumberStr := fmt.Sprintf("%d", chunkNumber)
	for i, chunk := range chunks {
		if strings.TrimSuffix(chunk.Name(), ".txt") == chunkNumberStr {
			return i
		}
	}
	return -1
}
