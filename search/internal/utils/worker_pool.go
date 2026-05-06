package utils

import (
	"context"
	"sync"
)

// WorkerPool manages a fixed number of workers for processing tasks
type WorkerPool struct {
	workerCount int
	taskQueue   chan func()
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewWorkerPool creates a new worker pool with the specified number of workers
func NewWorkerPool(ctx context.Context, workerCount int) *WorkerPool {
	if workerCount <= 0 {
		workerCount = 10 // Default to 10 workers
	}

	poolCtx, cancel := context.WithCancel(ctx)
	pool := &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan func(), workerCount*2), // Buffer 2x worker count
		ctx:         poolCtx,
		cancel:      cancel,
	}

	// Start workers
	for i := 0; i < workerCount; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	return pool
}

// Submit adds a task to the worker pool
func (p *WorkerPool) Submit(task func()) error {
	select {
	case p.taskQueue <- task:
		return nil
	case <-p.ctx.Done():
		return p.ctx.Err()
	}
}

// SubmitWithContext adds a task with a context check
func (p *WorkerPool) SubmitWithContext(ctx context.Context, task func()) error {
	select {
	case p.taskQueue <- task:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return p.ctx.Err()
	}
}

// worker processes tasks from the queue
func (p *WorkerPool) worker() {
	defer p.wg.Done()
	
	for {
		select {
		case task := <-p.taskQueue:
			if task != nil {
				task()
			}
		case <-p.ctx.Done():
			return
		}
	}
}

// Shutdown gracefully shuts down the worker pool
func (p *WorkerPool) Shutdown() {
	p.cancel()
	p.wg.Wait()
	close(p.taskQueue)
}

// ShutdownWithDrain shuts down after processing remaining tasks
func (p *WorkerPool) ShutdownWithDrain() {
	close(p.taskQueue)
	p.wg.Wait()
	p.cancel()
}