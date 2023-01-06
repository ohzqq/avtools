package meta

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/ohzqq/avtools"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type ProbeMeta struct {
	StreamEntry  []map[string]string `json:"streams"`
	Format       ProbeFormat         `json:"format"`
	ChapterEntry []ProbeChapter      `json:"chapters"`
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
	StartTime    float64           `json:"start"`
	EndTime      float64           `json:"end"`
	ChapterTitle string            `json:"title"`
	Tags         map[string]string `json:"tags"`
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

func (m ProbeMeta) Chapters() []avtools.ChapterMeta {
	var chaps []avtools.ChapterMeta
	for _, ch := range m.ChapterEntry {
		chaps = append(chaps, ch)
	}
	return chaps
}

func (m ProbeMeta) Streams() []map[string]string {
	return m.StreamEntry
}

func (m ProbeMeta) Tags() map[string]string {
	m.Format.Tags["filename"] = m.Format.Filename
	m.Format.Tags["duration"] = m.Format.Dur
	m.Format.Tags["size"] = m.Format.Size
	m.Format.Tags["bit_rate"] = m.Format.BitRate
	return m.Format.Tags
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

func (c ProbeChapter) Start() float64 {
	return c.StartTime
}

func (c ProbeChapter) End() float64 {
	return c.EndTime
}

func (c ProbeChapter) Timebase() float64 {
	if tb := c.Base; tb != "" {
		c.Base = strings.TrimPrefix(tb, "1/")
	}
	baseFloat, _ := strconv.ParseFloat(c.Base, 64)
	return baseFloat
}

var probeArgs = []ffmpeg.KwArgs{
	ffmpeg.KwArgs{"show_chapters": ""},
	ffmpeg.KwArgs{"select_streams": "a"},
	ffmpeg.KwArgs{"show_entries": "stream=codec_type,codec_name:format=filename, start_time, duration, size, bit_rate:format_tags"},
	ffmpeg.KwArgs{"of": "json"},
}
