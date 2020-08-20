package worker

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan JobInterface
	MaxWorkers int
	NewWorkerFunc func(workerPool chan chan JobInterface) WorkInterface
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan JobInterface, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		MaxWorkers: maxWorkers,
		NewWorkerFunc: NewWorker,
	}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := d.NewWorkerFunc(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// a job request has been received
			go func(job JobInterface) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
