package log

import (
	"fmt"
	"time"
)

type logDataDB struct {
	ID            int64  `gorm:"id"`
	ContainerName string `gorm:"container_name"`
	LogText       string `gorm:"log_text"`
	CreatedAt     int64  `gorm:"created_at"`
}

func (logDataDB) TableName() string {
	return "log"
}

type containerDB struct {
	Name string `gorm:"name"`
}

func (containerDB) TableName() string {
	return "container"
}

func (a *Adapter) getMappedName(containerName string) string {
	var result string
	a.DB.Raw("SELECT container_name FROM container_mapping WHERE long_name = ?", containerName).Scan(&result)
	return result
}

func (a *Adapter) setMappedName(containerName, shortName string) {
	a.DB.Exec("insert into container_mapping (long_name, container_name) values (?, ?)", containerName, shortName)
}

func (a *Adapter) insertLogs(logs []logDataDB) error {
	return a.DB.Create(logs).Error
}

func (a *Adapter) ensureContainer(shortName string) error {
	return a.DB.Exec("insert into container (name) values (?) on conflict do nothing", shortName).Error
}

func (a *Adapter) deleteOldLogs(beforeThan int64) error {
	return a.DB.Exec("delete from log where created_at < ?", beforeThan).Error
}

func (a *Adapter) DeleteContainersWithoutLogs() error {
	return a.DB.Exec(`delete from container where name in (select c.name as cname
from container c 
left join log on log.container_name = c.name
group by c.name
having count(log.id) = 0)`).Error
}

func (a *Adapter) searchLogs(shortName, like string) ([]logDataDB, error) {
	result := make([]logDataDB, 0, 100)
	err := a.DB.Where(`container_name = ? and log_text like ? order by id desc`, shortName, fmt.Sprintf("%%%s%%", like)).Find(&result).Error
	return result, err
}

func (a *Adapter) getLastTimestampFromDB(shortName string) *time.Time {
	var entry logDataDB
	err := a.DB.Where("container_name = ?", shortName).Order("id desc").Limit(1).First(&entry).Error
	if err != nil || entry.CreatedAt == 0 {
		return nil
	}
	result := time.Unix(entry.CreatedAt, 0)
	return &result
}

func (a *Adapter) getContainers() ([]string, error) {
	var containers []containerDB
	err := a.DB.Find(&containers).Error
	if err != nil {
		return []string{}, err
	}
	result := make([]string, len(containers))
	for i := range containers {
		result[i] = containers[i].Name
	}
	return result, nil
}

func (a *Adapter) getLogs(shortName string, greaterThan bool, cursor int64, limit int) ([]logDataDB, error) {
	var result []logDataDB
	query := a.DB.Where("container_name = ?", shortName)
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
