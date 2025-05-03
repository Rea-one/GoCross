package main

type MStack[T any] interface {
	Init()
	Push(T)
	Pop()
	The() T
	Empty() bool
	Size() int
}

type mstack[T any] struct {
	pool_   []T
	cursor_ int
}

func (tar *mstack[T]) Init() {
	tar.cursor_ = -1
}

func (tar *mstack[T]) Push(data T) {
	tar.cursor_++
	tar.pool_ = append(tar.pool_, data)
}

func (tar *mstack[T]) Pop() {
	if tar.Empty() {
		return
	}
	tar.cursor_--
}

func (tar *mstack[T]) The() T {
	if tar.Empty() {
		panic("栈为空，现在不让读")
	}
	return tar.pool_[tar.cursor_]
}

func (tar *mstack[T]) Empty() bool {
	return tar.cursor_ == -1
}

func (tar *mstack[T]) Size() int {
	return tar.cursor_ + 1
}
