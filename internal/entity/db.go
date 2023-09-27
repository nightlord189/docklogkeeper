package entity

type LogDataDB struct {
	ID                int64  `gorm:"column:id"`
	ContainerName     string `gorm:"column:container_name"`
	ContainerFullName string `gorm:"column:container_full_name"`
	LogText           string `gorm:"column:log_text"`
	CreatedAt         int64  `gorm:"column:created_at"`
}

func (LogDataDB) TableName() string {
	return "log"
}

type ContainerDB struct {
	Name string `gorm:"column:name"`
}

func (ContainerDB) TableName() string {
	return "container"
}
