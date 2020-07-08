package main

import (
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// DownloadFileTask holds information for downloading a file
type DownloadFileTask struct {
	url             string
	filenamePattern string
	dirPattern      string

	// context for filename pattern
	pageTitle string
	imgAlt    string
	linkText  string

	filter interface{} // TODO (@imafish) placeholder here. implement filter feature later.
}

// Execute downloads the file from url and save it to disk
func (t DownloadFileTask) Execute(ctx *Context) error {
	dirFormatted, err := formatString(t.dirPattern, t, ctx)
	if err != nil {
		ctx.log.Errorf("Failed to format dirPattern, err: %s", err.Error())
		return err
	}
	filenameFormatted, err := formatString(t.filenamePattern, t, ctx)
	if err != nil {
		ctx.log.Errorf("Failed to format filenamePattern, err: %s", err.Error())
		return err
	}

	dir := makeAbs(ctx.outDir, dirFormatted)
	os.MkdirAll(dir, os.ModePerm)
	fullPath := filepath.Join(dirFormatted, filenameFormatted)
	ctx.log.Debugf("got full path: %s", fullPath)

	err = downloadFile(t.url, fullPath)
	if err != nil {
		ctx.log.Errorf("Failed to download file from %s. err: %s", t.url, err.Error())
		return err
	}

	ctx.log.Infof("Downloaded from %s to %s", t.url, fullPath)
	return nil
}

func formatString(pattern string, t DownloadFileTask, ctx *Context) (string, error) {
	pattern = strings.ReplaceAll(pattern, "{pageTitle}", t.pageTitle)
	pattern = strings.ReplaceAll(pattern, "{imgAlt}", t.imgAlt)
	pattern = strings.ReplaceAll(pattern, "{linkText}", t.linkText)

	// file extension
	url, err := url.Parse(t.url)
	if err != nil {
		ctx.log.Errorf("Cannot parse url in download file task, err: %s", err.Error())
		return "", err
	}
	path := url.EscapedPath()
	ext := filepath.Ext(path)
	pattern = strings.ReplaceAll(pattern, "{.ext}", ext)

	// counter
	counterPattern := regexp.MustCompile(`\{_*i\}`)
	if counterPattern.MatchString(pattern) {
		i := <-ctx.counter
		iString := strconv.Itoa(i)
		iLength := len(iString)

		matches := counterPattern.FindAllString(pattern, -1)
		for _, match := range matches {
			matchLength := len(match) - 2
			replacement := iString
			if matchLength >= iLength {
				replacement = strings.Repeat("0", matchLength-iLength) + replacement
			}
			pattern = strings.ReplaceAll(pattern, match, replacement)
		}
	}

	return pattern, nil
}
