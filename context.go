package main

import (
	"sync"
)

// Context is execution context for this crawler
type Context struct {
	taskChan chan Task
	quitChan chan bool
	wait     sync.WaitGroup

	// application configurations:
	goroutineCount int
	outDir         string
	config         *configuration
	log            Logger

	mtx          sync.Mutex // mutex to protect task control data. protects the following variable
	linkMap      map[string]bool
	pendingCount int
}

func (ctx *Context) finishExecuting() {
	ctx.mtx.Lock()
	defer ctx.mtx.Unlock()

	ctx.pendingCount--
	ctx.log.Debugf("task finished execution, remaining task count: %d", ctx.pendingCount)

	if ctx.pendingCount == 0 {
		ctx.log.Info("no more tasks left, quitting...")
		for i := 0; i < ctx.goroutineCount; i++ {
			ctx.quitChan <- true
		}
	}
}

func (ctx *Context) addTask(task Task) {
	u := task.URL()

	ctx.mtx.Lock()
	defer ctx.mtx.Unlock()
	if _, ok := ctx.linkMap[u]; !ok {
		// task does not exist.
		ctx.linkMap[u] = false
		ctx.taskChan <- task
		ctx.pendingCount++
		ctx.log.Debugf("added task to queue: %#v", task)
	}
}
