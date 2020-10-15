package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	name string
	path string
	size int64
	date time.Time
}

func getFiles(path string) []File {
	var files []File
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			file := File{name: info.Name(), path: path, size: info.Size(), date: info.ModTime()}
			files = append(files, file)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func contains(files *[]File, file *File) (bool, *File) {
	for _, f := range *files {
		if f.name == file.name && f.size == file.size {
			return true, &f
		}
	}
	return false, nil
}

func getDuplicates(files *[]File) *[]File {
	// list with all unique files
	var unique []File
	// list with unique files that have at least 1 duplicate
	var uniqueDuplicates []File
	var duplicateFiles []File
	for _, file := range *files {
		cont, f := contains(&unique, &file)
		if !cont {
			// if file isnt registered in unique, add it
			unique = append(unique, file)
			continue
		}
		if cont {
			// if it exists in unique, add to duplicate
			duplicateFiles = append(duplicateFiles, file)
			contD, _ := contains(&uniqueDuplicates, f)
			if !contD {
				// if file doesnt exist in unique duplicates, add it
				uniqueDuplicates = append(uniqueDuplicates, *f)
			}
		}
	}
	for _, file := range uniqueDuplicates {
		duplicateFiles = append(duplicateFiles, file)
	}
	return &duplicateFiles
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	path := "./"
	files := getFiles(path)
	duplicates := getDuplicates(&files)

	os.MkdirAll("./temp", os.ModePerm)
	f, err := os.Create("./temp/results.txt")
	check(err)

	defer f.Close()

	for _, file := range *duplicates {
		text := fmt.Sprintf("name: %s; path: %s2 size: %d \n", file.name, file.path, file.size)
		_, err3 := f.WriteString(text)
		check(err3)

		f.Sync()
		// fmt.Println(file.name, "-", file.size)
	}
}
