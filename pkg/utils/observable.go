package utils

import "sync"

// Observable is a thread-safe event helper
type Observable struct {
	subscribers *sync.Map
}

// Event of the observale
type Event interface{}

// Subscriber of the observable
type Subscriber chan Event

// NewObservable creates a new observable
func NewObservable() *Observable {
	return &Observable{
		subscribers: &sync.Map{},
	}
}

// Publish event to all subscribers
func (o *Observable) Publish(e Event) {
	// no need to lock this
	o.subscribers.Range(func(_, s interface{}) bool {
		s.(Subscriber) <- e
		return true
	})
}

// Subscribe returns a subscriber to emit events
func (o *Observable) Subscribe() Subscriber {
	subscriber := make(Subscriber)

	o.subscribers.Store(&subscriber, subscriber)

	return subscriber
}

// Unsubscribe from the observable
func (o *Observable) Unsubscribe(s *Subscriber) {
	o.subscribers.Delete(s)
}

// Until filter returns true keep waiting
func (o *Observable) Until(filter func(e Event) bool) (e Event) {
	s := o.Subscribe()
	defer o.Unsubscribe(&s)

	for e = range s {
		if filter(e) {
			break
		}
	}

	return
}
