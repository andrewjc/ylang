package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var sb strings.Builder

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {

			// skip the file if its a cpu file
			if !strings.Contains(path, "arithmetic") && !strings.Contains(path, "branch") && !strings.Contains(path, "mov") && !strings.Contains(path, "stack") && !strings.Contains(path, "modrm") && !strings.Contains(path, "lod") && !strings.Contains(path, "compare") && !strings.Contains(path, "common") {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				sb.Write(content)
				sb.WriteString("\n")
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	outputFile := "concatenated.go"
	err = ioutil.WriteFile(outputFile, []byte(sb.String()), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Concatenated Go files into", outputFile)
}
