package main

import (
	"fmt"
	"io"
	"os"
)

type stdoutWriter struct{}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify a file to read")
		os.Exit(1)
	}

	info, err := os.Stat(os.Args[1])
	if err == nil {
		if info.IsDir() {
			fmt.Println("Specified path is a directory")
			os.Exit(1)
		}
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println("Failed to read file:", err)
			os.Exit(1)
		}
		io.Copy(os.Stdout, file)
		// io.Copy(stdoutWriter{}, file)
	} else if os.IsNotExist(err) {
		fmt.Println("Specified filed does't exist")
		os.Exit(1)
	} else {
		fmt.Println("Failed to open specified file")
		os.Exit(1)
	}

	os.Exit(0)
}

// For demonstration purposes
func (stdoutWriter) Write(bs []byte) (int, error) {
	fmt.Printf("Reading file %v...\n\n", os.Args[1])
	fmt.Println(string(bs))
	return len(bs), nil
}
