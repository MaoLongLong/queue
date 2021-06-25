package queue

import (
	"math/rand"
	"testing"
)

func BenchmarkQueue(b *testing.B) {

	queues := map[string]Queue{
		"lock_free_queue":    NewLockFreeQueue(),
		"simple_mutex_queue": newSimpleMutexQueue(),
	}

	for name, q := range queues {
		b.Run(name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					op := rand.Intn(100)
					if op < 40 {
						q.Dequeue()
					} else {
						q.Enqueue(1)
					}
				}
			})
		})
	}
}
