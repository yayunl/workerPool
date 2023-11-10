package main

import (
	"concurrentPool/workerPool"
	"fmt"
	"time"
)

var (
	taskFnc = func(args []any) (interface{}, error) {
		//fmt.Printf("Task 1 args %v\n", args)
		taskId := args[0]
		delay := args[1]
		arg2 := args[2]
		time.Sleep(time.Duration(delay.(int)) * time.Millisecond)
		result := fmt.Sprintf("task %d slept for %d ms. Its 2nd arg is %d.\n", taskId, delay, arg2)
		return result, nil
	}
)

func main() {
	workerNum, taskNum := 10, 500
	// Create a new workerPool with the specified number of workers
	wp := workerPool.New(workerNum)
	// Start the workerPool
	wp.Run()
	// Create a collection of tasks
	for j := 0; j < taskNum; j++ {
		i := j
		newJob := workerPool.NewJob(i, taskFnc, []any{i, i * 10, i + 2})
		wp.Add(newJob)
	}
	// Wait for all the tasks to be completed. Make sure to call Wait() before calling GetResult()
	wp.Wait()
	// Gather the results
	results := wp.GetResult()
	resultCnt := 0
	for result := range results {
		fmt.Printf("Task ID: %d, WorkerID: %d, Result: %v, Error: %v\n", result.TaskID, result.WorkerID, result.Value, result.Err)
		if resultCnt > 30 {
			break
		}
		resultCnt++
	}
	// Stop the dispatcher and all the workers
	wp.Stop()
	time.Sleep(3 * time.Second)
}
