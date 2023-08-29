package log

import (
	"context"
)

func (a *Adapter) GetAllContainers(ctx context.Context) ([]string, error) {
	return a.Repo.GetContainers()
}
