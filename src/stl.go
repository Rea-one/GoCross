package main

type ListNode interface {
	Init(interface{})
	Move_tail()
	B_next() *listNode
	Move_head()
	F_next() *listNode
	To_back(*listNode)
	To_fore(*listNode)
}

type listNode struct {
	fore *listNode
	back *listNode
	data interface{}
}

func (tar *listNode) Init(data interface{}) {
	tar.data = data

}

func (tar *listNode) Move_tail() {
	_cursor_ := tar
	for tar.back != nil {
		tar = tar.back
	}
	_cursor_.fore.back = _cursor_.back
	_cursor_.back.fore = _cursor_.fore
	tar.back = _cursor_
	_cursor_.fore = tar
	tar = _cursor_
}

func (tar *listNode) B_next() *listNode {
	return tar.back
}

func (tar *listNode) Move_head() {
	_cursor_ := tar
	for tar.fore != nil {
		tar = tar.fore
	}
	_cursor_.back.fore = _cursor_.fore
	_cursor_.fore.back = _cursor_.back
	tar.fore = _cursor_
	_cursor_.back = tar
	tar = _cursor_
}

func (tar *listNode) F_next() *listNode {
	return tar.fore
}
