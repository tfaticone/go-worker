package worker

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewDispatcher(t *testing.T) {
	t.Run("Retrieve object", func(t *testing.T) {
		newDispatcher := NewDispatcher(5)
		assert.Equal(t, 5, newDispatcher.MaxWorkers,
			"[TestNewDispatcher] Max workers not 5 in new dispatcher")
		assert.IsType(t, *new(chan chan JobInterface), newDispatcher.WorkerPool,
			"[TestNewDispatcher] Workerpool is not of type chan chan JobInterface")
	})
}

func TestDispatcher_Run(t *testing.T) {
	t.Run("Run with no workers", func(t *testing.T) {
		wasHitChecker := false
		testDispatcher := Dispatcher{
			WorkerPool:    make(chan chan JobInterface, 1),
			MaxWorkers:    0,
			NewWorkerFunc: func(workerPool chan chan JobInterface) WorkInterface {
				wasHitChecker = true
				return Worker{}
			},
		}
		testDispatcher.Run()
		assert.False(t, wasHitChecker)
	})

	t.Run("Run with single worker", func(t *testing.T) {
		worker := MockWorker{}
		testDispatcher := Dispatcher{
			WorkerPool:    make(chan chan JobInterface, 1),
			MaxWorkers:    1,
			NewWorkerFunc: func(workerPool chan chan JobInterface) WorkInterface {
				worker.ReturnChannel = make(chan bool, 1)
				return &worker
			},
		}
		testDispatcher.Run()
		assert.True(t, <-worker.ReturnChannel)
	})

	t.Run("Run with multiple worker", func(t *testing.T) {
		// Retain all return channels
		var returnChans []*chan bool
		testDispatcher := Dispatcher{
			WorkerPool:    make(chan chan JobInterface, 1),
			MaxWorkers:    10,
			NewWorkerFunc: func(workerPool chan chan JobInterface) WorkInterface {
				newWorker := MockWorker{}
				newWorker.ReturnChannel = make(chan bool, 1)
				returnChans = append(returnChans, &newWorker.ReturnChannel)
				return &newWorker
			},
		}
		testDispatcher.Run()

		// Loop through return channels
		for _, channel := range returnChans {
			// Blocks and Waits for True value to return
			assert.True(t, <-(*channel))
		}
	})
}