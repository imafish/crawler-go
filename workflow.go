package main

// url: the url to start crawl
// outDir: absolute directory as base dir for output
// concurrent: number of goroutines to handle crawling
// log: logger
func workflow(url string, outDir string, concurrent int, log Logger) {

	// initialize context
	context := &Context{
		taskChan: make(chan Task, 4096),
		quitChan: make(chan bool, concurrent),
		counter:  make(chan int),

		goroutineCount: concurrent,
		outDir:         outDir,
		log:            log,
		linkMap:        make(map[string]bool),
		pendingCount:   0,
	}
	context.wait.Add(concurrent)

	go func() {
		var count = 0
		for {
			count++
			if count < 0 {
				count = 0
			}
			context.counter <- count
		}
	}()

	context.log.Debugf("starting %d goroutines", concurrent)
	for i := 0; i < concurrent; i++ {
		go ruleExecutionFunc(context)
	}

	task := PageParseTask{
		url:   url,
		final: false,

		linkText: "1",
	}
	context.log.Debugf("adding initial tasks to queue")
	context.addTask(task)

	context.log.Debugf("waiting goroutines to complete")
	context.wait.Wait()
}

func ruleExecutionFunc(ctx *Context) {
	defer ctx.wait.Done()

	for {
		select {
		case task := <-ctx.taskChan:
			err := task.Execute(ctx)
			ctx.finishExecuting()
			if err != nil {
				ctx.log.Errorf("Failed to execute task, err: %s", err.Error())
			}

		case <-ctx.quitChan:
			ctx.log.Info("worker goroutine quited.")
			return
		}
	}
}
