package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type File struct {
	path string
	name string
	size int
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func dirSearch(arg string, format string, option int) map[int][]File {
	filesList := make(map[int][]File)
	err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Directory is not specified")
			return nil
		}
		if !info.IsDir() {
			if format == "" {
				filesList[int(info.Size())] = append(filesList[int(info.Size())], File{path, info.Name(), int(info.Size())})
			} else if strings.Trim(filepath.Ext(path), ".") == format {
				filesList[int(info.Size())] = append(filesList[int(info.Size())], File{path, info.Name(), int(info.Size())})
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}
	return filesList
}

func sortFiles(files map[int][]File, option int) (map[int][]File, []int) {
	keys := make([]int, 0, len(files))

	for k := range files {
		keys = append(keys, k)
	}
	if option == 2 {
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
	} else {
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] > keys[j]
		})
	}
	return files, keys
}

func filePrinter(files map[int][]File, option int) {
	files, keys := sortFiles(files, option)
	for _, k := range keys {
		fmt.Printf("\n%d bytes\n", k)
		for _, v := range files[k] {
			fmt.Println(v.path)
		}
	}
}

var options = []string{"Descending", "Ascending"}

func optionSetter(options []string) (string, int) {
	var format string
	var option int
	var correct = false
	fmt.Println("Enter file format:")
	format = readLine()
	fmt.Println("Size sorting options:")
	for i, v := range options {
		fmt.Println(i+1, v)
	}
	for !correct {
		option, _ = strconv.Atoi(readLine())
		if option > len(options) || option < 1 {
			fmt.Println("Wrong option")
		} else {
			correct = true
		}
	}
	return format, option
}

func main() {
	filesList := make(map[int][]File)
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		return
	} else {
		format, option := optionSetter(options)
		filesList = dirSearch(os.Args[1], format, option)
		filePrinter(filesList, option)
	}
}
