package entity

type LogDataDB struct {
	ID            int64  `gorm:"id"`
	ContainerName string `gorm:"container_name"`
	LogText       string `gorm:"log_text"`
	CreatedAt     int64  `gorm:"created_at"`
}

func (LogDataDB) TableName() string {
	return "log"
}

type ContainerDB struct {
	Name string `gorm:"name"`
}

func (ContainerDB) TableName() string {
	return "container"
}
