package sqlmap

import (
	"encoding/json"
	"os"
)

type ActReader interface {
	Init(string)
	To_SQL(string) string
}

type Reader struct {
	// 根据json文件初始化，形成sql过滤与翻译
	Classes map[string]string `json:"classes"`
	Rejects map[string]bool   `json:"rejects"`
}

func (tar *Reader) Init(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		os.Exit(1)
	}
	err = json.Unmarshal(data, tar)
	if err != nil {
		os.Exit(1)
	}
}

// 将缩写扩写， 比如 "a" -> "at"
func (tar *Reader) To_SQL(message string) string {
	_, r := tar.Rejects[message]
	mess, c := tar.Classes[message]
	if r {
		return "reject"
	} else if c {
		return mess
	} else {
		return message
	}
}
