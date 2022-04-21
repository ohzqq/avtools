package fftools

import (
	"log"
	"os"
	"bufio"
	"strings"
	"fmt"
	"regexp"
	"time"
	"strconv"
	//"reflect"
)
var _ = fmt.Printf

type Chapter struct {
	ID int
	TimeBase string
	StartTime string
	EndTime string
	Tags Tags
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
		indices []string
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
			durS := strconv.Itoa(int(dur.Seconds())) + ".000000"
			indices = append(indices, durS)
		}
	}

	var tracks Chapters
	for i, _ := range titles {
		t := Chapter{}
		t.Tags.Title = titles[i]
		t.StartTime = indices[i]
		tracks = append(tracks, &t)
	}

	return tracks
}
