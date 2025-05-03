package main

type MQueue[T any] interface {
	Init()
	Push(T)
	Pop()
	The() T
	Empty() bool
	Size() int
}

type mQueue[T any] struct {
	ique_ mstack[T]
	oque_ mstack[T]
}

func (tar *mQueue[T]) Init() {
}

func (tar *mQueue[T]) Push(data T) {
	tar.ique_.Push(data)
}

func (tar *mQueue[T]) Pop() {
	if tar.oque_.Empty() {
		for !tar.ique_.Empty() {
			tar.oque_.Push(tar.ique_.The())
			tar.ique_.Pop()
		}
	}
	tar.oque_.Pop()
}

func (tar *mQueue[T]) The() T {
	if tar.oque_.Empty() && tar.ique_.Empty() {
		panic("队列为空，现在不让读")
	}
	if tar.oque_.Empty() {
		for !tar.ique_.Empty() {
			tar.oque_.Push(tar.ique_.The())
			tar.ique_.Pop()
		}
	}
	return tar.oque_.The()
}

func (tar *mQueue[T]) Empty() bool {
	return tar.ique_.Empty() && tar.oque_.Empty()
}

func (tar *mQueue[T]) Size() int {
	return tar.ique_.Size() + tar.oque_.Size()
}
