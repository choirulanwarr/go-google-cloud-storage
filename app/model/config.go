package model

const (
	AppStatusKey = "app_status"
)

type Config struct {
	ID    uint   `gorm:"column:id;primaryKey"`
	Key   string `gorm:"column:key"`
	Value string `gorm:"column:value"`
}
