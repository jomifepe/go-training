package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
)

const (
	ColorReset     = "\u001b[0m"
	ColorFBWhite   = "\u001b[37;1m"
	ColorFBRed     = "\u001b[31;1m"
	ColorFBGreen   = "\u001b[32;1m"
	ColorFBYellow  = "\u001b[33;1m"
	ColorFBBlue    = "\u001b[34;1m"
	ColorFBMagenta = "\u001b[35;1m"
	ColorFBCyan    = "\u001b[36;1m"
)

var (
	cType       = true
	fMode       os.FileMode
	fModeColors map[os.FileMode]string
	dir         string
	expr        string
	matchRegex  bool
	matchSubstr bool
	printHR     bool
)

func init() {
	fModeColors = map[os.FileMode]string {
		os.ModeDir:        ColorFBBlue,
		os.ModeSymlink:    ColorFBMagenta,
		os.ModeSocket:     ColorFBCyan,
		os.ModeCharDevice: ColorFBWhite,
		os.ModeDevice:     ColorFBWhite,
		os.ModeNamedPipe:  ColorFBYellow,
	}
}

func main() {
	typePtr := flag.StringP("type", "t", "", `type of file(s) to find: 
f - regular file
d - directory
l - symbolic link
c - unix character device
b - block device
s - unix domain socket
p - named pipe (FIFO)`)
	flag.BoolVarP(&matchSubstr, "substring", "s", false, "match a substring")
	flag.BoolVarP(&matchRegex, "regex", "r", false, "match a regex")
	flag.BoolVarP(&printHR, "human", "H", false, "print human readable size")

	flag.Parse()

	switch *typePtr {
	case "f": fMode = 0
	case "d": fMode = os.ModeDir
	case "c": fMode = os.ModeCharDevice
	case "b": fMode = os.ModeDevice
	case "s": fMode = os.ModeSocket
	case "p": fMode = os.ModeNamedPipe
	case "l": fMode = os.ModeSymlink
	default:
		cType = false
		fMode = os.ModeType
	}

	switch args := flag.Args(); {
	case len(args) == 0:
		_, _ = fmt.Fprintln(os.Stderr, "Not enough arguments")
		os.Exit(1)
	case len(args) == 2:
		expr = args[1]
		fallthrough
	default: dir = args[0]
	}

	if matchSubstr && matchRegex {
		_, _ = fmt.Fprintln(os.Stderr, "Please specify only one matching rule: regex [-r], substring [-s] or none for the whole string")
		os.Exit(1)
	}

	var fullPath string
	if dir[0] != 47 {
		currentDir, err := os.Getwd()
		if err != nil {
			panic("Failed to get current directory")
		}
		fullPath = fmt.Sprintf("%v/%v", currentDir, dir)
	} else {
		fullPath = dir
	}

	find(fullPath)
}

func find(path string) {
	var regPattern *regexp.Regexp
	if matchRegex {
		regPattern, _ = regexp.Compile(expr)
	}
	fCount := 0
	err := filepath.Walk(path,
		func(path string, file os.FileInfo, err error) error {
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				return nil
			}

			d := strings.Replace(dir, "./", "", -1)
			if file.Name() == d {
				return nil
			}
			if (len(expr) == 0 || matchesExpression(file, regPattern)) && (!cType || matchesType(file)) {
				printFile(path, file)
				fCount++
			}

			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nFound a total of %v file(s)\n", fCount)
}

func printFile(path string, file os.FileInfo) {
	sizeString := fmt.Sprint(file.Size())
	if printHR {
		sizeString = getSizeString(file.Size())
	}

	fmt.Printf("fnd: %v %v\n", getModeColoredString(path, file.Mode()), sizeString)
}

func matchesExpression(file os.FileInfo, regex *regexp.Regexp) bool {
	return (!(matchSubstr || matchRegex) && expr == file.Name()) ||
		(matchSubstr && strings.Contains(file.Name(), expr)) ||
		(matchRegex && regex.MatchString(file.Name()))
}

func matchesType(file os.FileInfo) bool {
	if fMode == 0 && file.Mode().IsRegular() {
		return true
	}
	return file.Mode() &fMode != 0
}

func getModeColoredString(str string, fm os.FileMode) string {
	if c, found := fModeColors[fMode& fm]; found {
		return fmt.Sprint(c, str, ColorReset)
	}
	return fmt.Sprint(ColorFBWhite, str, ColorReset) // regular files and others
}

func getSizeString(bytes int64) string {
	if bytes == 0 {
		return "0B"
	}

	b, k := float64(bytes), 1024.0
	i := math.Floor(math.Log(b) / math.Log(k))
	return fmt.Sprintf("%v%v", math.Round(b / math.Pow(k, i)),
		[]string{"B", "K", "M", "G", "T", "P", "E", "Z", "Y"}[int(i)])
}