package log

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

func (a *Adapter) GetAllContainers(ctx context.Context) ([]string, error) {
	if _, err := os.Stat(a.Config.Dir); err != nil {
		log.Ctx(ctx).Err(err).Msg("stat root dir error")
		return nil, fmt.Errorf("stat root dir error: %w", err)
	}

	folders, err := os.ReadDir(a.Config.Dir)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("read root dir error")
		return nil, fmt.Errorf("read root dir error: %w", err)
	}

	result := make([]string, 0, len(folders))
	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}
		result = append(result, folder.Name())
	}
	return result, nil
}
