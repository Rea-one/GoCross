package main

type IOInterface interface {
	Init(host string, dom string, database string, password string)
	Start()
	Stop()
	Wait()
	ReadIn()
	ReadOut()
	WriteIn(message []string)
	WriteOut(message []string)
}

type IOMessage struct {
	in_mess_  chan []string
	out_mess_ chan []string
}

func (tar IOMessage) ReadIn() {
	if len(tar.in_mess_) > 0 {
		return <-tar.in_mess_
	} else {
		return " "
	}
}
