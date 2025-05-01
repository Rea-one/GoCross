package main

import (
	"time"
)

type task struct {
	ID       string // 任务唯一标识
	Query    string // SQL语句
	Result   string
	Deadline time.Time // 超时控制
}
