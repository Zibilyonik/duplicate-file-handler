package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type File struct {
	path string
	name string
	size int
}

func dirSearch(arg string, format string, option int) []File {
	var files []File
	err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Directory is not specified")
			return nil
		}
		if !info.IsDir() {
			if filepath.Ext(path) == format {
				files = append(files, File{path, info.Name(), int(info.Size())})
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}
	return files
}

func sortFiles(files []File, option int) []File {
	if option == 1 {
		sort.Slice(files, func(i, j int) bool {
			return files[i].size < files[j].size
		})
	} else {
		sort.Slice(files, func(i, j int) bool {
			return files[i].size > files[j].size
		})
	}
	return files
}

var options = []string{"Ascending", "Descending"}

func optionSetter(options []string) (string, int) {
	var format string
	var option int
	var correct = false
	fmt.Println("Enter file format:")
	fmt.Scan(&format)
	fmt.Println("Size sorting options:")
	for i, v := range options {
		fmt.Println(i+1, v)
	}
	for !correct {
		fmt.Scan(&option)
		if option > len(options) || option < 1 {
			fmt.Println("Wrong option")
		} else {
			correct = true
		}
	}
	return format, option
}

func main() {
	var size = 0
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		return
	} else {
		format, option := optionSetter(options)
		files := dirSearch(os.Args[1], format, option)
		files = sortFiles(files, option)
		for _, v := range files {
			if v.size == size {
				fmt.Println(v.name, "(", v.size, "b )")
			} else {
				fmt.Printf("%d\n", v.size)
				fmt.Println(v.path, v.name)
				size = v.size
			}
		}
	}
}
