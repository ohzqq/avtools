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
	seconds, err := strconv.Atoi(strings.Split(sec, ".")[0])
	if err != nil { log.Fatal(err) }
	//seconds := float64(s)
	h := seconds / 3600
	m := seconds % 3600 / 60
	s := seconds % 3600 % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	//return s
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

