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
	classes map[string]string `json:"classes"`
	rejects map[string]string `json:"rejects"`
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

func (tar *Reader) To_SQL(message string) string {
	_, r := tar.rejects[message]
	_, c := tar.classes[message]
	if r {
		return "reject"
	} else if c {
		return tar.classes[message]
	} else {
		return message
	}
}
