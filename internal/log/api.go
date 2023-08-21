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
)

func (a *Adapter) WriteMessage(ctx context.Context, containerName string, buf *bytes.Buffer) {
	if buf.Len() == 0 {
		return
	}

	fmt.Println("writeMessage", containerName, buf.Len())

	mergedName := getMergedContainerName(containerName)
	fileWriter, exists := a.currentFiles[mergedName]
	if !exists {
		ensureDir(path.Join(a.Config.Dir, mergedName))
		fileName := fmt.Sprintf("%s.txt", mergedName)
		fullFileName := path.Join(a.Config.Dir, mergedName, fileName)
		newFile, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Ctx(ctx).Err(err).Str("file_name", fullFileName).Msg("open new file error")
			return
		}
		fileWriter = newFile
		a.currentFiles[mergedName] = fileWriter
	}
	for {
		readBytes, err := buf.ReadBytes('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Ctx(ctx).Err(err).Msg("read bytes error")
			}
			break
		}
		fileWriter.Write(readBytes[8:]) // strip header from docker
	}
}

func getMergedContainerName(containerName string) string {
	if strings.Contains(containerName, ".") && !strings.HasPrefix(containerName, ".") {
		splitted := strings.Split(containerName, ".")
		return splitted[0]
	}
	return containerName
}
