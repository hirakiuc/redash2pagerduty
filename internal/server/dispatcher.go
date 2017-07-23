package server

// https://gist.github.com/lestrrat/c9b78369cf9b9c5d9b0c909ed1e2452e

import (
	"context"
	"log"
	"sync"
)

// Dispatcher describe a task dispatcher.
type Dispatcher struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

// WorkFunc describe the type of worker function.
type WorkFunc func(context.Context)

// NewDispatcher return a Dispatcher instance.
func NewDispatcher(max int) *Dispatcher {
	return &Dispatcher{
		sem: make(chan struct{}, max),
	}
}

// Wait waits end of worker tasks.
func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

// Work add new worker task.
func (d *Dispatcher) Work(ctx context.Context, proc WorkFunc) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		d.work(ctx, proc)
	}()
}

func (d *Dispatcher) work(ctx context.Context, proc WorkFunc) {
	select {
	case <-ctx.Done():
		log.Printf("cancel work")
		return
	case d.sem <- struct{}{}:
		// got semaphore
		defer func() { <-d.sem }()
	}

	log.Printf("Worker Start\n")
	proc(ctx)
	log.Printf("Worker End\n")
}
