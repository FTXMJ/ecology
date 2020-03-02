package models

import "time"

type DailyDividendTasks struct {
	Id             int       `gorm:"column:id;primary_key"`
	Time           string    `gorm:"column:time"`            // 任务基本时间
	State          string    `gorm:"column:state"`           // 任务完成状态
	CompletionTime time.Time `gorm:"column:completion_time"` // 任务具体的完成时间
}
