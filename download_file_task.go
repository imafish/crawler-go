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
	nc *namingContext

	filter interface{} // TODO (@imafish) placeholder here. implement filter feature later.
}

// URL returns the task's corresponding url
func (t DownloadFileTask) URL() string {
	return t.url
}

// Execute downloads the file from url and save it to disk
func (t DownloadFileTask) Execute(ctx *Context) error {
	dirFormatted, err := formatString(t.dirPattern, t.nc, ctx)
	if err != nil {
		ctx.log.Errorf("Failed to format dirPattern, err: %s", err.Error())
		return err
	}
	filenameFormatted, err := formatString(t.filenamePattern, t.nc, ctx)
	if err != nil {
		ctx.log.Errorf("Failed to format filenamePattern, err: %s", err.Error())
		return err
	}

	dirFormatted = makeAbs(ctx.outDir, dirFormatted)
	os.MkdirAll(dirFormatted, os.ModePerm)
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

func formatString(pattern string, nc *namingContext, ctx *Context) (string, error) {
	pattern = strings.ReplaceAll(pattern, "{pageTitle}", nc.pageTitle)
	pattern = strings.ReplaceAll(pattern, "{imgAlt}", nc.imgAlt)
	pattern = strings.ReplaceAll(pattern, "{linkText}", nc.text)
	pattern = strings.ReplaceAll(pattern, "{.ext}", nc.ext)

	// counter
	counterPattern := regexp.MustCompile(`\{_*i\}`)
	if counterPattern.MatchString(pattern) {
		i := <-nc.i.i
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

func getExtension(urlString string) (string, error) {
	// file extension
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	path := u.EscapedPath()
	ext := filepath.Ext(path)

	return ext, nil
}
