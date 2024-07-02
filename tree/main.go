package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func dirTree(out io.Writer, path string, includeFiles bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	return dirTreeHelper(out, path, info, includeFiles, "")
}

func dirTreeHelper(out io.Writer, path string, info os.FileInfo, includeFiles bool, prefix string) error {
	if !info.IsDir() {
		return fmt.Errorf("%s is`t a  directory", path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if includeFiles || entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	sort.Strings(names)

	for i, name := range names {
		entryPath := filepath.Join(path, name)
		entryInfo, err := os.Stat(entryPath)
		if err != nil {
			return err
		}

		var connector string
		if i == len(names)-1 {
			connector = "└───"
		} else {
			connector = "├───"
		}

		if entryInfo.IsDir() {
			fmt.Fprintln(out, prefix+connector+name)
			newPrefix := prefix
			if i == len(names)-1 {
				newPrefix += "        "
			} else {
				newPrefix += "│       "
			}
			err := dirTreeHelper(out, entryPath, entryInfo, includeFiles, newPrefix)
			if err != nil {
				return err
			}
		} else if includeFiles {
			size := entryInfo.Size()
			sizeInfo := fmt.Sprintf("%db", size)
			if size == 0 {
				sizeInfo = "empty"
			}
			fmt.Fprintln(out, prefix+connector+name+" ("+sizeInfo+")")
		}
	}
	return nil
}

func main() {
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(os.Stdout, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
