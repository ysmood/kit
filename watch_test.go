package gokit

import (
	"fmt"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestFsnWatch(t *testing.T) {
	str, _ := GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/watch/%s", str)

	OutputFile(p, str, nil)

	w, err := NewWatcher()
	E(err)

	go func() {
		for {
			e := <-w.Events

			if e.Op == fsnotify.Write {
				w.Close()
			}
		}
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		OutputFile(p, str, nil)
	}()

	w.Add(p)

	w.Start(300 * time.Millisecond)
}
func TestPollingWatch(t *testing.T) {
	str, _ := GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/watch/%s", str)

	OutputFile(p, str, nil)

	watcher, chEvents, chErrors := newWatcher()
	w := &Watcher{
		Events:  chEvents,
		Errors:  chErrors,
		watcher: watcher,
	}

	go func() {
		for {
			e := <-w.Events

			if e.Op == fsnotify.Write {
				w.Close()
			}
		}
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		OutputFile(p, str, nil)
	}()

	w.Add(p)

	w.Start(1 * time.Millisecond)
}
