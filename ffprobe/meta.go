package ffprobe

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type Meta struct {
	name    string
	Streams []*Stream `json:"streams"`
	Format  `json:"format"`
	Chaps   []Chapter `json:"chapters"`
}

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Filename string            `json:"filename"`
	Dur      string            `json:"duration"`
	Size     string            `json:"size"`
	BitRate  string            `json:"bit_rate"`
	Tags     map[string]string `json:"tags"`
}

type Chapter struct {
	Base         string            `json:"time_base",ini:"timebase"`
	StartTime    int               `json:"start",ini:"start"`
	EndTime      int               `json:"end",ini:"end"`
	ChapterTitle string            `json:"title", ini:"title"`
	Tags         map[string]string `json:"tags"`
}

func UnmarshalJSON(d []byte) Meta {
	var meta Meta
	err := json.Unmarshal(d, &meta)
	if err != nil {
		log.Fatal(err)
	}
	return meta
}

func (c Chapter) Title() string {
	if t, ok := c.Tags["title"]; ok {
		return t
	}
	return c.ChapterTitle
}

func (c Chapter) Start() int {
	return c.StartTime
}

func (c Chapter) End() int {
	return c.EndTime
}

func (c Chapter) Timebase() float64 {
	if tb := c.Base; tb != "" {
		c.Base = strings.TrimPrefix(tb, "1/")
	}
	baseFloat, _ := strconv.ParseFloat(c.Base, 64)
	return baseFloat
}
