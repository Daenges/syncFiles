package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	fmt.Println("Starting sync process.")
	for {
		startProc()
		time.Sleep(time.Millisecond * time.Duration(getArgs()))
	}
}

func startProc() {
	// Read config content or create one if it does not exist.
	fileData := strings.Split(string(getConfigOrCreateIt()), "\n")

	// Entry format: PATHFROM -> PathTo1 PathTo2 PathTo3
	if fileEmpty := true; len(fileData) > 0 {
		for _, line := range fileData {
			if len(line) > 0 && line[0] != '#' { 	// Allow comments
				go detectAndExecuteOperation(line) 	// todo make go later
				fileEmpty = false					// Check that file has content.
			}
		}
		if !fileEmpty {
			wg.Wait()
		} else {
			log.Fatalf("The file contains no commands!")
		}
	} else {
		log.Fatalf("The file was empty!")
	}
}

func detectAndExecuteOperation(line string) {

	if strings.Contains(line, "<->") {
		pathsTo := getPathsInQuotes(line)
		mostRecent := getNewestFile(pathsTo)

		if mostRecent != "" {
			for _, path := range pathsTo {
				if len(path) > 0 {
					if path[len(path)-1:] != getPathSeperator() && !isFile(path) {
						path += getPathSeperator()
					}
					if path[len(path)-1:] == getPathSeperator() {
						path += filepath.Base(mostRecent)
					}

					if isFile(path) {
						if !isSameFile(mostRecent, path) && path != mostRecent {
							wg.Add(1)
							go fileCopy(mostRecent, path)
						}
					} else {
						wg.Add(1)
						go fileCopy(mostRecent, path)
					}
				}
			}
		}
	} else if strings.Contains(line, "|->") {
		pathFrom := getPathsInQuotes(strings.Split(line, "->")[0])[0]
		pathsTo := getPathsInQuotes(strings.Split(line, "|->")[1])

		for _, path := range pathsTo {
			if len(path) > 0 {
				if path[len(path)-1:] != getPathSeperator() && !isFile(path) {
					path += getPathSeperator()
				}
				if path[len(path)-1:] == getPathSeperator() {
					path += filepath.Base(pathFrom)
				}

				if isFile(path) {
					if !isSameFile(pathFrom, path) {
						wg.Add(1)
						go fileCopy(pathFrom, path)
					}
				} else {
					wg.Add(1)
					go fileCopy(pathFrom, path)
				}
			}
		}
	} else if strings.Contains(line, "->") {

		pathFrom := getPathsInQuotes(strings.Split(line, "->")[0])[0]
		pathsTo := getPathsInQuotes(strings.Split(line, "->")[1])

		for _, path := range pathsTo {
			if len(path) > 0 {
				if path[len(path)-1:] != getPathSeperator() && !isFile(path) {
					path += getPathSeperator()
				}
				if path[len(path)-1:] == getPathSeperator() {
					path += filepath.Base(pathFrom)
				}

				if isFile(path) {
					if !isSameFile(pathFrom, path) && a_NEWER_b(pathFrom, path) {
						wg.Add(1)
						go fileCopy(pathFrom, path)
					}
				} else {
					wg.Add(1)
					go fileCopy(pathFrom, path)
				}
			}
		}
	}
}
