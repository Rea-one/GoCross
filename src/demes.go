package main

type DePass interface {
	get_chan() chan *string
}

type DeMes struct {
	side_    bool
	pi_que_  []string
	na_que_  []string
	pi_chan_ chan *string
	na_chan_ chan *string
}

func (tar DeMes) get_chan() chan *string {
	if tar.side_ {
		return tar.pi_chan_
	} else {
		return tar.na_chan_
	}
}

func (tar DeMes)
