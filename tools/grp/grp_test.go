package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	binaryName = "grp"
	binaryPath string
	tmpFile *os.File
	tmpFilePath string
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	fmt.Println("Setting up...")
	cmd := exec.Command("go", "build", "-o", binaryName)
	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to execute test setup: ", err.Error())
	}
	p, err := filepath.Abs(binaryName)
	if err != nil {
		log.Fatal("Failed to execute test setup: ", err.Error())
	}
	binaryPath = p

	tmpFile = openTmpFile(binaryName)
	p, err = filepath.Abs(tmpFile.Name())
	if err != nil {
		log.Fatal("Failed to execute test setup: ", err.Error())
	}
	tmpFilePath = p
}

func teardown() {
	fmt.Println("Tearing down...")
	os.Remove(tmpFile.Name())
	// teardown code
}

func TestCaseSensitiveMatching(t *testing.T) {
	writeToTmpFile("user\n123\nUser\n456\nUSER")

	cmdArgs := []string{"user", tmpFilePath}
	cmd := exec.Command(binaryPath, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	expectedLines := 1
	lines := countRune(string(output), '\n') + 1
	if lines != expectedLines {
		t.Errorf("Expected %v matched lines, got %v.\nCommand: %v %v\nOutput:\n%v",
			expectedLines, lines, binaryPath, strings.Join(cmdArgs, " "), string(output))
	}
}

func TestCaseInsensitiveMatching(t *testing.T) {
	writeToTmpFile("user\n123\nUser\n456\nUSER")

	cmdArgs := []string{"-i", "user", tmpFilePath}
	cmd := exec.Command(binaryPath, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	expectedLines := 3
	lines := countRune(string(output), '\n') + 1
	if lines != expectedLines {
		t.Errorf("Expected %v matched lines, got %v.\nCommand: %v %v\nOutput:\n%v",
			expectedLines, lines, binaryPath, strings.Join(cmdArgs, " "), string(output))
	}
}

func TestMaximumMatchedLines(t *testing.T) {
	writeToTmpFile("user\n123\nUser\n456\nUSER")

	cmdArgs := []string{"-il", "2", "user", tmpFilePath}
	cmd := exec.Command(binaryPath, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	expectedLines := 2
	lines := countRune(string(output), '\n') + 1
	if lines != expectedLines {
		t.Errorf("Expected %v matched lines, got %v.\nCommand: %v %v\nOutput:\n%v",
			expectedLines, lines, binaryPath, strings.Join(cmdArgs, " "), string(output))
	}
}

func TestMaximumMatchedStrings(t *testing.T) {
	writeToTmpFile("user\n123\nUser:superuser\n456\nUSER")

	cmdArgs := []string{"-im", "3", "user", tmpFilePath}
	cmd := exec.Command(binaryPath, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	expectedLines := 2
	lines := countRune(string(output), '\n') + 1
	if lines != expectedLines {
		t.Errorf("Expected %v matched lines, got %v.\nCommand: %v %v\nOutput:\n%v",
			expectedLines, lines, binaryPath, strings.Join(cmdArgs, " "), string(output))
	}
}

func TestShowLineNumbers(t *testing.T) {
	writeToTmpFile("user\n123\nUser:superuser\n456\nUSER\nusers")

	cmdArgs := []string{"-in", "user", tmpFilePath}
	cmd := exec.Command(binaryPath, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.SplitAfter(string(output), "\n")
	for _, line := range lines {
		if _, err := strconv.Atoi(line[:1]); err != nil {
			t.Errorf("Found a matching line that isn't starting with a number.\nCommand: %v %v\nOutput:\n%v",
				binaryPath, strings.Join(cmdArgs, " "), string(output))
			break
		}
	}
}

func TestShowNoColoredHighlighting(t *testing.T) {
	writeToTmpFile("user\n123\nUser")

	cmdArgs := []string{"--no-color", "user", tmpFilePath}
	cmd := exec.Command(binaryPath, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(output), "\u001B[") {
		t.Errorf("Found a matching line highlighted with color.\nCommand: %v %v\nOutput:\n%v",
			binaryPath, strings.Join(cmdArgs, " "), string(output))
	}
}

// Helpers

func openTmpFile(fileBaseName string) *os.File {
	f, err := ioutil.TempFile("", fileBaseName)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func writeToTmpFile(content string) {
	if err := tmpFile.Truncate(0); err != nil {
		log.Fatal(err)
	}
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
}

func writeTmpFile(content []byte, fileBaseName string) *os.File {
	tmpFile, err := ioutil.TempFile("", fileBaseName)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		log.Fatal(err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	return tmpFile
}

func countRune(s string, r rune) (count int) {
	for _, c := range s {
		if c == r {
			count++
		}
	}
	return
}