package main

type DeMess interface {
	Get_chan() chan *string
}

type MessQue interface {
	Write(message string)
	Read() *string
}

type messque struct {
	in_chan_  chan *string
	out_chan_ chan *string
	queue_    []string
}
type demess struct {
	side_   bool
	pi_que_ map[string]*messque
	na_que_ map[string]*messque
}

func (tar *messque) Write(message string) {
	tar.in_chan_ <- &message
}

func (tar *messque) Read() *string {
	return <-tar.out_chan_
}
