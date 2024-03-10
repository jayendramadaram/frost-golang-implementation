package collections

import "sync"

type OrderedList[T comparable] struct {
	Items []T
	mu    sync.RWMutex
}

// todo: use HashSet
func NewOrderedList[T comparable]() *OrderedList[T] {
	return &OrderedList[T]{
		Items: []T{},
		mu:    sync.RWMutex{},
	}
}

func (o *OrderedList[T]) Add(item T) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.Items = append(o.Items, item)
}

func (o *OrderedList[T]) Remove(item T) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for i, v := range o.Items {
		if v == item {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			break
		}
	}
}

func (o *OrderedList[T]) RemoveAll(item T) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for i, v := range o.Items {
		if v == item {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
		}
	}
}

func (o *OrderedList[T]) Contains(item T) bool {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for _, v := range o.Items {
		if v == item {
			return true
		}
	}

	return false
}

func (o *OrderedList[T]) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return len(o.Items)
}

func (o *OrderedList[T]) RemoveAt(index int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.Items = append(o.Items[:index], o.Items[index+1:]...)
}

func (o *OrderedList[T]) Get(index int) T {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.Items[index]
}
