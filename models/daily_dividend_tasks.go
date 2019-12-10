package models

import "time"

type DailyDividendTasks struct {
	Id             int       `orm:"column(id);pk;auto"`
	Time           string    `orm:"column(time)"` // 任务基本时间
	State          string    `orm:"column(state)"` // 任务完成状态
	CompletionTime time.Time `orm:"column(completion_time)"` // 任务具体的完成时间
}
