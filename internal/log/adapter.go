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

func New(cfg config.LogConfig) (*Adapter, error) {
	adapter := &Adapter{
		Config:       cfg,
		names:        make(map[string]string, 10),
		currentFiles: make(map[string]*os.File, 10),
	}
	if err := adapter.init(); err != nil {
		return nil, err
	}
	return adapter, nil
}

func (a *Adapter) init() error {
	if _, err := os.Stat(a.Config.Dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(a.Config.Dir, os.ModePerm)
		if err != nil {
			return err
		}
		log.Ctx(context.Background()).Info().Msgf("output directory created: %s", a.Config.Dir)
	}
	log.Ctx(context.Background()).Info().Msgf("output directory already exists: %s", a.Config.Dir)
	return nil
}

func (a *Adapter) Close() {
	for _, f := range a.currentFiles {
		f.Close()
	}
}
