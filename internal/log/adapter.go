package log

import (
	"context"
	"errors"
	"github.com/nightlord189/docklogkeeper/internal/config"
	"github.com/rs/zerolog/log"
	"os"
)

type Adapter struct {
	Config       config.LogConfig
	names        map[string]string   //srv-captain--jsonbeautifier.1.qa9gcu6usinw06lqcfu286wsc -> srv-captain--jsonbeautifier
	currentFiles map[string]*os.File //srv-captain--jsonbeautifier -> srv-captain--jsonbeautifier-log-1.txt
}

func New(cfg config.LogConfig) *Adapter {
	adapter := &Adapter{
		Config:       cfg,
		names:        make(map[string]string, 10),
		currentFiles: make(map[string]*os.File, 10),
	}
	ensureDir(adapter.Config.Dir)
	return adapter
}

func ensureDir(dir string) {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Ctx(context.Background()).Info().Msgf("directory create error: %s %v", dir, err)
			return
		}
		log.Ctx(context.Background()).Info().Msgf("directory created: %s", dir)
		return
	}
	log.Ctx(context.Background()).Info().Msgf("directory already exists: %s", dir)
}

func (a *Adapter) Close() {
	for _, f := range a.currentFiles {
		f.Close()
	}
}
