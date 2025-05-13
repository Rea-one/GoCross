package sqlmap

import (
	"time"
)

type ActTask interface {
	GetID() string
	GetState() string
	GetFeedback() string
	SetState(string)
	GetDeadline() time.Time
}

type Task struct {
	ID        string
	At        string
	State     string
	Sender    string
	Receiver  string
	Password  string
	Ttype     string
	Message   string
	Feedback  string
	ImageID   string
	ImageURL  string
	TimeStamp string
	Deadline  time.Time // 超时控制
}

func (tar *Task) GetID() string {
	return tar.ID
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

func (tar *Task) GetFeedback() string {
	return tar.Feedback
}
