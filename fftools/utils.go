package fftools

import (
	"strconv"
	"strings"
	"log"
	"fmt"
	"os"
	"path/filepath"
)

var _ = fmt.Sprintf("%v", "")

func secsToHHMMSS(sec string) string {
	seconds := secsAtoi(sec)
	h := seconds / 3600
	m := seconds % 3600 / 60
	s := seconds % 3600 % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	//return s
}

func secsToMMSS(sec string) string {
	seconds := secsAtoi(sec)
	m := seconds / 60
	s := seconds % 60
	return fmt.Sprintf("%02d:%02d:00", m, s)
	//return s
}

func secsAtoi(sec string) int {
	seconds, err := strconv.Atoi(strings.Split(sec, ".")[0])
	if err != nil {
		log.Fatal(err)
	}
	return seconds
}

func timeBaseFloat(timebase string) float64 {
	tb := strings.ReplaceAll(timebase, "1/", "")
	baseint, _ := strconv.ParseFloat(tb, 64)
	return baseint
}

func ffChapstoSeconds(timebase, start, end string) (string, string) {
	tb := timeBaseFloat(timebase)
	ss, _ := strconv.ParseFloat(start, 64)
	to, _ := strconv.ParseFloat(end, 64)
	s := strconv.FormatFloat(ss / tb, 'f', 6, 64)
	e := strconv.FormatFloat(to / tb, 'f', 6, 64)
	return s, e
}

func find(ext string) []string {
	var files []string

	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range entries {
		if filepath.Ext(f.Name()) == ext {
			file, err := filepath.Abs(f.Name())
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, file)
		}
	}
	return files
}

