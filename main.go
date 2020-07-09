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
var logGlobal Logger

func main() {
	configFile := flag.String("c", "", "Manditory: path to the configuration yaml file")
	outDir := flag.String("d", ".", "output directory for downloaded resources")
	concurrent := flag.Int("-p", 5, "concurrent count")
	logPath := flag.String("-l", "out.log", "log file path")
	flag.Parse()
	if *configFile == "" {
		usage()
		os.Exit(1)
	}

	initLog(*logPath)

	logGlobal.Info("starting crawler workflow...")
	absDir := makeAbs(filepath.Dir(os.Args[0]), *outDir)
	workflow(*configFile, absDir, *concurrent, logGlobal)

	logGlobal.Info("ALL DONE.")
}

func initLog(logPath string) {
	if logPath == "" {
		switch runtime.GOOS {
		case "windows":
			dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			logPath = path.Join(dir, logPath)
		case "linux":
			logPath = path.Join("/etc/var/log/", logPath)
		default:
		}
	}

	// TODO implement log
	writer, err := os.Create(logPath)
	if err != nil {
		writer = nil
	}
	logGlobal = ConsoleLog{
		w: writer,
	}
	if err != nil {
		logGlobal.Warningf("Failed to create log file: %s. err: %s", logPath, err.Error())
	}
}

// TODO (@imafish) move to separate file?
func makeAbs(absolute, path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	abs := filepath.Join(absolute, path)
	return abs
}

func usage() {
	fmt.Fprint(os.Stderr, "Usage:\n")
	flag.PrintDefaults()
}
