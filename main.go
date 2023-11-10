package main

import (
	"concurrentPool/workerPool"
	"fmt"
	"time"
)

var (
	taskFnc = func(args []any) (interface{}, error) {
		var result string
		taskId := args[0].(int)
		delay := args[1].(int)

		time.Sleep(time.Duration(delay) * time.Millisecond)
		if args[2] == true {
			result = fmt.Sprintf("task %d slept for %d ms. create a sub-task %d. \n", taskId, delay, taskId*-1)
		} else {
			result = fmt.Sprintf("task %d slept for %d ms. \n", taskId, delay)
		}

		return result, nil
	}
)

func main() {
	workerNum, taskNum := 3, 5
	// Create a new workerPool with the specified number of workers
	wp := workerPool.New(workerNum)
	// Start the workerPool
	wp.Run()
	// Create a collection of tasks
	for j := 1; j <= taskNum; j++ {
		i := j
		if i%2 == 0 { // create a sub-task for task with even id
			newJob := workerPool.NewJob(i, taskFnc, []any{i, i * 10, true})
			wp.Add(newJob)
		} else {
			newJob := workerPool.NewJob(i, taskFnc, []any{i, i * 10, false})
			wp.Add(newJob)
		}

	}
	// Wait for all the tasks to be completed. Make sure to call Wait() before calling GetResult()
	wp.Wait()
	// Gather the results
	for result := range wp.GetResult() {
		fmt.Printf("Task ID: %d, WorkerID: %d, Result: %v, Error: %v\n", result.TaskID, result.WorkerID, result.Value, result.Err)
	}
	// Stop the dispatcher and all the workers
	wp.Stop()
	time.Sleep(3 * time.Second)
}
