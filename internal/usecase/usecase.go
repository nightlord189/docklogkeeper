package usecase

import (
	docker2 "github.com/nightlord189/docklogkeeper/internal/docker"
	"github.com/nightlord189/docklogkeeper/internal/log"
	"github.com/nightlord189/docklogkeeper/internal/repo"
	"github.com/nightlord189/docklogkeeper/internal/trigger"
)

type Usecase struct {
	Repo    *repo.Repo
	Docker  *docker2.Adapter
	Log     *log.Adapter
	Trigger *trigger.Adapter
}

func New(repoInst *repo.Repo, dockerInst *docker2.Adapter, logAdapter *log.Adapter, triggerAdapter *trigger.Adapter) *Usecase {
	return &Usecase{
		Repo:    repoInst,
		Docker:  dockerInst,
		Log:     logAdapter,
		Trigger: triggerAdapter,
	}
}
