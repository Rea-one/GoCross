package sqlmap

type ActSqlMap interface {
	Init()
	To_SQL(string) string
	analyze(string) string
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
