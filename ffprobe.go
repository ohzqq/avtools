package avtools

import (
	"log"
	"strings"

	"github.com/ohzqq/avtools/timestamp"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type chapterEntry struct {
	Base  string            `json:"time_base",ini:"timebase"`
	Start int               `json:"start",ini:"start"`
	End   int               `json:"end",ini:"end"`
	Title string            `json:"title", ini:"title"`
	Tags  map[string]string `json:"tags"`
}

var probeArgs = []ffmpeg.KwArgs{
	ffmpeg.KwArgs{"show_chapters": ""},
	ffmpeg.KwArgs{"select_streams": "a"},
	ffmpeg.KwArgs{"show_entries": "stream=codec_type,codec_name:format=filename, start_time, duration, size, bit_rate:format_tags"},
	ffmpeg.KwArgs{"of": "json"},
}

func Probe(input string) []byte {
	args := ffmpeg.MergeKwArgs(probeArgs)
	info, err := ffmpeg.ProbeWithTimeoutExec(input, 0, args)
	if err != nil {
		log.Fatal(err)
	}

	return []byte(info)
}

func (c chapterEntry) ToChapter() *Chapter {
	base := strings.Split(c.Base, "/")[1]
	b := int(timestamp.StringToFloat(base))
	ch := NewChapter(c.Start, c.End, b)

	if t, ok := c.Tags["title"]; ok {
		ch.Title = t
	}
	return ch
}
