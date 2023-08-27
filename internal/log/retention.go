package log

import (
	"context"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"time"
)

func (a *Adapter) ClearOldFiles(ctx context.Context) {
	if _, err := os.Stat(a.Config.Dir); err != nil {
		log.Ctx(ctx).Err(err).Msg("stat root dir error")
		return
	}
	folders, err := os.ReadDir(a.Config.Dir)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("read root dir error")
		return
	}

	deletedFiles := 0

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}
		files, err := os.ReadDir(path.Join(a.Config.Dir, folder.Name()))
		if err != nil {
			log.Ctx(ctx).Err(err).Msgf("read dir error [dir %s]", folder.Name())
			continue
		}
		if len(files) == 0 { // delete empty directories without log files
			if err := os.Remove(path.Join(a.Config.Dir, folder.Name())); err != nil {
				log.Ctx(ctx).Err(err).Msgf("delete empty directory error [dir %s]", folder.Name())
			}
		}
		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				log.Ctx(ctx).Err(err).Msgf("get file info error [dir %s, file %s]", folder.Name(), file.Name())
				continue
			}
			diff := time.Now().Sub(info.ModTime())
			if diff.Seconds() >= float64(a.Config.Retention) {
				if err := os.Remove(path.Join(a.Config.Dir, folder.Name(), file.Name())); err != nil {
					log.Ctx(ctx).Err(err).Msgf("delete file error [dir %s, file %s]", folder.Name(), file.Name())
				}
				deletedFiles++
			}
		}
	}
	if deletedFiles > 0 {
		log.Ctx(ctx).Info().Msgf("deleted old files: %d", deletedFiles)
	}
}
