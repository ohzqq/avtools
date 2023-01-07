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
	StreamEntry  []map[string]any `json:"streams"`
	Format       ProbeFormat      `json:"format"`
	ChapterEntry []ProbeChapter   `json:"chapters"`
}

type ProbeFormat struct {
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
}

type ProbeChapter struct {
	Base      string `json:"time_base"`
	StartTime string `json:"start_time"`
	//StartTime    float64           `json:"start"`
	//EndTime      float64           `json:"end"`
	EndTime      string            `json:"end_time"`
	ChapterTitle string            `json:"title"`
	Tags         map[string]string `json:"tags"`
}

func FFProbe(input string) ProbeMeta {
	args := ffmpeg.MergeKwArgs(probeArgs)
	info, err := ffmpeg.ProbeWithTimeoutExec(input, 0, args)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(info)

	data := []byte(info)

	var meta ProbeMeta
	err = json.Unmarshal(data, &meta)
	if err != nil {
		log.Fatal(err)
	}

	return meta
}

func (ch ProbeChapter) Chap() *avtools.Chapter {
	return &avtools.Chapter{
		StartTime: avtools.Timestamp(avtools.ParseDuration(ch.StartTime + "s")),
		EndTime:   avtools.Timestamp(avtools.ParseDuration(ch.EndTime + "s")),
		ChTitle:   ch.Title(),
	}
}

func (m ProbeMeta) Chapters() []*avtools.Chapter {
	var ch []*avtools.Chapter
	for _, c := range m.ChapterEntry {
		ch = append(ch, c.Chap())
	}
	return ch
}

func (m ProbeMeta) Streams() []map[string]string {
	var streams []map[string]string
	for _, stream := range m.StreamEntry {
		meta := make(map[string]string)
		for key, raw := range stream {
			switch val := raw.(type) {
			case float64:
				meta[key] = strconv.Itoa(int(val))
			case string:
				meta[key] = val
			case map[string]any:
				if key == "disposition" {
					if val["attached_pic"].(float64) == 0 {
						meta["cover"] = "false"
					} else {
						meta["cover"] = "true"
					}
				}
			}
		}
		streams = append(streams, meta)
	}
	return streams
}

func (m ProbeMeta) Tags() map[string]string {
	m.Format.Tags["filename"] = m.Format.Filename
	m.Format.Tags["duration"] = m.Format.Dur
	m.Format.Tags["size"] = m.Format.Size
	m.Format.Tags["bit_rate"] = m.Format.BitRate
	return m.Format.Tags
}

func (c ProbeChapter) Title() string {
	if t, ok := c.Tags["title"]; ok {
		return t
	}
	return c.ChapterTitle
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
	//ffmpeg.KwArgs{"select_streams": "a"},
	ffmpeg.KwArgs{"show_entries": "stream:format=filename, start_time, duration, size, bit_rate:format_tags"},
	ffmpeg.KwArgs{"of": "json"},
}
