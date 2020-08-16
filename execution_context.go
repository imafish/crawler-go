package main

import (
	"path/filepath"
	"regexp"
	"sync"
)

// ExecutionContext is execution context for this crawler
type ExecutionContext struct {
	taskChan chan Task
	quitChan chan bool
	wait     sync.WaitGroup

	// global config
	goroutineCount int
	config         *configuration
	baseDir        string

	// runtime context
	groups       []*GroupContext
	defaultGroup *GroupContext

	mtx          sync.Mutex      // Protects task control data. Protects the following variable
	linkMap      map[string]bool // Keeps track of links that has been processed.
	pendingCount int             // Tasks left in the channels
}

func (ctx *ExecutionContext) finishExecuting() {
	ctx.mtx.Lock()
	defer ctx.mtx.Unlock()

	ctx.pendingCount--
	GetLogger().Debugf("task finished execution, remaining task count: %d", ctx.pendingCount)

	if ctx.pendingCount == 0 {
		GetLogger().Info("no more tasks left, quitting...")
		for i := 0; i < ctx.goroutineCount; i++ {
			ctx.quitChan <- true
		}
	}
}

// AddTask adds the task to execution channel if it hasn't been processed already.
func (ctx *ExecutionContext) AddTask(task Task) {
	u := task.URL()

	ctx.mtx.Lock()
	defer ctx.mtx.Unlock()
	if _, ok := ctx.linkMap[u]; !ok {
		// task does not exist.
		ctx.linkMap[u] = false
		ctx.taskChan <- task
		ctx.pendingCount++
		GetLogger().Debugf("added task to queue, pending tasks: %d, task: %#v", ctx.pendingCount, task)
	} else {
		GetLogger().Debugf("Task already exists.")
	}
}

// ExecuteTask runs the given task, changes execution context on task finish.
func (ctx *ExecutionContext) ExecuteTask(task Task) {
	defer ctx.finishExecuting()
	err := task.Execute(ctx)
	if err != nil {
		GetLogger().Debugf("Failed to execute task %#v, error is: %s", task, err.Error())
	}
}

// FindOrNewGroupContext tries to find group context in execution context
// If not found, create a new one or return nil base on createNew and groupFormat parameter
func (ctx *ExecutionContext) FindOrNewGroupContext(taskContext *TaskContext, createNew bool, groupFormat *group) *GroupContext {
	var group *GroupContext
	for _, g := range ctx.groups {
		if g.IsMatch(taskContext) {
			group = g
			break
		}
	}

	if group != nil && createNew {
		GetLogger().Warning("createNew is true but found a matching group. Not creating new group")
	}

	if group == nil && createNew {
		group = &GroupContext{
			groupBy: groupFormat.GroupBy,
			name:    FormatString(groupFormat.GroupBy, taskContext),
			counter: newCounter(),
			dir:     filepath.Join(ctx.baseDir, FormatString(groupFormat.DirPattern, taskContext)),
		}
	}

	if group == nil {
		group = ctx.defaultGroup
	}

	return group
}

// FindMatchingRules finds rules that matches page's title and URL from execution context.
func (ctx *ExecutionContext) FindMatchingRules(taskContext *TaskContext) []Rule {
	matchingRules := make([]Rule, 0)

	for _, r := range ctx.config.Rules {
		for _, m := range r.Matches {
			isMatch := false
			if m.URL != "" {
				regex, err := regexp.Compile(m.URL)
				if err != nil {
					GetLogger().Warningf("Cannot compile regex from string %s", m.URL)
				} else if regex.MatchString(taskContext.pageURL) {
					isMatch = true
				}
			}
			if m.Title != "" {
				regex, err := regexp.Compile(m.Title)
				if err != nil {
					GetLogger().Warningf("Cannot compile regex from string %s", m.Title)
				} else if regex.MatchString(taskContext.pageTitle) {
					isMatch = true
				}
			}

			if isMatch {
				matchingRules = append(matchingRules, r)
				break
			}
		}
	}

	return matchingRules
}
