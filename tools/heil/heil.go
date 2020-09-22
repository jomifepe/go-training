package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

var (
	tail      bool
	argLines  uint64
	argBytes  uint64
	showLines bool
)

func main() {
	flag.BoolVarP(&showLines, "line-numbers", "n", false, "show the number at the beginning of each line")
	flag.BoolVarP(&tail, "tail", "t", false, "tail: starts reading from the end of the file")
	flag.Uint64VarP(&argLines, "lines", "l", 0, "number of lines to read")
	flag.Uint64VarP(&argBytes, "bytes", "b", 0, "number of bytes to read")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fatal("Please specify a file to read")
	}

	absPath, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		fatal("Couldn't get specified file -", err.Error())
	}

	file, err := os.Stat(absPath)
	if err != nil {
		fatal("Couldn't read specified file -", err.Error())
	} else if file.IsDir() {
		fatal("Specified path corresponds to a directory")
	}

	read(absPath, file.Size())
}

func read(path string, size int64) {
	file, err := os.Open(path)
	if err != nil {
		fatal("Couldn't open specified file -", err.Error())
	}
	defer file.Close()

	if argLines > 0 {
		readLines(file, size)
	} else if argBytes > 0 {
		readBytes(file, size)
	} else {
		argLines = 10
		readLines(file, size)
	}
}

func readBytes(file *os.File, size int64) {
	if tail {
		if _, err := file.Seek(-(int64(argBytes) % (size + 1)), 2); err != nil {
			fatal("Couldn't read file -", err.Error())
		}
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, argBytes)
	if _, err := reader.Read(buffer); err != nil && err != io.EOF {
		fatal("Couldn't read file -", err.Error())
	}
	fmt.Print(string(buffer))
}

func readLines(file *os.File, size int64) {
	if tail {
		tailLines(file, size)
	} else {
		headLines(file)
	}
}

func headLines(file *os.File) {
	n, scanner := uint64(0), bufio.NewScanner(file)
	for scanner.Scan() {
		if showLines {
			fmt.Printf("%v: %v\n", n + 1, scanner.Text())
		} else {
			fmt.Println(scanner.Text())
		}
		n++
		if n == argLines {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func tailLines(file *os.File, size int64) {
	var (
		lines []string
		line string
		buff = make([]byte, 1)
	)

	prependLine := func(elem string) {
		lines = append([]string{elem}, lines...)
		line = ""
	}

	for pos := size; pos >= 0; pos-- {
		rc, err := file.ReadAt(buff, pos)
		if err != nil && err != io.EOF {
			fatal("Couldn't read file -", err.Error())
		}

		if len(lines) == int(argLines) /* reached the define number of lines to read*/ {
			break
		}
		if string(buff) == "\n" /* changed line */ {
			if len(line) > 0 {
				prependLine(line)
			}
			rc = 0
		}
		if rc > 0 /* read something from the file */ {
			line = string(buff) + line
			if pos == 0 /* reached the top of the file */ {
				prependLine(line)
			}
		}
	}

	var lineCount int
	if showLines {
		lineCount = getLineCount(file)
	}
	for i, l := range lines {
		if showLines {
			fmt.Printf("%v: %v\n", (lineCount - len(lines) + i) + 1, l)
			continue
		}
		fmt.Println(l)
	}
}

func getLineCount(file *os.File) (count int) {
	_, _ = file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanner.Text()
		count++
	}
	return
}

func fatal(msg ...string) {
	_, _ = fmt.Fprintln(os.Stderr, strings.Join(msg, " "))
	os.Exit(1)
}