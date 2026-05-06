package application

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ConcurrencyEnhancer provides comprehensive concurrency improvements
type ConcurrencyEnhancer struct {
	workerPool     *WorkerPool
	goroutineCount int64
	activeWorkers  int64
	taskQueue      chan Task
	shutdown       chan struct{}
	wg             sync.WaitGroup
}

// Task represents a unit of work
type Task struct {
	ID       string
	Function func(ctx context.Context) interface{}
	Context  context.Context
	Priority int
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	size     int
	workers  []*Worker
	taskChan chan Task
	shutdown chan struct{}
	wg       sync.WaitGroup
}

// Worker represents a single worker
type Worker struct {
	id       int
	taskChan chan Task
	quit     chan struct{}
}

// NewConcurrencyEnhancer creates a new concurrency enhancer
func NewConcurrencyEnhancer() *ConcurrencyEnhancer {
	size := runtime.NumCPU() * 2
	ce := &ConcurrencyEnhancer{
		taskQueue: make(chan Task, 1000),
		shutdown:  make(chan struct{}),
	}

	ce.workerPool = NewWorkerPool(size, ce.taskQueue)
	go ce.startMonitoring()

	return ce
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(size int, taskChan chan Task) *WorkerPool {
	wp := &WorkerPool{
		size:     size,
		taskChan: taskChan,
		shutdown: make(chan struct{}),
	}

	for i := 0; i < size; i++ {
		worker := &Worker{
			id:       i,
			taskChan: taskChan,
			quit:     make(chan struct{}),
		}
		wp.workers = append(wp.workers, worker)
		wp.wg.Add(1)
		go worker.start(&wp.wg)
	}

	return wp
}

// start starts the worker
func (w *Worker) start(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task := <-w.taskChan:
			w.processTask(task)
		case <-w.quit:
			return
		}
	}
}

// processTask processes a single task
func (w *Worker) processTask(task Task) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Worker %d: Task %s panicked: %v", w.id, task.ID, r)
		}
	}()

	task.Function(task.Context)
}

// SafeGoroutine starts a goroutine with panic recovery
func (ce *ConcurrencyEnhancer) SafeGoroutine(name string, fn func()) {
	atomic.AddInt64(&ce.goroutineCount, 1)

	go func() {
		defer func() {
			atomic.AddInt64(&ce.goroutineCount, -1)
			if r := recover(); r != nil {
				log.Printf("Goroutine %s panicked: %v", name, r)
			}
		}()

		fn()
	}()
}

// SubmitTask submits a task to the worker pool
func (ce *ConcurrencyEnhancer) SubmitTask(task Task) error {
	select {
	case ce.taskQueue <- task:
		return nil
	default:
		return fmt.Errorf("task queue is full")
	}
}

// startMonitoring starts monitoring goroutines and performance
func (ce *ConcurrencyEnhancer) startMonitoring() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			count := runtime.NumGoroutine()
			log.Printf("Concurrency Monitor: %d goroutines active", count)

			if count > 1000 {
				log.Printf("WARNING: High goroutine count detected: %d", count)
			}
		case <-ce.shutdown:
			return
		}
	}
}

// Shutdown gracefully shuts down the concurrency enhancer
func (ce *ConcurrencyEnhancer) Shutdown(ctx context.Context) error {
	close(ce.shutdown)

	// Shutdown worker pool
	for _, worker := range ce.workerPool.workers {
		close(worker.quit)
	}

	ce.workerPool.wg.Wait()
	return nil
}
