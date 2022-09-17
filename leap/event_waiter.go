package leap

import "sync"

type resultOrError[T any] struct {
	result T
	err    error
}

// eventWaiter is used when you have an item that needs a response of some description that is async.
type eventWaiter[T any] struct {
	map_    map[string][]chan resultOrError[T]
	mapLock sync.Mutex
}

// wait is used to add into the map and wait for a response.
func (e *eventWaiter[T]) wait(tag string) (result T, err error) {
	// Create a channel to wait on.
	c := make(chan resultOrError[T])
	defer close(c)

	// Add the channel to the map.
	e.mapLock.Lock()
	if e.map_ == nil {
		e.map_ = map[string][]chan resultOrError[T]{}
	}
	e.map_[tag] = append(e.map_[tag], c)
	e.mapLock.Unlock()

	// Wait for the result.
	x := <-c
	return x.result, x.err
}

// signal is used to signal a result. Returns true if any channels were signalled to.
func (e *eventWaiter[T]) signal(tag string, result T, err error) bool {
	e.mapLock.Lock()
	if e.map_ == nil {
		e.map_ = map[string][]chan resultOrError[T]{}
	}
	s, ok := e.map_[tag]
	delete(e.map_, tag)
	e.mapLock.Unlock()
	if !ok {
		return false
	}
	for _, c := range s {
		c <- resultOrError[T]{result: result, err: err}
	}
	return true
}

// close is used to close all the channels and give an error to them all.
func (e *eventWaiter[T]) close(err error) {
	e.mapLock.Lock()
	if e.map_ == nil {
		e.map_ = map[string][]chan resultOrError[T]{}
	}
	for _, s := range e.map_ {
		for _, c := range s {
			c <- resultOrError[T]{err: err}
		}
	}
	e.map_ = nil
	e.mapLock.Unlock()
}
