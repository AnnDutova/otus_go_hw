package list

import "fmt"

type List interface {
	Len() int
	Front() *Item
	Back() *Item
	PushFront(v interface{}) *Item
	PushBack(v interface{}) *Item
	Remove(i *Item)
	MoveToFront(i *Item)
	PrintList()
}

type Item struct {
	Value interface{}
	Next  *Item
	Prev  *Item
}

type list struct {
	head *Item
	tail *Item
	len  int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *Item {
	return l.head
}

func (l *list) Back() *Item {
	return l.tail
}

func (l *list) PushFront(v interface{}) *Item {
	item := &Item{Value: v, Next: nil, Prev: nil}
	if l.head == nil {
		l.head = item
		l.tail = item
	} else {
		item.Next = l.head
		l.head.Prev = item
		l.head = item
	}
	l.len++
	return l.head
}

func (l *list) PushBack(v interface{}) *Item {
	newTail := &Item{Value: v, Next: nil, Prev: l.tail}
	if l.head == nil {
		l.head = newTail
		l.tail = l.head
	} else {
		l.tail.Next = newTail
		l.tail = newTail
	}

	l.len++
	return l.tail
}

func (l *list) Remove(i *Item) {
	elem := l.head
	for elem != nil {
		if elem.Value == i.Value {
			switch {
			case l.head == elem:
				newHead := l.head.Next
				l.head = newHead
			case l.tail == elem:
				l.tail = elem.Prev
				elem.Prev.Next = nil
			default:
				pr := elem.Prev
				pr.Next = elem.Next
				elem.Next.Prev = pr
			}
			l.len--
			return
		}
		elem = elem.Next
	}
}

func (l *list) MoveToFront(i *Item) {
	elem := l.head
	for elem != nil {
		if elem.Value == i.Value {
			switch {
			case l.head == elem:
			case l.tail == elem:
				l.tail = l.tail.Prev
				l.tail.Next = nil

				elem.Prev = nil
				elem.Next = l.head
				l.head.Prev = elem
				l.head = elem
			default:
				pr := elem.Prev
				pr.Next = elem.Next
				elem.Next.Prev = pr

				l.head.Prev = elem
				elem.Next = l.head
				l.head = elem
			}
			return
		}
		elem = elem.Next
	}
}

func NewList() List {
	return new(list)
}

func (l *list) PrintList() {
	start := l.head
	var res string
	var i int
	for start != nil {
		res += fmt.Sprintf("elem[%d]: %v ", i, start.Value.(int))
		i++
		start = start.Next
	}
	fmt.Println(res)
}
