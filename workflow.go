package main

// rulePath: path to rule file (.yaml)
// outDir: absolute directory as base dir for output
// concurrent: number of goroutines to handle crawling
// log: logger
func workflow(rulePath string, outDir string, concurrent int) {

	config, err := createConfigurationFromYaml(rulePath)
	if err != nil {
		GetLogger().Errorf("Failed to load rules from file %s. err: %s", rulePath, err.Error())
		return
	}

	// initialize context
	context := &ExecutionContext{
		taskChan:       make(chan Task, 2048*concurrent),
		quitChan:       make(chan bool, concurrent),
		baseDir:        outDir,
		config:         config,
		goroutineCount: concurrent,
		groups:         make([]*GroupContext, 0),
		linkMap:        make(map[string]bool),
		pendingCount:   0,
		defaultGroup: &GroupContext{
			counter: newCounter(),
		},
	}
	GetLogger().Debugf("context: %#v", context)

	// first look for targets with URL, they'll act as starter for the crawler
	gotTask := false
	for _, s := range config.Startpages {
		task := PageParseTask{
			url:         s.URL,
			final:       false,
			newGroup:    true,
			groupFormat: &s.Group,
			taskContext: &TaskContext{},
		}
		context.AddTask(task)
		gotTask = true
	}

	if !gotTask {
		GetLogger().Warning("No startpage entry is found in the config file.")
		return
	}

	GetLogger().Debugf("starting %d goroutines", concurrent)
	for i := 0; i < concurrent; i++ {
		go workerFunc(context)
	}

	context.wait.Add(concurrent)
	GetLogger().Debugf("waiting %d goroutines to complete", concurrent)
	context.wait.Wait()
}

func workerFunc(ctx *ExecutionContext) {
	defer ctx.wait.Done()

	for {
		select {
		case task := <-ctx.taskChan:
			ctx.ExecuteTask(task)

		case <-ctx.quitChan:
			GetLogger().Info("worker goroutine quited.")
			return
		}
	}
}
