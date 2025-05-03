package main

type MList[T any] interface {
	Push_tail(*mListNode[T])
	Push_head(*mListNode[T])
	Pop_tail()
	Pop_head()
	Head() *mListNode[T]
	Tail() *mListNode[T]
	Empty() bool
	Till() *mListNode[T]
	F_till() *mListNode[T]
	B_till() *mListNode[T]
	tar() *T

	Move_head(*mListNode[T])
	Move_tail(*mListNode[T])
	Delete(*mListNode[T])

	Init_with_num(int)
}

type mList[T any] struct {
	head   *mListNode[T]
	tail   *mListNode[T]
	cursor *mListNode[T]
	size   int
	itl    bool
}

func (tar *mList[T]) Push_tail(node *mListNode[T]) {
	if tar.head == nil {
		tar.head = node
		tar.tail = node
	} else {
		tar.tail.To_back(node)
		tar.tail = node
	}
}

func (tar *mList[T]) Push_head(node *mListNode[T]) {
	if tar.head == nil {
		tar.head = node
		tar.tail = node
	} else {
		tar.head.To_fore(node)
		tar.head = node
	}
}

func (tar *mList[T]) Pop_tail() {
	if tar.tail != nil {
		tar.tail = tar.tail.fore
	}
}

func (tar *mList[T]) Pop_head() {
	if tar.head != nil {
		tar.head = tar.head.back
	}
}

func (tar *mList[T]) Head() *mListNode[T] {
	return tar.head
}

func (tar *mList[T]) Tail() *mListNode[T] {
	return tar.tail
}

func (tar *mList[T]) Empty() bool {
	return tar.size <= 0
}

func (tar *mList[T]) Till() *mListNode[T] {
	return tar.B_till()
}

func (tar *mList[T]) F_till() *mListNode[T] {
	if !tar.itl {
		tar.cursor = tar.tail
	}
	tar.itl = true
	result := tar.cursor
	tar.cursor = tar.cursor.F_next()
	return result
}

func (tar *mList[T]) B_till() *mListNode[T] {
	if !tar.itl {
		tar.cursor = tar.head
	}
	tar.itl = true
	result := tar.cursor
	tar.cursor = tar.cursor.B_next()
	return result
}

func (tar *mList[T]) Move_head(node *mListNode[T]) {
	if tar.head != node {
		node.fore.back = node.back
		tar.head.fore = node
		node.back = tar.head
		node.fore = nil
		tar.head = node
	}
}

func (tar *mList[T]) Move_tail(node *mListNode[T]) {
	if tar.tail != node {
		node.back.fore = node.fore
		tar.tail.back = node
		node.fore = tar.tail
		node.back = nil
		tar.tail = node
	}
}

func (tar *mList[T]) Delete(node *mListNode[T]) {
	if node.back != nil {
		node.back.fore = node.fore
	} else {
		tar.head = node.fore
	}
}

func (tar *mList[T]) Init_with_num(num int) {
	tar.cursor = new(mListNode[T])
	tar.head = tar.cursor
	tar.tail = tar.cursor
	for range num - 1 {
		tar.Push_tail(&mListNode[T]{})
	}
}
