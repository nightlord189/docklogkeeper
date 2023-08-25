package usecase

import (
	docker2 "github.com/nightlord189/docklogkeeper/internal/docker"
	"github.com/nightlord189/docklogkeeper/internal/log"
)

type Usecase struct {
	Docker *docker2.Adapter
	Log    *log.Adapter
}

func New(dockerInst *docker2.Adapter, logAdapter *log.Adapter) *Usecase {
	return &Usecase{
		Docker: dockerInst,
		Log:    logAdapter,
	}
}
