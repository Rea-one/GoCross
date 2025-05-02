package main

import (
	"time"
)

type Task interface {
	GetID() string
	GetQuery() string
	GetResult() string
	SetResult(string)
	GetDeadline() time.Time
}

type task struct {
	ID       string // 任务唯一标识
	Query    string // SQL语句
	Result   string
	Deadline time.Time // 超时控制
}

func (tar *task) GetID() string {
	return tar.ID
}

func (tar *task) GetQuery() string {
	return tar.Query
}

func (tar *task) GetResult() string {
	return tar.Result
}

func (tar *task) SetResult(result string) {
	tar.Result = result
}

func (tar *task) GetDeadline() time.Time {
	return tar.Deadline
}
