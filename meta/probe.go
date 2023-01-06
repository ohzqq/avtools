package meta

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type ProbeMeta struct {
	Streams  []*ProbeStream `json:"streams"`
	Format   ProbeFormat    `json:"format"`
	Chapters []ProbeChapter `json:"chapters"`
}

type ProbeStream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type ProbeFormat struct {
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
}

type ProbeChapter struct {
	Base         string            `json:"time_base"`
	StartTime    int               `json:"start"`
	EndTime      int               `json:"end"`
	ChapterTitle string            `json:"title"`
	Tags         map[string]string `json:"tags"`
}

var probeArgs = []ffmpeg.KwArgs{
	ffmpeg.KwArgs{"show_chapters": ""},
	ffmpeg.KwArgs{"select_streams": "a"},
	ffmpeg.KwArgs{"show_entries": "stream=codec_type,codec_name:format=filename, start_time, duration, size, bit_rate:format_tags"},
	ffmpeg.KwArgs{"of": "json"},
}

func FFProbe(input string) ProbeMeta {
	args := ffmpeg.MergeKwArgs(probeArgs)
	info, err := ffmpeg.ProbeWithTimeoutExec(input, 0, args)
	if err != nil {
		log.Fatal(err)
	}

	data := []byte(info)

	var meta ProbeMeta
	err = json.Unmarshal(data, &meta)
	if err != nil {
		log.Fatal(err)
	}

	return meta
}

func UnmarshalJSON(d []byte) ProbeMeta {
	var meta ProbeMeta
	err := json.Unmarshal(d, &meta)
	if err != nil {
		log.Fatal(err)
	}
	return meta
}

func (c ProbeChapter) Title() string {
	if t, ok := c.Tags["title"]; ok {
		return t
	}
	return c.ChapterTitle
}

func (c ProbeChapter) Start() int {
	return c.StartTime
}

func (c ProbeChapter) End() int {
	return c.EndTime
}

func (c ProbeChapter) Timebase() float64 {
	if tb := c.Base; tb != "" {
		c.Base = strings.TrimPrefix(tb, "1/")
	}
	baseFloat, _ := strconv.ParseFloat(c.Base, 64)
	return baseFloat
}
