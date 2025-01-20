package priorityqueue

type Item[T any] struct {
	Value    T   // The value of the item; arbitrary.
	priority int // The priority of the item in the queue.
}

func NewItem[T any](value T, priority int) *Item[T] {
	return &Item[T]{
		Value:    value,
		priority: priority,
	}
}
