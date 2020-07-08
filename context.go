package main

import (
	"sync"
)

// Context is execution context for this crawler
type Context struct {
	taskChan chan Task
	quitChan chan bool
	wait     sync.WaitGroup
	counter  chan int

	// application configurations:
	goroutineCount int
	outDir         string
	log            Logger

	mtx          sync.Mutex // mutex to protect task control data. protects the following variable
	linkMap      map[string]bool
	pendingCount int
}

func (ctx *Context) startExecuting(url string) bool {
	ctx.mtx.Lock()
	defer ctx.mtx.Unlock()

	_, ok := ctx.linkMap[url]
	if ok {
		return false
	}

	ctx.linkMap[url] = false
	return true
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
	ctx.taskChan <- task
	ctx.log.Debugf("added task to queue: %#v", task)

	ctx.mtx.Lock()
	defer ctx.mtx.Unlock()
	ctx.pendingCount++
}
