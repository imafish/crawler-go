package main

import (
	"os"
	"path/filepath"
)

// DownloadFileTask holds information for downloading a file
type DownloadFileTask struct {
	url             string
	filenamePattern string
	dirPattern      string

	taskContext *TaskContext

	filters []interface{} // TODO (@imafish) placeholder here. implement filter feature later.
}

// URL returns the task's corresponding url
func (t DownloadFileTask) URL() string {
	return t.url
}

// Execute downloads the file from url and save it to disk
func (t DownloadFileTask) Execute(ctx *ExecutionContext) error {
	dirFormatted := FormatString(t.dirPattern, t.taskContext)
	filenameFormatted := FormatString(t.filenamePattern, t.taskContext)

	groupContext := ctx.FindOrNewGroupContext(t.taskContext, false, nil)
	t.taskContext.counter = groupContext.counter
	dirFormatted = filepath.Join(makeAbs(ctx.baseDir, groupContext.dir), dirFormatted)

	os.MkdirAll(dirFormatted, os.ModePerm)
	fullPath := filepath.Join(dirFormatted, filenameFormatted)
	GetLogger().Debugf("got full path: %s", fullPath)

	err := downloadFile(t.url, fullPath)
	if err != nil {
		GetLogger().Errorf("Failed to download file from %s. err: %s", t.url, err.Error())
		return err
	}

	GetLogger().Infof("Downloaded from %s to %s", t.url, fullPath)
	return nil
}
