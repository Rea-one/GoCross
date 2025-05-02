package main

type MList interface {
	Push(*mListNode)
	Pop_tail()
	Pop_head()
	Head() *mListNode
	Tail() *mListNode
}

type mList struct {
	head *mListNode
	tail *mListNode
}

func (tar *mList) Push(node *mListNode) {
	if tar.head == nil {
		tar.head = node
		tar.tail = node
	} else {
		tar.tail.To_back(node)
		tar.tail = node
	}
}

func (tar *mList) Pop_tail() {
	if tar.tail != nil {
		tar.tail = tar.tail.fore
	}
}

func (tar *mList) Pop_head() {
	if tar.head != nil {
		tar.head = tar.head.back
	}
}

func (tar *mList) Head() *mListNode {
	return tar.head
}

func (tar *mList) Tail() *mListNode {
	return tar.tail
}
