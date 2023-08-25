package log

import (
	"bufio"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"path"
	"strings"
)

func (a *Adapter) SearchLines(ctx context.Context, shortName string, req SearchRequest) ([]string, error) {
	fileEntries := a.getSortedFilesByDir(shortName)
	if len(fileEntries) == 0 {
		return nil, fmt.Errorf("there is no log files for this container")
	}

	lines := make([]string, 0, 10)
	for _, fileEntry := range fileEntries {
		file := openFile(ctx, path.Join(a.Config.Dir, shortName, fileEntry.Name()), true)
		if file == nil {
			continue
		}
		// Splits on newlines by default.
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), req.Contains) {
				lines = append(lines, scanner.Text())
			}
		}
		if err := file.Close(); err != nil {
			log.Ctx(ctx).Err(err).Msg("close file error")
		}
	}
	return lines, nil
}
