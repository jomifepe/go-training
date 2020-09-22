package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
)

const (
	ColorReset = "\u001b[0m"
	ColorFRed  = "\u001b[31;1m"
)

var (
	argExpr           string
	argFilePath       string
	argInsensitive    bool
	argMaxMatchCount  int
	argMaxLineCount   int
	argShowLineNumber bool
	argNoColor        bool
	argAfterContext   int
	argBeforeContext  int

	errLogger = log.New(os.Stderr, "", 0)
)

func main() {
	parseArgs()

	args := flag.Args()
	if len(args) == 0 {
		fatal("insufficient number of arguments")
	}

	argExpr = args[0]
	if len(args) == 2 {
		absPath, err := filepath.Abs(args[1])
		if err != nil {
			fatal("couldn't get file path -", err.Error())
		}
		argFilePath = absPath
	}

	highlightedLines := search()
	if len(highlightedLines) == 0 {
		os.Exit(1)
	}
	fmt.Print(strings.Join(highlightedLines, "\n"))
}

// parseArgs parses the flags and arguments passed in and assigns them to the correct variables
func parseArgs() {
	flag.BoolVarP(&argInsensitive, "ignore-case", "i", false,
		"Perform case insensitive matching. By default, grp is case sensitive.")
	flag.IntVarP(&argMaxMatchCount, "max-count", "m", 0,
		"Stop reading the file after num matches.")
	flag.IntVarP(&argMaxLineCount, "max-lines", "l", 0,
		"Stop reading the file after num line with matches.")
	flag.BoolVarP(&argShowLineNumber, "line-number", "n", false,
		"Each output line is preceded by its relative line number in the file, starting at line 1.")
	flag.BoolVar(&argNoColor, "no-color", false,
		"If specified, matched words will not be highlighted with color.")

	// TODO
	flag.IntVarP(&argAfterContext, "after-context", "A", 0,
		"TODO: Print num lines of trailing context after each match.")
	flag.IntVarP(&argBeforeContext, "before-context", "B", 0,
		"TODO: Print num lines of leading context before each match.")
	flag.Parse()
}

// search for lines on the argFilePath or os.Stdin (if no file was given) that match the
// specified expression (argExpr), according to the rules defined by the specified flag(s).
// Returns an array with the matched lines
func search() (result []string) {
	var (
		scanner        *bufio.Scanner
		reader         io.Reader = os.Stdin
		matchCount     int
		matchLineCount int
	)

	if len(argFilePath) > 0 {
		file, err := os.Open(argFilePath)
		if err != nil {
			fatal("couldn't open file -", err.Error())
		}
		reader = file
	}
	scanner = bufio.NewScanner(reader)

	appendLine := func(line string, lineNum int) {
		if argShowLineNumber {
			result = append(result, fmt.Sprintf("%v: %v", lineNum, line))
		} else {
			result = append(result, line)
		}
	}

	lc := 1
	for scanner.Scan() {
		line := scanner.Text()
		highlightedLine, matches := handleLine(line)

		if matches > 0 {
			appendLine(highlightedLine, lc)
			matchCount += matches
			matchLineCount++

			if (argMaxMatchCount > 0 && matchCount >= argMaxMatchCount) ||
				(argMaxLineCount > 0 && matchLineCount >= argMaxLineCount) {
				break
			}
		}
		lc++
	}

	if err := scanner.Err(); err != nil {
		fatal("failed to read file -", err.Error())
	}

	return
}

// handleLine matches a single line against the specified expression (argExpr) and flag(s).
// Returns the same line (color highlighted or not) and a number of matches
func handleLine(text string) (result string, matches int) {
	if argInsensitive && strings.Contains(strings.ToLower(text), argExpr) {
		indexes := regexp.MustCompile(`(?i)` + argExpr).FindAllStringIndex(text, -1)
		if !argNoColor {
			for i := len(indexes) - 1; i >= 0; i-- {
				text = text[:indexes[i][1]] + ColorReset + text[indexes[i][1]:]
				text = text[:indexes[i][0]] + ColorFRed + text[indexes[i][0]:]
			}
		}
		return text, len(indexes)
	}
	if strings.Contains(text, argExpr) {
		indexes := regexp.MustCompile(argExpr).FindAllStringIndex(text, -1)
		if !argNoColor {
			text = strings.ReplaceAll(text, argExpr, ColorFRed + argExpr + ColorReset)
		}
		return text, len(indexes)
	}

	return "", 0
}

// fatal sends an error message to stderr and exits with and exit code of 1
func fatal(msg ...string) {
	errLogger.Fatal("grp: ", strings.Join(msg, " "))
}