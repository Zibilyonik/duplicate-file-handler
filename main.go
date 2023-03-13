package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		return
	} else {
		err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Directory is not specified")
				return nil
			}
			if !info.IsDir() {
				fmt.Println(path)
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
}
