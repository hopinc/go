package leap

import (
	"sync"
	"sync/atomic"
)

// queueDispatcher is used to handle dispatching messages into a managed channel. This is useful for un-buffered
// channels or channels where it is unknown if the buffer is saturated since it will stop blocking on the main
// read loop (very bad).
type queueDispatcher[T any] struct {
	channel chan T

	active uintptr

	events     []T
	eventsLock sync.Mutex
}

func (q *queueDispatcher[T]) dispatch(item T) {
	// Add the event to the queue.
	q.eventsLock.Lock()
	q.events = append(q.events, item)
	q.eventsLock.Unlock()

	// If we are already active, then we don't need to do anything else.
	if atomic.SwapUintptr(&q.active, 1) == 1 {
		return
	}

	// Start a goroutine to handle flushing.
	go func() {
		// When we are done, we need to set active to 0.
		defer atomic.StoreUintptr(&q.active, 0)

		// This intentionally goes round twice to handle a flood of events.
		for {
			// Get the events.
			q.eventsLock.Lock()
			events := q.events
			q.events = nil
			q.eventsLock.Unlock()

			// If events is nil, then we are done.
			if events == nil {
				return
			}

			// Dispatch the events.
			for _, event := range events {
				q.channel <- event
			}
		}
	}()
}

func newQueueDispatcher[T any](c chan T) *queueDispatcher[T] {
	x := &queueDispatcher[T]{channel: c}
	return x
}
