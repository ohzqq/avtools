package avtools

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/ohzqq/avtools/timestamp"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type ProbeMeta struct {
	name        string
	Streams     []*StreamEntry `json:"streams"`
	FormatEntry `json:"format"`
	Chaps       []map[string]json.RawMessage `json:"chapters"`
	Chapters    []*Chapter
}

type StreamEntry struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type FormatEntry struct {
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
}

type ChapterEntry struct {
	Base         string            `json:"time_base",ini:"timebase"`
	StartTime    int               `json:"start",ini:"start"`
	EndTime      int               `json:"end",ini:"end"`
	ChapterTitle string            `json:"title", ini:"title"`
	Tags         map[string]string `json:"tags"`
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

func ReadEmbeddedMeta(input string) *ProbeMeta {
	info := Probe(input)

	var meta ProbeMeta
	err := json.Unmarshal(info, &meta)
	if err != nil {
		log.Fatal(err)
	}

	var chapters []*Chapter
	for _, chapter := range meta.Chaps {
		for key, val := range chapter {
			var start int
			var end int
			var base int
			var title string
			switch key {
			case "start":
				err := json.Unmarshal(val, &start)
				if err != nil {
					log.Fatal(err)
				}
			case "end":
				err := json.Unmarshal(val, &end)
				if err != nil {
					log.Fatal(err)
				}
			case "time_base":
				var b string
				err := json.Unmarshal(val, &b)
				if err != nil {
					log.Fatal(err)
				}
				b = strings.Split(b, "/")[1]
				base = int(timestamp.StringToFloat(b))
			case "tags":
				tags := make(map[string]string)
				err := json.Unmarshal(val, &tags)
				if err != nil {
					log.Fatal(err)
				}
				title = tags["title"]
			}
			ch := NewChapter(start, end, base)
			ch.Title = title
			chapters = append(chapters, ch)
		}
	}
	meta.Chapters = chapters

	return &meta
}
