package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type File struct {
	path  string
	name  string
	size  int
	order int
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
				filesList[int(info.Size())] = append(filesList[int(info.Size())], File{path, info.Name(), int(info.Size()), 0})
			} else if strings.Trim(filepath.Ext(path), ".") == format {
				filesList[int(info.Size())] = append(filesList[int(info.Size())], File{path, info.Name(), int(info.Size()), 0})
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

func md5sum(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func checkDuplicates(files map[int][]File, option int) map[int]File {
	var count = 1
	duplicates := map[int]map[string][]File{}
	for _, v := range files {
		for _, v1 := range v {
			if _, ok := duplicates[v1.size]; !ok {
				duplicates[v1.size] = map[string][]File{}
			}
			duplicates[v1.size][md5sum(v1.path)] = append(duplicates[v1.size][md5sum(v1.path)], v1)
		}
	}
	var sorted = make([]int, 0, len(duplicates))
	for k := range duplicates {
		sorted = append(sorted, k)
	}
	if option == 2 {
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i] < sorted[j]
		})
	} else if option == 1 {
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i] > sorted[j]
		})
	}
	var filesMap = make(map[int]File)
	for sortedCount := 0; sortedCount < len(sorted); {
		fmt.Printf("\n%d bytes\n", sorted[sortedCount])
		for hash, v1 := range duplicates[sorted[sortedCount]] {
			if len(v1) > 1 {
				fmt.Printf("Hash: %s\n", hash)
				for _, v2 := range v1 {
					v2.order = count
					filesMap[count] = v2
					fmt.Printf("%d. %s\n", count, v2.path)
					count++
				}
			}
		}
		sortedCount++
	}
	return filesMap
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func deleteDuplicates(files map[int]File) {
	var list []int
	var correct = false
	var totalSize = 0
	for !correct {
		fmt.Println("Enter file numbers to remove:")
		listString := strings.Split(readLine(), " ")
		for _, v := range listString {
			num, _ := strconv.Atoi(v)
			if num > len(files) || num < 1 {
				fmt.Println("Wrong option")
				break
			} else {
				list = append(list, num)
				correct = true
			}
		}
	}
	for _, v := range list {
		for _, v1 := range files {
			if v1.order == v {
				totalSize += v1.size
				deleteFile(v1.path)
				delete(files, v)
			}
		}
	}
	fmt.Printf("Total freed up space: %d bytes", totalSize)
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
		fmt.Println("Check for duplicates?")
		var filesMap map[int]File
		for {
			var input = readLine()
			if input == "yes" {
				filesMap = checkDuplicates(filesList, option)
				break
			} else if input == "no" {
				return
			} else {
				fmt.Println("Wrong option")
			}
		}
		deleteDuplicates(filesMap)
	}
}
