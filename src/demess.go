package gocross

type DeMess interface {
	Get_chan() chan *string
}

type MessQue interface {
	Write(message string)
	Read() *string
	Finish() bool
	Inish() bool
	Onish() bool
	IClose()
	OClose()
}

type messque struct {
	inish     bool
	onish     bool
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

func (tar *messque) Finish() bool {
	return tar.inish && tar.onish
}

func (tar *messque) Inish() bool {
	return tar.inish
}

func (tar *messque) Onish() bool {
	return tar.onish
}

func (tar *messque) IClose() {
	tar.inish = true
}

func (tar *messque) OClose() {
	tar.onish = true
}
