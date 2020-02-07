package utils

import (
	"context"
	"sync"
	"sync/atomic"
)

// Observable is a thread-safe event helper
type Observable struct {
	subscribers *sync.Map
	idCount     int64
}

// Event of the observale
type Event interface{}

// Subscriber of the observable
type Subscriber struct {
	C  chan Event
	id int64
}

// NewObservable creates a new observable
func NewObservable() *Observable {
	return &Observable{
		subscribers: &sync.Map{},
	}
}

// Publish event to all subscribers
func (o *Observable) Publish(e Event) {
	o.subscribers.Range(func(_, s interface{}) (goOn bool) {
		goOn = true
		defer func() { _ = recover() }()
		s.(*Subscriber).C <- e
		return
	})
}

// Subscribe returns a subscriber to emit events
func (o *Observable) Subscribe() *Subscriber {
	id := atomic.AddInt64(&o.idCount, 1)

	subscriber := &Subscriber{
		C:  make(chan Event),
		id: id,
	}

	o.subscribers.Store(id, subscriber)

	return subscriber
}

// Unsubscribe from the observable
func (o *Observable) Unsubscribe(s *Subscriber) {
	close(s.C)
	o.subscribers.Delete(s.id)
}

// Count returns the number of subscribers
func (o *Observable) Count() int {
	c := 0
	o.subscribers.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

// Until check returns true keep waiting
func (o *Observable) Until(ctx context.Context, check func(e Event) bool) (Event, error) {
	s := o.Subscribe()
	defer o.Unsubscribe(s)

	for {
		select {
		case e := <-s.C:
			if check(e) {
				return e, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
