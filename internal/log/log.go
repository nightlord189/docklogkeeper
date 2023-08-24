package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"io/fs"
	"os"
	"path"
	"time"
)

const defaultFileName = "1.txt"

func (a *Adapter) WriteMessage(ctx context.Context, containerName string, buf *bytes.Buffer) {
	if buf.Len() == 0 {
		return
	}

	shortName := a.getMergedContainerName(containerName)

	fmt.Println("writeMessage", shortName, buf.Len())

	ctx = log.Ctx(ctx).With().Str("short_name", shortName).Logger().WithContext(ctx)

	fileData := a.getFileData(ctx, shortName)
	if fileData == nil {
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
		n, _ := fileData.Writer.Write(readBytes[8:]) // strip header from docker
		fileData.Size += int64(n)
		linesCount++
	}
	fmt.Printf("lines count: %d, lastLine: %s\n", linesCount, lastLine)

	a.currentFiles[shortName] = fileData

	a.checkCurrentChunkSize(ctx, shortName)

	timestampFromLog := getLastTimestampFromLog(lastLine)
	if timestampFromLog != nil {
		newTimestamp := timestampFromLog.Add(1 * time.Second)
		a.lastTimestamps[shortName] = &newTimestamp
		fmt.Println("new timestamp", timestampFromLog, newTimestamp)
	}
}

func (a *Adapter) getFileData(ctx context.Context, shortName string) *FileData {
	fileData, exists := a.currentFiles[shortName]
	if !exists {
		ensureDir(path.Join(a.Config.Dir, shortName))
		fileWriter := a.getLastChunkFromDir(ctx, shortName)
		if fileWriter == nil {
			return nil
		}
		fileData = &FileData{
			Writer: fileWriter,
			Size:   getSize(ctx, fileWriter),
		}
		a.currentFiles[shortName] = fileData
	}
	return fileData
}

func (a *Adapter) checkCurrentChunkSize(ctx context.Context, shortName string) {
	fileData := a.currentFiles[shortName]
	if fileData.Size >= a.Config.ChunkSize {
		log.Ctx(ctx).Info().Int64("current_size", fileData.Size).
			Str("current_chunk", fileData.Writer.Name()).Msg("current chunk size is too big, creating new chunk")
		// close file
		if err := fileData.Writer.Close(); err != nil {
			log.Ctx(ctx).Err(err).Str("filename", fileData.Writer.Name()).Msg("close current chunk error")
		}
		// increase number
		nextChunkName := getNextFileName(ctx, fileData.Writer.Name())
		// open new file
		newWriter := openFile(ctx, path.Join(a.Config.Dir, shortName, nextChunkName), false)
		log.Ctx(ctx).Info().Int64("current_size", fileData.Size).Str("new_chunk", nextChunkName).Msg("new chunk")
		if newWriter == nil {
			return
		}
		a.currentFiles[shortName] = &FileData{
			Writer: newWriter,
			Size:   getSize(ctx, newWriter),
		}
	}
}

func (a *Adapter) getLastChunkFromDir(ctx context.Context, shortName string) *os.File {
	fileName := defaultFileName

	lastFileInfo := a.getLastFileNameFromDir(ctx, shortName)
	if lastFileInfo != nil {
		if lastFileInfo.Size() >= a.Config.ChunkSize {
			fileName = getNextFileName(ctx, lastFileInfo.Name())
			log.Ctx(ctx).Info().Str("old_filename", lastFileInfo.Name()).Msgf("last chunk size is too big, creating next chunk %s", fileName)
		} else {
			fileName = lastFileInfo.Name()
			log.Ctx(ctx).Info().Str("filename", fileName).Msgf("opening last existing chunk")
		}
	} else {
		log.Ctx(ctx).Info().Str("filename", fileName).Msgf("opening first chunk")
	}

	fileWriter := openFile(ctx, path.Join(a.Config.Dir, shortName, fileName), false)

	return fileWriter
}

func (a *Adapter) getLastFileNameFromDir(ctx context.Context, shortName string) fs.FileInfo {
	fSortedFiles := a.getSortedFilesByDir(shortName)

	if len(fSortedFiles) == 0 {
		return nil
	}

	lastFile := fSortedFiles[len(fSortedFiles)-1]

	fileInfo, err := lastFile.Info()
	if err != nil {
		log.Ctx(ctx).Err(err).Str("file_name", lastFile.Name()).Msg("get file info error")
	}

	return fileInfo
}

func (a *Adapter) getSortedFilesByDir(shortName string) []os.DirEntry {
	dir := path.Join(a.Config.Dir, shortName)
	if _, err := os.Stat(dir); err != nil {
		return nil
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	if len(files) == 0 {
		return nil
	}

	fSortedFiles := getFilteredAndSortedFiles(files)

	return fSortedFiles
}

func (a *Adapter) getMergedContainerName(containerName string) string {
	result := a.names[containerName]
	if result == "" {
		result = calcMergedContainerName(containerName)
		a.names[containerName] = result
	}
	return result
}
