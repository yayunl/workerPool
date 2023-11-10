# WorkerPool

This repository contains the version 2 implementation of the worker pool pattern discussed in the blog post [GoLang Series 1.3.4: The Worker Pool Pattern](https://www.yayunliu90.blog/post/golang-series-1-3-4-the-worker-pool-pattern).
The implementation uses a dispatcher to distribute work to workers.

## Components
### The dispatcher
The dispatcher is responsible for creating the worker pool and distributing work to the workers. 

### Workers
The workers are responsible for executing the work and returning the result to the dispatcher.

### The client
The client interacts with the dispatcher to submit work and retrieve the results.


## Results
```bash
$ go run main.go
Dispatcher: creating worker 0
Dispatcher: creating worker 1
Dispatcher: creating worker 2
Dispatcher: adding job 1
Dispatcher: adding job 2
Dispatcher: adding job 3
Dispatcher: adding job 4
Dispatcher: adding job 5
Task ID: 1, WorkerID: 0, Result: task 1 slept for 10 ms. 
, Error: <nil>
Task ID: 3, WorkerID: 2, Result: task 3 slept for 30 ms. 
, Error: <nil>
Task ID: 2, WorkerID: 1, Result: task 2 slept for 20 ms. create a sub-task -2.
, Error: <nil>
Task ID: 4, WorkerID: 0, Result: task 4 slept for 40 ms. create a sub-task -4. 
, Error: <nil>
Task ID: 5, WorkerID: 2, Result: task 5 slept for 50 ms. 
, Error: <nil>
Task ID: -2, WorkerID: 1, Result: task -2 slept for 200 ms. 
, Error: <nil>
Dispatcher: all jobs done
Task ID: -4, WorkerID: 0, Result: task -4 slept for 400 ms. 
, Error: <nil>
Dispatcher: stopping worker 0
Dispatcher: stopping worker 1
Dispatcher: stopping worker 2
Dispatcher done
worker 1 done
worker 0 done
worker 2 done
```