package log

import (
	"github.com/nightlord189/docklogkeeper/internal/config"
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/nightlord189/docklogkeeper/internal/repo"
	"time"
)

type Adapter struct {
	Config         config.LogConfig
	Repo           *repo.Repo
	names          map[string]string     //srv-captain--jsonbeautifier.1.qa9gcu6usinw06lqcfu286wsc -> srv-captain--jsonbeautifier
	lastTimestamps map[string]*time.Time //srv-captain--jsonbeautifier->"timestamp..."
	outputChannels []chan entity.LogDataDB
}

func New(cfg config.LogConfig, repoInst *repo.Repo, output []chan entity.LogDataDB) *Adapter {
	return &Adapter{
		Config:         cfg,
		Repo:           repoInst,
		names:          make(map[string]string, 10),
		lastTimestamps: make(map[string]*time.Time, 10),
		outputChannels: output,
	}
}
