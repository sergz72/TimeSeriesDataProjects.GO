package core

type LruItem[T any] struct {
	Data *T
	Key  int
	next *LruItem[T]
	prev *LruItem[T]
}

type LruManager[T any] struct {
	head        *LruItem[T]
	tail        *LruItem[T]
	activeItems int
}

func (m *LruManager[T]) Add(key int, data *T) *LruItem[T] {
	i := &LruItem[T]{data, key, nil, m.head}
	m.Attach(i)
	return i
}

func (m *LruManager[T]) MoveToFront(item *LruItem[T]) {
	if item == m.head {
		return
	}
	m.Detach(item)
	m.Attach(item)
}

func (m *LruManager[T]) Attach(item *LruItem[T]) {
	item.prev = nil
	item.next = m.head
	if m.head != nil {
		m.head.prev = item
	} else {
		m.tail = item
	}
	m.head = item
	m.activeItems++
}

func (m *LruManager[T]) Detach(item *LruItem[T]) {
	if item.prev != nil {
		item.prev.next = item.next
	} else {
		m.head = item.next
	}
	if item.next != nil {
		item.next.prev = item.prev
	} else {
		m.tail = item.prev
	}
	m.activeItems--
}

func (m *LruManager[T]) GetTail() *LruItem[T] {
	return m.tail
}
