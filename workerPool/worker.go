package workerPool

import (
	"fmt"
)

type (
	// result is the result of a task
	result struct {
		// WorkerID is the ID of the worker that executed the task
		WorkerID int
		// TaskID is the ID of the task
		TaskID interface{}
		// Value is the result of the task
		Value interface{}
		// Err is the error of the task
		Err error
	}

	// worker type defines the data struct of a worker
	worker struct {
		Id         int
		dispatcher *Dispatcher // each worker has a pointer to the dispatcher
		jobChan    JobChannel  // each worker has its own job channel.
		done       chan struct{}
	}

	JobChannel    chan *Job
	JobQueue      chan chan *Job
	ResultChannel chan *result
)

func newWorker(dispatcher *Dispatcher, workerID int, jobChan JobChannel, done chan struct{}) *worker {
	return &worker{
		Id:         workerID,
		dispatcher: dispatcher,
		jobChan:    jobChan,
		done:       done,
	}
}

// Start starts the worker
func (wr *worker) Start() {
	go func() {
		defer fmt.Printf("worker %d done\n", wr.Id)
		defer close(wr.jobChan)
		for {
			select {
			case <-wr.done:
				return
			case wr.dispatcher.Queue <- wr.jobChan: // Add the worker's job channel to the queue of the workerPool to let it know the worker is available
				select {
				case <-wr.done:
					return
				// Wait for a job to be assigned by the workerPool through the worker's job channel
				case t := <-wr.jobChan:
					// Process the task
					val, err := t.execute()
					res := &result{
						TaskID:   t.ID,
						Value:    val,
						Err:      err,
						WorkerID: wr.Id,
					}
					// Create a new task
					if t.Args[2].(bool) {
						newJobId := t.ID * -1
						newDelayTime := t.Args[1].(int) * 10
						newJob := NewJob(
							newJobId,
							t.Fn,
							[]any{newJobId, newDelayTime, false})

						wr.dispatcher.WorkChan <- newJob
						wr.dispatcher.Wg.Add(1)
					}

					// Send the result to the workerPool's result channel
				loop:
					for {
						select {
						case wr.dispatcher.ResultChan <- res:
							wr.dispatcher.Wg.Done()
							break loop
						case <-wr.done:
							return
						}
					}

				}
			}

		}
	}()
}

// Stop closes the Done channel on the worker, causing the goroutine to exit.
func (wr *worker) Stop() {
	close(wr.done)
}
