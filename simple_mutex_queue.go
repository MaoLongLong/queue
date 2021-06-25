package queue

import "sync"

type simpleMutexQueue struct {
	data []interface{}
	mu   sync.RWMutex
}

func newSimpleMutexQueue() *simpleMutexQueue {
	return &simpleMutexQueue{data: make([]interface{}, 0)}
}

func (q *simpleMutexQueue) Enqueue(v interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data, v)
}

func (q *simpleMutexQueue) Dequeue() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) == 0 {
		return nil
	}
	v := q.data[0]
	q.data = q.data[1:]
	return v
}

func (q *simpleMutexQueue) Empty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.data) != 0
}
