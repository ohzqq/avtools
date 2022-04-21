package fftools

import (
	"log"
	"os"
	"bufio"
	"strings"
	"fmt"
	"regexp"
	"time"
	//"reflect"
)
var _ = fmt.Printf

type Chapter struct {
	Title string
	Start time.Duration
	End time.Duration
}

type Chapters []*Chapter

func (c *Chapter) Timestamps() {
}

func ReadCueSheet(file string) Chapters {
	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	var (
		titles []string
		indices []time.Duration
	)
	scanner := bufio.NewScanner(contents)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if strings.Contains(s, "TITLE") {
			t := strings.TrimPrefix(s, "TITLE ")
			t = strings.Trim(t, "'")
			t = strings.Trim(t, `"`)
			titles = append(titles, t)
		} else if strings.Contains(s, "INDEX") {
			start := strings.TrimPrefix(s, "INDEX 01 ")
			rmFrames := regexp.MustCompile(`:\d\d$`)
			start = rmFrames.ReplaceAllString(start, "s")
			start = strings.ReplaceAll(start, ":", "m")
			dur, _ := time.ParseDuration(start)
			indices = append(indices, dur)
		}
	}

	var tracks Chapters
	for i, _ := range titles {
		t := Chapter{}
		t.Title = titles[i]
		t.Start = indices[i]
		tracks = append(tracks, &t)
	}

	return tracks
}
