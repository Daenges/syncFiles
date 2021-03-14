package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	fileData := strings.Split(string(getConfigOrCreateIt()), "\n")

	// Entry format: PATHFROM -> PathTo1 PathTo2 PathTo3
	if !isAllComented(fileData) && len(fileData) > 0 {
		for _, line := range fileData {
			if len(line) > 0 && line[0] != '#' { // Allow comments
				go detectAndExecuteOperation(line) // todo make go later
			}
		}
		wg.Wait()
	} else {
		log.Fatalf("The file was empty!")
	}
}

func getConfigOrCreateIt() (fileContent []byte) {
	if wd, err := os.Getwd(); isFile(wd+getPathSeperator()+"config.txt") && err == nil {
		fileContent, err = ioutil.ReadFile(wd + getPathSeperator() + "config.txt")
		check(err)
	} else if len(os.Args[1:]) > 0 {
		fileContent, err = ioutil.ReadFile(os.Args[1])
		check(err)
	} else {

		if path, err := os.Getwd(); err == nil {
			emptyFile, err := os.Create(path + getPathSeperator() + "config.txt")
			check(err)
			defer emptyFile.Close()

			fileContent = []byte("# Enter configuration (one File per line) as described in the help page.")
			emptyFile.Write(fileContent)

			log.Printf("Could not find config. Created one at %v%vconfig.txt", path, getPathSeperator())
		} else {
			log.Fatal("Could not get working directory. Please create a config file, so the program can operate properly.")
		}
	}
	return
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
