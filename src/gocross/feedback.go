package gocross

type ActFeedback interface {
}

type feedback struct {
	At        string `json:"at"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Timestamp string `json:"timestamp"`
	State     string `json:"state"`
	Message   string `json:"message"`
	Image     string `json:"image"`
}
