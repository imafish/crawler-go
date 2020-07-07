package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// Log is the global logger
var Log Logger

func main() {
	url := flag.String("u", "", "URL for the crawler to parse")
	outDir := flag.String("d", "", "output directory for downloaded resources")
	logPath := flag.String("-log", "", "log file path")
	flag.Parse()
	if *url == "" || *outDir == "" {
		usage()
		os.Exit(1)
	}

	initLog(*logPath)

}

func initLog(logPath string) {
	if logPath == "" {
		switch runtime.GOOS {
		case "windows":
			dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			logPath = path.Join(dir, "log.log")
		case "linux":
			logPath = "/etc/var/log/"
		default:
		}
	}

	// TODO implement log
	Log = nil
}

func usage() {
	fmt.Fprint(os.Stderr, "Usage:\n")
	flag.PrintDefaults()
}
