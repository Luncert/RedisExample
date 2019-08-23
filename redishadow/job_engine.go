package main

import (
	"errors"
)

// JobEngine ...
type JobEngine struct {
	concurrent int
	taskQueue  chan task
	stopSignal chan bool
	running    bool
}

type task struct {
	callable func(...interface{})
	args     []interface{}
}

// NewJobEngine ...
func NewJobEngine(taskQueueSize, concurrentLevel int) *JobEngine {
	j := &JobEngine{
		concurrent: concurrentLevel,
		taskQueue:  make(chan task, taskQueueSize),
		stopSignal: make(chan bool),
		running:    true,
	}
	j.init()
	return j
}

// create dispatcher
func (j *JobEngine) init() {
	go func() {
		running := 0
		wait := make(chan bool, j.concurrent)
		waitOneFinish := func() {
			<-wait
			running--
		}

		for task := range j.taskQueue {
			go j.wrappedCall(task, wait)
			running++
			// if the number of concurrencies reaches ex.concurrent
			// then wait a task to finish
			if running == j.concurrent {
				waitOneFinish()
			}
		}

		// if there are still some task running
		// wait them
		for running > 0 {
			waitOneFinish()
		}
		close(wait)

		j.stopSignal <- true
		j.running = false
	}()
}

func (j *JobEngine) wrappedCall(t task, signal chan<- bool) {
	defer func() {
		signal <- true
	}()
	t.callable(t.args...)
}

// Execute ...
func (j *JobEngine) Execute(callable func(...interface{}), args ...interface{}) (err error) {
	if !j.running {
		err = errors.New("JobEngine has been stopped")
	} else {
		j.taskQueue <- task{callable: callable, args: args}
	}
	return
}

// Shutdown will block until no task in task queue or being executing,
// after this call, this JobEngine is broken
func (j *JobEngine) Shutdown() {
	close(j.taskQueue)
	<-j.stopSignal
}
