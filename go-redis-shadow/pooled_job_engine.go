package main

import (
	"errors"
	"log"
	"sync/atomic"
	"time"
	"unsafe"

	cmap "github.com/orcaman/concurrent-map"
	uuid "github.com/satori/go.uuid"
)

// PooledJobEngine Unstable!!!
type PooledJobEngine struct {
	info            jobEngineInfo      // store relating configuration
	taskQueue       chan *task         // store all task need to execute, but no worker available when it comes in
	workerPool      cmap.ConcurrentMap // store all worker, id -> *worker
	idleWorkerQueue chan *worker       // once a worker finished its task, it will be changed into idle state and be add to this queue
	idleWorkerTimer *Timer             /* this timer will be activated to clean redundant idle worker
	when idleWorkerQueue's size is bigger than info.corePoolSize */
	running    bool
	stopSignal chan bool
}

type jobEngineInfo struct {
	taskQueueSize int
	corePoolSize  int
	maxPoolSize   int
	keepAliveTime time.Duration
}

// NewPooledJobEngine ...
func NewPooledJobEngine(taskQueueSize, corePoolSize, maxPoolSize int, keepAliveTime time.Duration) *PooledJobEngine {
	j := &PooledJobEngine{
		info: jobEngineInfo{
			taskQueueSize: taskQueueSize,
			corePoolSize:  corePoolSize,
			maxPoolSize:   maxPoolSize,
			keepAliveTime: keepAliveTime,
		},
		taskQueue:       make(chan *task, taskQueueSize),
		workerPool:      cmap.New(),
		idleWorkerQueue: make(chan *worker, maxPoolSize),
		idleWorkerTimer: nil,
		running:         true,
		stopSignal:      make(chan bool, 1),
	}
	j.idleWorkerTimer = NewTimer(keepAliveTime, j.cleanIdleWorker)
	// initialize workerPool
	for i := 0; i < corePoolSize; i++ {
		j.addWorker()
	}
	return j
}

// create new worker, new worker will be add to idle queue automatically, with in function worker.start
func (j *PooledJobEngine) addWorker() *worker {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	id := uuid.String()
	w := newWorker(id, j.onWorkerStateChange)
	w.start()

	j.workerPool.Set(id, w)

	return w
}

func (j *PooledJobEngine) removeWorker(w *worker) {
	w.stop()
}

// worker state listener, when worker state changed, this function will be invoked
func (j *PooledJobEngine) onWorkerStateChange(id string, state int32) {
	switch state {
	case eStateIdle:
		if tmp, ok := j.workerPool.Get(id); ok {
			w := tmp.(*worker)
			if len(j.taskQueue) > 0 {
				w.execute(<-j.taskQueue)
			} else {
				j.idleWorkerQueue <- w
				if len(j.idleWorkerQueue) > j.info.corePoolSize {
					// if timeout, invoke cleanIdleWorker to clean redundant worker
					j.idleWorkerTimer.Start()
				}
			}
		} else {
			log.Fatal("no worker found with id = " + id)
		}
	case eStateStopped:
		j.workerPool.Remove(id)
		if j.workerPool.IsEmpty() {
			j.stopSignal <- true
		}
	default: // pass
	}
}

func (j *PooledJobEngine) cleanIdleWorker() {
	var w *worker
	for len(j.idleWorkerQueue) > j.info.corePoolSize {
		w = <-j.idleWorkerQueue
		j.removeWorker(w)
	}
}

// Execute ... TODO: return value
func (j *PooledJobEngine) Execute(callable func(...interface{}), args ...interface{}) (err error) {
	if callable == nil {
		err = errors.New("Invalid argument, callable must be non-nil")
	} else if !j.running {
		err = errors.New("PooledJobEngine has been stopped")
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
		if j.workerPool.Count() < j.info.maxPoolSize {
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
func (j *PooledJobEngine) Shutdown() {
	j.running = false
	close(j.taskQueue)
	for {
		select {
		case w := <-j.idleWorkerQueue:
			w.stop()
		case <-j.stopSignal:
			goto end
		}
	}
end:
	close(j.idleWorkerQueue)
	j.workerPool = nil
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
	if t == nil {
		return errors.New("Invalid argument, task must be non-nil")
	} else if w.state != workerState(eStateIdle) {
		return errors.New("Invalid state " + w.state.toString() + ", worker must be in idle state")
	}
	w.taskEntry <- t
	return nil
}

func (w *worker) stop() {
	close(w.taskEntry)
}
