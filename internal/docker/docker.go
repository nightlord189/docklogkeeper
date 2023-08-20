package docker

import "github.com/nightlord189/docklogkeeper/internal/config"

type Adapter struct {
	Config config.Config
}

func New(cfg config.Config) *Adapter {
	return &Adapter{Config: cfg}
}
