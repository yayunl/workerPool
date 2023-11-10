package workerPool

import (
	"fmt"
	"sync"
)

type Dispatcher struct {
	Workers    []*worker  // this is the list of workers that workerPool tracks
	WorkChan   JobChannel // client submits a job to this channel
	Queue      JobQueue   // this is the shared JobPool between the workers
	ResultChan ResultChannel
	Done       chan struct{}
	Wg         sync.WaitGroup // this is used to wait for all jobs to be done
}

// New creates a new workerPool
func New(numWorkers int) *Dispatcher {
	dispatcher := &Dispatcher{
		Workers:    make([]*worker, numWorkers),
		WorkChan:   make(JobChannel),
		Queue:      make(JobQueue),
		ResultChan: make(ResultChannel),
		Done:       make(chan struct{}),
	}
	// create workers
	for i := 0; i < numWorkers; i++ {
		// create a worker. d.Queue and d.ResultChan are shared between all workers and dispatchers.
		wrk := newWorker(dispatcher, i, make(JobChannel), make(chan struct{}))
		dispatcher.Workers[i] = wrk
	}
	return dispatcher
}

func (d *Dispatcher) Run() {

	// start workers
	for _, wrk := range d.Workers {
		fmt.Printf("Dispatcher: creating worker %d\n", wrk.Id)
		wrk.Start()
	}

	// start the workerPool by creating a goroutine
	go func() {
		defer fmt.Printf("Dispatcher done\n")
		for {
			select {
			case task, ok := <-d.WorkChan: // listen to a submitted job on WorkChannel
				if !ok { // if the WorkChannel is closed, return
					return
				}
				jobChanOfIdleWorker := <-d.Queue // pull out an available job channel from queue. The job channel indicates which worker is available.
				jobChanOfIdleWorker <- task      // submit the job on the available job channel and the worker of the job channel will pick it up.

			case <-d.Done:
				return
			}
		}
	}()
}

// Add adds a job to the workerPool
func (d *Dispatcher) Add(job *Job) {
	d.Wg.Add(1)
	fmt.Printf("Dispatcher: adding job %d\n", job.ID)
	go func() { d.WorkChan <- job }()
}

func (d *Dispatcher) Wait() {
	go func() {
		d.Wg.Wait()
		fmt.Printf("Dispatcher: all jobs done\n")
		close(d.ResultChan)
	}()
}

func (d *Dispatcher) Stop() {
	for _, wrk := range d.Workers {
		fmt.Printf("Dispatcher: stopping worker %d\n", wrk.Id)
		wrk.Stop()
	}
	close(d.Done) // Close the Done channel after all workers are stopped
}

// GetResult gets the result from the result channel
func (d *Dispatcher) GetResult() <-chan *result {
	return d.ResultChan
}
