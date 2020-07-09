package main

// rulePath: path to rule file (.yaml)
// outDir: absolute directory as base dir for output
// concurrent: number of goroutines to handle crawling
// log: logger
func workflow(rulePath string, outDir string, concurrent int, log Logger) {

	configuration, err := createConfigurationFromYaml(rulePath)
	if err != nil {
		log.Errorf("Failed to load rules from file %s. err: %s", rulePath, err.Error())
		return
	}

	// initialize context
	context := &Context{
		taskChan: make(chan Task, 2048*concurrent),
		quitChan: make(chan bool, concurrent),

		goroutineCount: concurrent,
		outDir:         outDir,
		config:         configuration,

		log:          log,
		linkMap:      make(map[string]bool),
		pendingCount: 0,
	}
	log.Debugf("context: %#v", context)

	// first look for targets with URL, they'll act as starter for the crawler
	gotTask := false
	for _, s := range configuration.Starters {
		for _, t := range s.Targets {
			start(t, s.StartGroup, context)
			gotTask = true
		}
	}

	if !gotTask {
		log.Warning("No task is found in the rule file.")
		return
	}

	context.log.Debugf("starting %d goroutines", concurrent)
	for i := 0; i < concurrent; i++ {
		go ruleExecutionFunc(context)
	}

	context.wait.Add(concurrent)
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

func start(url string, g *startGroupAction, ctx *Context) {
	ct := newCounter()

	task := PageParseTask{
		url:   url,
		final: false,

		gc: groupContext{
			i:          ct,
			dir:        g.DirPattern,
			name:       g.GroupName,
			firstParse: true,
		},
	}

	ctx.addTask(task)
}
