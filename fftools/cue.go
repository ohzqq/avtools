package fftools

import (
	"log"
	"os"
	"bufio"
	"strings"
	"fmt"
	"regexp"
	"time"
)

type track struct {
	title string
	startTime time.Duration
}

func ReadCueSheet(file string) {
	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	scanner := bufio.NewScanner(contents)
	var titles []string
	var indices []time.Duration
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if strings.Contains(s, "TITLE") {
			titles = append(titles, strings.TrimPrefix(s, "TITLE "))
		} else if strings.Contains(s, "INDEX") {
			rmFrames := regexp.MustCompile(`:\d\d$`)
			start := strings.TrimPrefix(s, "INDEX 01 ")
			start = rmFrames.ReplaceAllString(start, "s")
			start = strings.ReplaceAll(start, ":", "m")
			dur, _ := time.ParseDuration(start)
			indices = append(indices, dur)
		}
	}

	var tracks []track
	for i, _ := range titles {
		t := track{}
		t.title = titles[i]
		t.startTime = indices[i]
		tracks = append(tracks, t)
	}

	fmt.Printf("%v", tracks)
}
