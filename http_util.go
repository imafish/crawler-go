package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// DownloadFile downloads a file from the given url and save to the destinatin file
func DownloadFile(url, destination string) error {
	Log.Infof("Downloading file from %s to %s", url, destination)

	resp, err := http.Get(url)
	if err != nil {
		Log.Errorf("Failed to complete request. err: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destination)
	if err != nil {
		Log.Errorf("Failed to create destination file. err: %s", err.Error())
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		Log.Errorf("Failed to write to destination file. err: %s", err.Error())
	}
	return err
}

// GrabText en..
func GrabText(pageContent, xpath string) (string, error) {
	return "", fmt.Errorf("Unimplemented")
}

// GrabPage grabs the entire page content as string
func GrabPage(url string) (string, error) {
	Log.Infof("Getting page content from url: %s", url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	resp, err := client.Do(req)
	if err != nil {
		Log.Errorf("Failed to complete request. err: %s", err.Error())
		return "", err
	}

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Got unexpected status code: %d", resp.StatusCode)
		Log.Errorf(msg)
		return "", errors.New(msg)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Errorf("Failed to read from response body. err: %s", err.Error())
		return "", err
	}

	return string(body), nil
}
