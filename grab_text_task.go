package main

import (
	"fmt"
	"os"
	"path/filepath"

	iconv "github.com/djimenez/iconv-go"
	"golang.org/x/net/html"
)

// GrabTextTask grabs tasks from page content and save to file
type GrabTextTask struct {
	taskContext     *TaskContext
	node            *html.Node
	dirPattern      string
	filenamePattern string
	target          target
	url             string
	converter       *iconv.Converter
}

// URL returns task's unique URL
func (t GrabTextTask) URL() string {
	return fmt.Sprintf("%s %v", t.url, t.target)
}

// Execute runs this task.
func (t GrabTextTask) Execute(ctx *ExecutionContext) error {
	dirFormatted := FormatString(t.dirPattern, t.taskContext)
	filenameFormatted := FormatString(t.filenamePattern, t.taskContext)

	groupContext := ctx.FindOrNewGroupContext(t.taskContext, false, nil)
	t.taskContext.counter = groupContext.counter
	dirFormatted = filepath.Join(makeAbs(ctx.baseDir, groupContext.dir), dirFormatted)

	os.MkdirAll(dirFormatted, os.ModePerm)
	fullPath := filepath.Join(dirFormatted, filenameFormatted)
	GetLogger().Debugf("got full path: %s", fullPath)

	targetNode, err := GetMatchingNodes(t.node, t.target)
	if err != nil {
		GetLogger().Errorf("Error getting matching target node, error: %s", err.Error())
		return err
	}
	if len(targetNode) != 1 {
		err = fmt.Errorf("couldn't find exact one matching inner node for task %s", t.URL())
		GetLogger().Error(err.Error())
		return err
	}

	text := CollectText(targetNode[0], t.converter)
	f, err := os.Create(fullPath)
	if err != nil {
		GetLogger().Errorf("Error creating text file %s, error: %s", fullPath, err.Error())
		return err
	}
	defer f.Close()
	f.WriteString(text)
	f.WriteString("\n")

	return nil
}
