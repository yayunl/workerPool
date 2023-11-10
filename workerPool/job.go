package workerPool

type jobFnc func(args []any) (interface{}, error)

// Job is a data structure that represents a job
type Job struct {
	ID int
	// Fn is the function to be executed
	Fn jobFnc
	// Args is the arguments of the Fn
	Args []any
}

// NewJob creates a new task
func NewJob(id int, fn jobFnc, args []any) *Job {
	return &Job{
		ID:   id,
		Fn:   fn,
		Args: args,
	}
}

// Execute lets the task be executed in a worker
func (t *Job) execute() (interface{}, error) {
	return t.Fn(t.Args)
}
