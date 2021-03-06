package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileCopy(src, dst string) {
	defer wg.Done()
	fmt.Printf("Copy %v -> %v", src, dst)

	sourceFileStat, err := os.Stat(src)
	check(err)

	if !sourceFileStat.Mode().IsRegular() {
		check(fmt.Errorf("%s is not a regular file", src))
	}

	source, err := os.Open(src)
	check(err)
	defer source.Close()

	destination, err := os.Create(dst)
	check(err)
	defer destination.Close()

	_, err = io.Copy(destination, source)
	check(err)
}

func isAllComented(fileContent []string) bool {
	for _, line := range fileContent {
		if line[0] != '#' {
			return false
		}
	}
	return true
}

func getArgs() int {
	argSlice := os.Args[1:]

	for num, arg := range argSlice {
		if arg == "-t" {
			return int(arg[num+1])
		}
		if arg == "-h" || arg == "--help" {
			printHelp()
			os.Exit(1)
		}
	}
	return 5000
}

func printHelp() {
	fmt.Println("--- Help Page ---")
	fmt.Println("Execution: syncFile.exe [Path of Config] (arg) (val)")
	fmt.Println("-t [value] = time to sleep between checks")
	fmt.Println("-h / --help = get this page")
	fmt.Println("[] - necessary | () - optional")
	fmt.Println()
	fmt.Println("---- Config Format ---")
	fmt.Println("\"Path from file\" -> \"Path to copy to (1)\" \"Path to copy to (2)\"...")
	fmt.Println("- Copies the file if it is newer than the file at the destination. -")
	fmt.Println("\"Path from file\" |-> \"Path to copy to (1)\" \"Path to copy to (2)\"...")
	fmt.Println("- Always copy if something in the first file has changed. -")
	fmt.Println("\"Path to copy to (1)\" <-> \"Path to copy to (2)\" <-> \"Path to copy to (3)\"")
	fmt.Println("- Takes the last edited file and copies it to all other locations. -")
	fmt.Println("# - To Comment a line")
}

func isSameFile(pathFileA, pathFileB string) bool {
	fileContentA, err := os.Open(pathFileA)
	if err != nil{
		return false
	} else {
		defer fileContentA.Close()
	}

	fileContentB, err := os.Open(pathFileB)
	if err != nil{
		return false
	} else {
		defer fileContentA.Close()
	}

	hashfuncA := sha256.New()
	hashfuncB := sha256.New()

	if _, err := io.Copy(hashfuncA, fileContentA); err != nil {
		check(err)
	}

	hashValA := hex.EncodeToString(hashfuncA.Sum(nil))

	if _, err := io.Copy(hashfuncB, fileContentB); err != nil {
		check(err)
	}

	hashValB := hex.EncodeToString(hashfuncB.Sum(nil))

	return hashValA == hashValB
}

func a_NEWER_b(pathA, pathB string) bool {

	fileA, err := os.Stat(pathA)
	check(err)
	fileB, err := os.Stat(pathB)
	check(err)

	return time.Now().Sub(fileA.ModTime()) < time.Now().Sub(fileB.ModTime())
}

func getPathsInQuotes(pathsInQuotes string) (pathArray []string) {

	currentPosInArray := -1
	beginCopy := false

	for _, letter := range pathsInQuotes {
		if letter == '"' {

			if beginCopy {
				beginCopy = false
			} else {
				beginCopy = true
				pathArray = append(pathArray, "")
				currentPosInArray++
			}
		} else {
			if beginCopy && pathArray != nil {
				pathArray[currentPosInArray] += string(letter)
			}
		}
	}

	return
}

func isFile(path string) bool {
	file, err := os.Stat(path)
	if err != nil {			// unable to open file
		return false
	} else if file.Mode().IsRegular() {	// is a valid file!
		return true
	} else {	// whatever else could happen
		return false
	}
}

func getPathSeperator() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

func getNewestFile(pathFiles []string) (pathNewest string) {

	for _, path := range pathFiles {
		if isFile(path) {
			pathNewest = path
			break
		}
	}

	for _, path := range pathFiles {
		if len(path) > 0 {
			if path[len(path)-1:] != getPathSeperator() && !isFile(path) {
				path += getPathSeperator() + filepath.Base(pathNewest)
			}


			if a_NEWER_b(path, pathNewest) && isFile(path) && path != pathNewest {
				pathNewest = path
			}
		}
	}
	return
}