package usecase

import (
	"context"
	"fmt"
	"sort"

	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/rs/zerolog/log"
)

func (u *Usecase) GetContainers(ctx context.Context) ([]entity.ContainerInfo, error) {
	allContainers, err := u.Repo.GetContainers()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("get all containers error")
		return nil, fmt.Errorf("get all containers error: %w", err)
	}

	aliveContainers, err := u.Docker.GetAliveContainers(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("get alive containers error")
		return nil, fmt.Errorf("get alive containers error: %w", err)
	}

	u.convertToShortNames(aliveContainers)

	allMap := arrToMap(allContainers)
	aliveMap := arrToMap(aliveContainers) // move to shot name

	result := make([]entity.ContainerInfo, 0, len(aliveContainers)+len(allContainers))
	for cont := range allMap {
		info := entity.ContainerInfo{
			ShortName: cont,
			IsAlive:   false,
		}
		if _, ok := aliveMap[cont]; ok {
			info.IsAlive = true
		}
		result = append(result, info)
	}

	// for case when some container didn't have any logs (not in allContainers)
	for cont := range aliveMap {
		if _, ok := allMap[cont]; !ok {
			result = append(result, entity.ContainerInfo{
				ShortName: cont,
				IsAlive:   true,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ShortName < result[j].ShortName
	})

	return result, nil
}
