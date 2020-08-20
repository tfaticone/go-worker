package worker

import "log"

type JobInterface interface {
	GetPayload() interface{}
	PanicOnFailure() bool
	Handler() error
}

type WorkInterface interface {
	Start()
}

// A buffered channel that we can send work requests on.
var JobQueue chan JobInterface

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool  chan chan JobInterface
	JobChannel  chan JobInterface
	quit    	chan bool
}

func NewWorker(workerPool chan chan JobInterface) WorkInterface {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan JobInterface),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				err := job.Handler()
				if err != nil && job.PanicOnFailure() {
					panic(err)
				} else if err != nil && !job.PanicOnFailure() {
					log.Panic("Error occurred during handle: ", err.Error())
				}
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
