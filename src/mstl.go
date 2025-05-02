package main

type MListNode interface {
	Init(receiver)
	Move_tail()
	B_next() *mListNode
	Move_head()
	F_next() *mListNode
	To_back(*mListNode)
	To_fore(*mListNode)
	Get() *receiver
}

type mListNode struct {
	fore *mListNode
	back *mListNode
	data receiver
}

func (tar *mListNode) Init(data receiver) {
	tar.data = data

}

func (tar *mListNode) Move_tail() {
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

func (tar *mListNode) B_next() *mListNode {
	return tar.back
}

func (tar *mListNode) Move_head() {
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

func (tar *mListNode) F_next() *mListNode {
	return tar.fore
}

func (tar *mListNode) To_back(node *mListNode) {
	tar.back = node
}

func (tar *mListNode) To_fore(node *mListNode) {
	tar.fore = node
}

func (tar *mListNode) Get() *receiver {
	return &tar.data
}
