package sqlmap

import "strings"

type ActSqlMap interface {
	Init()
	To_SQL(string) string
	analyze(string) string
	Ana(*Task)
}

type SqlMap struct {
	themap Reader
}

func (tar *SqlMap) Init() {
	tar.themap.Init("./src/sql_map/SQL.json")
}

func (tar *SqlMap) To_SQL(message string) string {
	return tar.analyze(message)
}

func (tar *SqlMap) analyze(message string) string {
	var rec string
	var result string
	for word := range message {
		if word == ' ' {
			if rec != "" {
				now := tar.themap.To_SQL(rec)
				if now != "reject" {
					result += now
				} else {
					return "reject"
				}
			}
		} else {
			rec += string(word)
		}
	}
	// 处理最后一个非空片段
	now := tar.themap.To_SQL(rec)
	if now != "reject" {
		result += now
	} else {
		return "reject"
	}
	return result
}

func (tar *SqlMap) Ana(task *Task) {
	var cursor *string

	// 将 Query 分割为 token
	tokens := strings.Fields(task.Message)
	switch len(tokens) {
	case 0:
		task.State = "reject"
		return
	case 1:
		task.SetState(tokens[0])
		return
	}

	for _, word := range tokens {
		now := tar.themap.To_SQL(word)
		if cursor == nil {
			switch now {
			case "at":
				cursor = &task.At
			case "sender":
				cursor = &task.Sender
			case "receiver":
				cursor = &task.Receiver
			case "type":
				cursor = &task.Ttype
			case "password":
				cursor = &task.Password
			case "message":
				cursor = &task.Message
			default:
				// 如果不是字段名，直接跳过
				continue
			}
		} else {
			*cursor = now
			cursor = nil
		}
	}

	// 如果最后仍有未赋值的 cursor，说明语法不完整
	if cursor != nil {
		task.State = "reject"
	}
}
