package sqlmap

import (
	"time"
)

type ActTask interface {
	GetID() float64
	GetQuery() string
	GetState() string
	SetState(string)
	GetDeadline() time.Time
}

type Task struct {
	ID       float64 // 任务唯一标识
	Sender   string
	Receiver string
	Password string
	State    string
	Ttype    string
	Query    string // 自定义语句，遵循sql_map的处理规则
	SQL_fb   string
	Message  string
	Feedback string
	Deadline time.Time // 超时控制
}

func (tar *Task) GetID() float64 {
	return tar.ID
}

func (tar *Task) GetQuery() string {
	return tar.Query
}

func (tar *Task) GetState() string {
	return tar.State
}

func (tar *Task) SetState(result string) {
	tar.State = result
}

func (tar *Task) GetDeadline() time.Time {
	return tar.Deadline
}
