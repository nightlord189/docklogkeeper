package repo

import (
	"fmt"
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"time"
)

func (r *Repo) GetMappedName(containerName string) string {
	var result string
	r.DB.Raw("SELECT container_name FROM container_mapping WHERE long_name = ?", containerName).Scan(&result)
	return result
}

func (r *Repo) SetMappedName(containerName, shortName string) {
	r.DB.Exec("insert into container_mapping (long_name, container_name) values (?, ?)", containerName, shortName)
}

func (r *Repo) InsertLogs(logs []entity.LogDataDB) error {
	return r.DB.Create(logs).Error
}

func (r *Repo) EnsureContainer(shortName string) error {
	return r.DB.Exec("insert into container (name) values (?) on conflict do nothing", shortName).Error
}

func (r *Repo) DeleteOldLogs(beforeThan int64) error {
	return r.DB.Exec("delete from log where created_at < ?", beforeThan).Error
}

func (r *Repo) DeleteContainersWithoutLogs() error {
	return r.DB.Exec(`delete from container where name in (select c.name as cname
from container c 
left join log on log.container_name = c.name
group by c.name
having count(log.id) = 0)`).Error
}

func (r *Repo) SearchLogs(shortName, like string) ([]entity.LogDataDB, error) {
	result := make([]entity.LogDataDB, 0, 100)
	err := r.DB.Where(`container_name = ? and log_text like ? order by id desc`, shortName, fmt.Sprintf("%%%s%%", like)).Find(&result).Error
	return result, err
}

func (r *Repo) GetLastTimestamp(shortName string) *time.Time {
	var entry entity.LogDataDB
	err := r.DB.Where("container_name = ?", shortName).Order("id desc").Limit(1).First(&entry).Error
	if err != nil || entry.CreatedAt == 0 {
		return nil
	}
	result := time.Unix(entry.CreatedAt, 0)
	return &result
}

func (r *Repo) GetContainers() ([]string, error) {
	var containers []entity.ContainerDB
	err := r.DB.Find(&containers).Error
	if err != nil {
		return []string{}, err
	}
	result := make([]string, len(containers))
	for i := range containers {
		result[i] = containers[i].Name
	}
	return result, nil
}

func (r *Repo) GetLogs(shortName string, greaterThan bool, cursor int64, limit int) ([]entity.LogDataDB, error) {
	var result []entity.LogDataDB
	query := r.DB.Where("container_name = ?", shortName)
	if cursor > 0 {
		if greaterThan {
			query = query.Where("id > ?", cursor)
		} else {
			query = query.Where("id < ?", cursor)
		}
	}
	query = query.Order("id desc").Limit(limit)
	err := query.Find(&result).Error
	return result, err
}
