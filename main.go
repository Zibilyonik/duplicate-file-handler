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

func dirSearch(arg string, format string, option int) []File {
	var files []File
	err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Directory is not specified")
			return nil
		}
		if !info.IsDir() {
			if format == "" {
				files = append(files, File{strings.SplitN(path, string(os.PathSeparator), 2)[1], info.Name(), int(info.Size())})
			} else if strings.Trim(filepath.Ext(path), ".") == format {
				files = append(files, File{strings.SplitN(path, string(os.PathSeparator), 2)[1], info.Name(), int(info.Size())})
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
	if option == 2 {
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

func filePrinter(files []File) {
	size := 0
	for _, v := range files {
		if v.size == size {
			fmt.Println(v.path)
		} else {
			fmt.Printf("\n%d bytes\n\n", v.size)
			fmt.Println(v.path)
			size = v.size
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
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		return
	} else {
		format, option := optionSetter(options)
		files := dirSearch(os.Args[1], format, option)
		files = sortFiles(files, option)
		filePrinter(files)
	}
}
