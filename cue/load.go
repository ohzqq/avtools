package cue

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func Load(file string) *Sheet {
	var sheet Sheet

	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	scanner := bufio.NewScanner(contents)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.Contains(line, "TITLE"):
			sheet.titles = append(sheet.titles, title(line))
		case strings.Contains(line, "INDEX 01"):
			sheet.startTimes = append(sheet.startTimes, start(line))
		}
	}

	e := 1
	for i := 0; i < len(sheet.titles); i++ {
		var t Track
		t.title = sheet.titles[i]
		t.start = sheet.startTimes[i]
		if e < len(sheet.titles) {
			t.end = sheet.startTimes[e]
		}
		e++
		sheet.Tracks = append(sheet.Tracks, t)
	}

	return &sheet
}

func file(line string) string {
	fileRegexp := regexp.MustCompile(`^(\w+ )('|")(?P<title>.*)("|')( .*)$`)
	matches := fileRegexp.FindStringSubmatch(line)
	title := matches[fileRegexp.SubexpIndex("title")]
	return title
}

func title(line string) string {
	t := strings.TrimPrefix(line, "TITLE ")
	t = strings.Trim(t, "'")
	t = strings.Trim(t, `"`)
	return t
}

func start(line string) int {
	stamp := strings.TrimPrefix(line, "INDEX 01 ")
	split := strings.Split(stamp, ":")
	dur, err := time.ParseDuration(split[0] + "m" + split[1] + "s")
	if err != nil {
		log.Fatal(err)
	}
	return int(dur.Seconds())
}
