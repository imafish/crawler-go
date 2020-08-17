package main

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// FormatString formats pattern with given task context
func FormatString(pattern string, taskContext *TaskContext) string {
	pattern = strings.ReplaceAll(pattern, "{startpage.title}", taskContext.startPageTitle)
	pattern = strings.ReplaceAll(pattern, "{startpage.URL}", taskContext.startPageURL)
	pattern = strings.ReplaceAll(pattern, "{title}", taskContext.pageTitle)
	pattern = strings.ReplaceAll(pattern, "{URL}", taskContext.pageURL)
	pattern = strings.ReplaceAll(pattern, "{imgAlt}", taskContext.imgAlt)
	pattern = strings.ReplaceAll(pattern, "{linkURL}", taskContext.linkURL)
	pattern = strings.ReplaceAll(pattern, "{linkText}", taskContext.linkText)
	pattern = strings.ReplaceAll(pattern, "{.ext}", taskContext.extension)

	// counter
	counterPattern := regexp.MustCompile(`\{_*i\}`)
	if counterPattern.MatchString(pattern) {
		i := taskContext.counter.GetCount()
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

	return pattern
}

// GetExtension extracts extension from an URL
func GetExtension(urlString string) (string, error) {
	// file extension
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	path := u.EscapedPath()
	ext := filepath.Ext(path)

	return ext, nil
}
