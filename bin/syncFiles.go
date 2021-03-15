package main

import (
	"fmt"
	"log"
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
				/*go*/ detectAndExecuteOperation(line) 	// todo make go later
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

	if strings.Contains(line, "<->") { 		// Get the newest file version and replace all the others with it.
		pathsTo := getPathsInQuotes(line)
		mostRecent := getNewestFile(pathsTo)

		if mostRecent != "" {
			startCopyProcess(mostRecent, "<->", pathsTo)
		}

	} else if strings.Contains(line, "|->") {		// Always copy this file, if the other one changed.
		pathFrom := getPathsInQuotes(strings.Split(line, "|->")[0])[0]
		pathsTo := getPathsInQuotes(strings.Split(line, "|->")[1])
		startCopyProcess(pathFrom, "|->", pathsTo)

	} else if strings.Contains(line, "->") {

		pathFrom := getPathsInQuotes(strings.Split(line, "->")[0])[0]
		pathsTo := getPathsInQuotes(strings.Split(line, "->")[1])
		startCopyProcess(pathFrom, "->", pathsTo)
	}
}

func startCopyProcess(pathFrom, operator string, pathsTo []string) {
	for _, path := range pathsTo {
		if len(path) > 0 {
			path = preparePath(path, pathFrom)

			if same, NotexistentErr := isSameFile(pathFrom, path); NotexistentErr != nil {
				addToWgAndCopy(pathFrom, path)
			} else {
				if operator == "->" && !same && a_NEWER_b(pathFrom, path) {
					addToWgAndCopy(pathFrom, path)
				}
				if operator == "|->" && !same {
					addToWgAndCopy(pathFrom, path)
				}
				if operator == "<->" && !same && path != pathFrom {
					addToWgAndCopy(path, pathFrom)
				}
			}

		}
	}
}
