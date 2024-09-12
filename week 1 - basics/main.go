package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var prefix = ""

func dirTree(out io.Writer, dirPath string, printFiles bool) error {

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	if !printFiles {
		files = func(f []fs.DirEntry) (ret []fs.DirEntry) {
			for _, val := range f {
				if val.IsDir() {
					ret = append(ret, val)
				}
			}
			return
		}(files)
	}

	filesLen := len(files)

	for idx, file := range files {
		newPath := filepath.Join(dirPath, file.Name())
		output := ""

		if filesLen == idx+1 {
			output += "└───"
		} else {
			output += "├───"
		}
		output += filepath.Base(newPath)

		if !file.IsDir() && printFiles {
			countBytes := "empty"
			fileInfo, err := file.Info()
			if err == nil && fileInfo.Size() != 0 {
				countBytes = fmt.Sprintf("%vb", fileInfo.Size())
			}

			fmt.Fprintf(out, "%v%v (%v)\n", prefix, output, countBytes)
		}

		if file.IsDir() {
			fmt.Fprintf(out, "%v%v\n", prefix, output)

			if filesLen != idx+1 {
				prefix += "│\t"
			} else {
				prefix += "\t"
			}

			dirTree(out, newPath, printFiles)

			if filesLen != idx+1 {
				// Слайс по байтам, нужно учитывать руны
				prefix = prefix[:len(prefix)-4]
			} else {
				prefix = prefix[:len(prefix)-1]
			}

		}
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
