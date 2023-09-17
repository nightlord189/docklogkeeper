package usecase

import (
	"context"
	"fmt"
	"github.com/nightlord189/docklogkeeper/internal/entity"
)

func (u *Usecase) CreateTrigger(ctx context.Context, item *entity.TriggerDB) error {
	if err := u.Repo.CreateEntity(item); err != nil {
		return fmt.Errorf("create db entity error: %w", err)
	}
	u.Trigger.ReloadCache(ctx)
	return nil
}

func (u *Usecase) UpdateTrigger(ctx context.Context, item *entity.TriggerDB) error {
	if err := u.Repo.UpdateEntity(item); err != nil {
		return fmt.Errorf("update db entity error: %w", err)
	}
	u.Trigger.ReloadCache(ctx)
	return nil
}

func (u *Usecase) DeleteTrigger(ctx context.Context, id int64) error {
	if err := u.Repo.DeleteTrigger(id); err != nil {
		return fmt.Errorf("delete db entity error: %w", err)
	}
	u.Trigger.ReloadCache(ctx)
	return nil
}
