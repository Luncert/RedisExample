package main

import (
	"errors"
	"github.com/satori/go.uuid"
	"log"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
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

// JobEngineX ...
type JobEngineX struct {
	info            jobEngineInfo // store relating configuration
	mutex           sync.Mutex
	taskQueue       chan *task         // store all task need to execute, but no worker available when it comes in
	workerPool      map[string]*worker // store all worker
	idleWorkerQueue chan *worker       // once a worker finished its task, it will be changed into idle state and be add to this queue
	idleWorkerTimer *Timer             /* this timer will be activated to clean redundant idle worker
	when idleWorkerQueue's size is bigger than info.corePoolSize */
	running bool
}

type jobEngineInfo struct {
	taskQueueSize int
	corePoolSize  int
	maxPoolSize   int
	keepAliveTime time.Duration
}

// NewJobEngineX ...
func NewJobEngineX(taskQueueSize, corePoolSize, maxPoolSize int, keepAliveTime time.Duration) *JobEngineX {
	j := &JobEngineX{
		info: jobEngineInfo{
			taskQueueSize: taskQueueSize,
			corePoolSize:  corePoolSize,
			maxPoolSize:   maxPoolSize,
			keepAliveTime: keepAliveTime,
		},
		mutex:           sync.Mutex{},
		taskQueue:       make(chan *task, taskQueueSize),
		workerPool:      map[string]*worker{},
		idleWorkerQueue: make(chan *worker, maxPoolSize),
		idleWorkerTimer: nil,
		running:         true,
	}
	j.idleWorkerTimer = NewTimer(keepAliveTime, j.cleanIdleWorker)
	// initialize workerPool
	for i := 0; i < corePoolSize; i++ {
		j.addWorker()
	}
	return j
}

// create new worker, new worker will be add to idle queue automatically, with in function worker.start
func (j *JobEngineX) addWorker() *worker {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	id := uuid.String()
	w := newWorker(id, j.onWorkerStateChange)
	w.start()

	j.mutex.Lock()
	j.workerPool[id] = w
	j.mutex.Unlock()

	return w
}

func (j *JobEngineX) removeWorker(w *worker) {
	w.stop()

	j.mutex.Lock()
	delete(j.workerPool, w.id)
	j.mutex.Unlock()
}

// worker state listener, when worker state changed, this function will be invoked
func (j *JobEngineX) onWorkerStateChange(id string, state int32) {
	switch state {
	case eStateIdle:
		w, _ := j.workerPool[id]
		if len(j.taskQueue) > 0 {
			w.execute(<-j.taskQueue)
		} else {
			j.idleWorkerQueue <- w
			if len(j.idleWorkerQueue) > j.info.corePoolSize {
				// if timeout, invoke cleanIdleWorker to clean redundant worker
				j.idleWorkerTimer.Start()
			}
		}
	default: // pass
	}
}

func (j *JobEngineX) cleanIdleWorker() {
	var w *worker
	for len(j.idleWorkerQueue) > j.info.corePoolSize {
		w = <-j.idleWorkerQueue
		j.removeWorker(w)
	}
}

// Execute ... TODO: return value
func (j *JobEngineX) Execute(callable func(...interface{}), args ...interface{}) (err error) {
	if callable == nil {
		err = errors.New("Invalid argument, callable must be non-nil")
	} else if !j.running {
		err = errors.New("JobEngineX has been stopped")
	} else {
		var w *worker
		var task = &task{callable: callable, args: args}
		// if there is worker in idle state, get it from idleWorkerQueue to run task directly
		for len(j.idleWorkerQueue) > 0 {
			w = <-j.idleWorkerQueue
			if !w.state.equals(eStateStopped) {
				err = w.execute(task)
				return
			}
		}
		// when step into this branch, it means all workers are running,
		// then try to create new worker if our worker pool haven't touch the upper limit maxPoolSize,
		// or else just add new task to taskQueue
		if len(j.workerPool) < j.info.maxPoolSize {
			j.addWorker()
			w = <-j.idleWorkerQueue // only idle worker could execute new task
			err = w.execute(task)
		} else {
			j.taskQueue <- task
		}
	}
	return
}

// Shutdown ...
func (j *JobEngineX) Shutdown() {
	close(j.taskQueue)
	stoppedWorker := 0
	workerNum := len(j.workerPool)
	for w := range j.idleWorkerQueue {
		w.stop()
		stoppedWorker++
		if stoppedWorker == workerNum {
			// all worker have been stopped
			j.workerPool = nil
			close(j.idleWorkerQueue)
			break
		}
	}
	j.running = false
}

// worker states

type workerState int32

const (
	eStateReady int32 = iota
	eStateRunning
	eStateIdle
	eStateStopped
)

func (ws workerState) equals(s int32) bool {
	return int32(ws) == s
}

func (ws workerState) toString() string {
	switch int32(ws) {
	case eStateReady:
		return "Ready"
	case eStateRunning:
		return "Running"
	case eStateIdle:
		return "Idle"
	case eStateStopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

// worker

type worker struct {
	id                  string
	taskEntry           chan *task
	state               workerState
	stateChangeListener func(id string, state int32)
}

func newWorker(id string, onStateChange func(id string, state int32)) *worker {
	return &worker{
		id:                  id,
		taskEntry:           make(chan *task, 1),
		state:               workerState(eStateReady),
		stateChangeListener: onStateChange,
	}
}

func (w *worker) start() {
	go func() {
		w.updateState(eStateIdle)
		for t := range w.taskEntry {
			w.updateState(eStateRunning)
			if t == nil {
				log.Println("nil task")
			}
			t.callable(t.args...)
			w.updateState(eStateIdle)
		}
		w.updateState(eStateStopped)
	}()
}

func (w *worker) updateState(state int32) {
	pState := (*int32)(unsafe.Pointer(&w.state))
	atomic.CompareAndSwapInt32(pState, int32(w.state), state)
	w.stateChangeListener(w.id, state)
}

func (w *worker) execute(t *task) error {
	if w.state != workerState(eStateIdle) {
		return errors.New("Invalid state " + w.state.toString() + ", worker must be in idle state")
	}
	w.taskEntry <- t
	return nil
}

func (w *worker) stop() {
	close(w.taskEntry)
}
