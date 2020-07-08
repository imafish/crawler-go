package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// downloadFile downloads a file from the given url and save to the destinatin file
func downloadFile(url, destination string) error {
	logGlobal.Infof("Downloading file from %s to %s", url, destination)

	resp, err := http.Get(url)
	if err != nil {
		logGlobal.Errorf("Failed to complete request. err: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destination)
	if err != nil {
		logGlobal.Errorf("Failed to create destination file. err: %s", err.Error())
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logGlobal.Errorf("Failed to write to destination file. err: %s", err.Error())
	}
	return err
}

// GrabText en..
func grabText(pageContent, xpath string) (string, error) {
	return "", fmt.Errorf("Unimplemented")
}

func grabPageReader(url string) (io.ReadCloser, error) {
	logGlobal.Infof("Getting page content from url: %s", url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	resp, err := client.Do(req)
	if err != nil {
		logGlobal.Errorf("Failed to complete request. err: %s", err.Error())
		return nil, err
	}

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Got unexpected status code: %d", resp.StatusCode)
		logGlobal.Errorf(msg)
		return nil, errors.New(msg)
	}

	return resp.Body, nil
}

// grabPage grabs the entire page content as string
func grabPage(url string) (string, error) {
	reader, err := grabPageReader(url)
	if err != nil {
		logGlobal.Errorf("Failed to grab page, err: %s", err.Error())
	}
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		logGlobal.Errorf("Failed to read from response body. err: %s", err.Error())
		return "", err
	}

	return string(body), nil
}

func convertToAbsoluteURL(absolute, relative string) (string, error) {
	url, err := url.Parse(relative)
	if err != nil {
		return "", err
	}
	if !url.IsAbs() {
		baseURL, err := url.Parse(absolute)
		if err != nil {
			return "", err
		}
		url = baseURL.ResolveReference(url)
	}

	return url.String(), nil
}
