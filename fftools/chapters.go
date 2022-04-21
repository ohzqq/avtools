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

type Chapters []Chapter

type Chapter struct {
	ID int
	TimeBase string `json:"time_base"`
	StartTime string `json:"start_time"`
	ss string
	to string
	EndTime string `json:"end_time"`
	Tags Tags `json:"tags"`
}

func (c Chapters) Timestamps() {
	var end string
	eCh := 0
	for i := 1; i < len(c); i++ {
		c[eCh].ss = c[eCh].StartTime
		switch e := c[i].EndTime; e {
		default:
			end = c[eCh].EndTime
		case "":
			end = c[i].StartTime
		}
		c[eCh].to = end
		eCh++
	}
		fmt.Printf("%v\n", c)
}

func ReadCueSheet(file string) MediaMeta {
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
			durS := strconv.Itoa(int(dur.Seconds()))
			indices = append(indices, durS)
		}
	}

	var tracks Chapters
	for i, _ := range titles {
		t := Chapter{}
		t.Tags.Title = titles[i]
		t.StartTime = indices[i]
		tracks = append(tracks, t)
	}

	return MediaMeta{Chapters: tracks}
}
