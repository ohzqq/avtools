package ffmeta

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type FFMetaChapter struct {
	Base      string  `ini:"TIMEBASE"`
	StartTime float64 `ini:"START"`
	EndTime   float64 `ini:"END"`
	ChTitle   string  `ini:"title"`
}

func (ch FFMetaChapter) Start() time.Duration {
	return calculateSecs(ch.StartTime, ch.Base)
}

func (ch FFMetaChapter) End() time.Duration {
	return calculateSecs(ch.EndTime, ch.Base)
}

func (ch FFMetaChapter) Title() string {
	return ch.ChTitle
}

func calculateSecs(num float64, base string) time.Duration {
	b := timebase(base)
	t := num / b * float64(time.Second)
	return time.Duration(t)
}

func timebase(b string) float64 {
	base, err := strconv.Atoi(strings.TrimPrefix(b, "1/"))
	if err != nil {
		log.Fatal(err)
	}
	return float64(base)
}
