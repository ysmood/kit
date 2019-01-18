package gokit

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/radovskyb/watcher"
)

// Watcher cross platform fs watch, it will auto-fallback to polling if the
// fsnotify doesn't work with the OS
type Watcher struct {
	Events chan fsnotify.Event
	Errors chan error

	fsnWatcher *fsnotify.Watcher
	watcher    *watcher.Watcher
}

// NewWatcher ...
func NewWatcher() (*Watcher, error) {
	fsnWatcher := newFSNWatcher()

	var chEvents chan fsnotify.Event
	var chErrors chan error
	var w *watcher.Watcher

	if fsnWatcher == nil {
		w, chEvents, chErrors = newWatcher()
	} else {
		chEvents = make(chan fsnotify.Event)
		chErrors = make(chan error)
	}

	return &Watcher{
		Events:     chEvents,
		Errors:     chErrors,
		fsnWatcher: fsnWatcher,
		watcher:    w,
	}, nil
}

func newFSNWatcher() *fsnotify.Watcher {
	fsnWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil
	}

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		return nil
	}

	fsnWatcher.Add(tmpfile.Name())

	changed := false

	go func() {
		e, ok := <-fsnWatcher.Events

		if !ok {
			return
		}

		if e.Op&fsnotify.Remove == fsnotify.Remove {
			changed = true
		}
	}()

	os.Remove(tmpfile.Name())

	count := 300
	for count > 0 {
		count = count - 1
		time.Sleep(1 * time.Millisecond)

		if changed {
			return fsnWatcher
		}
	}

	return nil
}

func newWatcher() (*watcher.Watcher, chan fsnotify.Event, chan error) {
	w := watcher.New()
	chEvents := make(chan fsnotify.Event)
	chErrors := make(chan error)

	go func() {
		for {
			select {
			case event := <-w.Event:
				e := fsnotify.Event{
					Name: event.Path,
				}

				switch event.Op {
				case watcher.Create:
					e.Op = fsnotify.Create
				case watcher.Write:
					e.Op = fsnotify.Write
				case watcher.Remove:
					e.Op = fsnotify.Remove
				case watcher.Rename, watcher.Move:
					e.Op = fsnotify.Rename
				case watcher.Chmod:
					e.Op = fsnotify.Chmod
				}

				chEvents <- e
			case err := <-w.Error:
				chErrors <- err

			case <-w.Closed:
				return
			}
		}
	}()

	return w, chEvents, chErrors
}

// Start ...
func (w *Watcher) Start(interval time.Duration) error {
	if w.IsPolling() {
		return w.watcher.Start(interval)
	}

	for {
		select {
		case event, ok := <-w.fsnWatcher.Events:
			if !ok {
				return nil
			}
			w.Events <- event
		case err, ok := <-w.fsnWatcher.Errors:
			if !ok {
				return nil
			}
			w.Errors <- err
		}
	}
}

// Add ...
func (w *Watcher) Add(name string) error {
	if w.fsnWatcher != nil {
		return w.fsnWatcher.Add(name)
	}

	return w.watcher.Add(name)
}

// Remove ...
func (w *Watcher) Remove(name string) error {
	if w.fsnWatcher != nil {
		return w.fsnWatcher.Remove(name)
	}

	return w.watcher.Remove(name)
}

// Close ...
func (w *Watcher) Close() error {
	if w.fsnWatcher != nil {
		return w.fsnWatcher.Close()
	}

	w.watcher.Close()
	return nil
}

// IsPolling ...
func (w *Watcher) IsPolling() bool {
	return w.fsnWatcher == nil
}
