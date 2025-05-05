package gocross

type MListNode[T any] interface {
	Init(T)
	Move_tail()
	B_next() *mListNode[T]
	Move_head()
	F_next() *mListNode[T]
	To_back(*mListNode[T])
	To_fore(*mListNode[T])
	Get() T
	Empty() bool
	Till() *mListNode[T]
	F_till() *mListNode[T]
	B_till() *mListNode[T]
}

type mListNode[T any] struct {
	fore *mListNode[T]
	back *mListNode[T]
	data T
}

func (tar *mListNode[T]) Init(data T) {
	tar.data = data
	tar.back = nil
	tar.fore = nil
}

func (tar *mListNode[T]) Move_tail() {
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

func (tar *mListNode[T]) B_next() *mListNode[T] {
	return tar.back
}

func (tar *mListNode[T]) Move_head() {
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

func (tar *mListNode[T]) F_next() *mListNode[T] {
	return tar.fore
}

func (tar *mListNode[T]) To_back(node *mListNode[T]) {
	tar.back = node
}

func (tar *mListNode[T]) To_fore(node *mListNode[T]) {
	tar.fore = node
}

func (tar *mListNode[T]) Get() T {
	return tar.data
}

func (tar *mListNode[T]) Empty() bool {
	return tar.back == nil && tar.fore == nil
}

func (tar *mListNode[T]) F_till() *mListNode[T] {
	result := tar
	tar = tar.fore
	return result
}

func (tar *mListNode[T]) B_till() *mListNode[T] {
	result := tar
	tar = tar.back
	return result
}

func (tar *mListNode[T]) Till() *mListNode[T] {
	return tar.B_till()
}
