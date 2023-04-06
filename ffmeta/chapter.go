package ffmeta

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ohzqq/dur"
	"gopkg.in/ini.v1"
)

type Chapter struct {
	*ini.Section
}

func (ch Chapter) Start() time.Duration {
	if ch.HasKey("START") {
		t, err := dur.Parse(ch.Key("START").String())
		if err != nil {
			log.Fatal(err)
		}
		return calculateSecs(t.Dur.Seconds(), ch.GetTag("TIMEBASE"))
	}
	return time.Duration(0)
}

func (ch Chapter) End() time.Duration {
	if ch.HasKey("START") {
		t, err := dur.Parse(ch.Key("END").String())
		if err != nil {
			log.Fatal(err)
		}
		return calculateSecs(t.Dur.Seconds(), ch.GetTag("TIMEBASE"))
	}
	return time.Duration(0)
}

func (ch Chapter) Title() string {
	if ch.HasKey("title") {
		return ch.Key("title").String()
	}
	return ""
}

func (ch Chapter) GetTag(t string) string {
	if ch.HasKey(t) {
		return ch.Key(t).String()
	}
	return ""
}

func (ch Chapter) Tags() map[string]string {
	tags := make(map[string]string)
	for _, k := range ch.Keys() {
		switch n := k.Name(); n {
		case "TIMEBASE", "title", "START", "END":
		default:
			tags[n] = k.String()
		}
	}
	return tags
}

func calculateSecs(num float64, base string) time.Duration {
	b := timebase(base)
	t := num / b * float64(time.Second)
	return time.Duration(t)
}

func timebase(b string) float64 {
	if b == "" {
		b = "1/1"
	}
	base, err := strconv.Atoi(strings.TrimPrefix(b, "1/"))
	if err != nil {
		log.Fatal(err)
	}
	return float64(base)
}
