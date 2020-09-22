package util

import (
	"crypto/rand"
	"fmt"
	"math"
	"os"
	"time"
)

type TimeRep struct {
	time time.Duration
	rep string
}

func GetRandStringBytes(n int) string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

	for k, v := range buf {
		buf[k] = characters[v % byte(len(characters))]
	}
	return string(buf)
}

func GetSizeString(bytes int) string {
	if bytes <= 0 {
		return "0B"
	}

	b, k := float64(bytes), 1024.0
	i := math.Floor(math.Log(b) / math.Log(k))
	return fmt.Sprintf("%v%v", math.Round(b / math.Pow(k, i)),
		[]string{"B", "K", "M", "G", "T", "P", "E", "Z", "Y"}[int(i)])
}

func GetTimeString(nanoseconds int64) string {
	m := []TimeRep {
		{time.Hour, "h"},
		{time.Minute, "min"},
		{time.Second, "s"},
		{time.Millisecond, "ms"},
		{time.Microsecond, "μs"},
	}

	for _, r := range m {
		if t := nanoseconds / int64(r.time); t > 1 {
			return fmt.Sprintf("%v%v", t, r.rep)
		}
	}

	return fmt.Sprintf("%vns", nanoseconds)
}

func AppendToFile(filePath string, textToAppend string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(textToAppend); err != nil {
		return err
	}
	return nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func OmitEmptyFields(s map[string]interface{}) map[string]interface{} {
	n := make(map[string]interface{})
	for key, value := range s {
		if value != nil {
			switch t := value.(type) {
			case string:
				if len(t) > 0 {
					n[key] = value
				}
			default:
				n[key] = value
			}
		}
	}

	return n
}