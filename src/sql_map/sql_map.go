package sqlmap

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
	var rec string
	var cursor *string
	for word := range task.Query {
		if word == ' ' {
			if rec != "" {
				now := tar.themap.To_SQL(rec)
				if now != "reject" {
					if cursor == nil {
						switch now {
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
						}
					} else {
						*cursor = now
						cursor = nil
					}
				} else {
					task.State = "reject"
					return
				}
			}
		} else {
			rec += string(word)
		}
	}
	if rec != "" {
		now := tar.themap.To_SQL(rec)
		if now != "reject" {
			if cursor != nil {
				*cursor = now
				cursor = nil
			}
		} else {
			task.State = "reject"
		}
	}
}
