package tool

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/ohzqq/avtools/chap"
)

const ffProbeMeta = `format=filename,start_time,duration,size,bit_rate:stream=codec_type,codec_name:format_tags`

type MediaMeta struct {
	data     []byte
	Chapters Chapters
	Ch       chap.Chapters
	Streams  []*Stream
	Format   *Format
}

func (m *MediaMeta) MarshalTo(format string) *FileFormat {
	return NewFormat(format).SetMeta(m)
}

func (m *MediaMeta) SetChapters(ch Chapters) {
	m.Chapters = ch
}

func (m *MediaMeta) SetTags(tags map[string]string) {
	m.Format.Tags = tags
}

func (m *MediaMeta) GetTag(tag string) string {
	if t := m.Format.Tags[tag]; t != "" {
		return t
	}
	return ""
}

func (m MediaMeta) Tags() map[string]string {
	return m.Format.Tags
}

func (m *MediaMeta) LastChapterEnd() {
	if m.Format.Duration != "" && m.HasChapters() {
		lastCh := m.Chapters[len(m.Chapters)-1]
		lastCh.End = m.Format.DurationSecs(lastCh.TimebaseFloat())
	}
}

func (m *MediaMeta) HasChapters() bool {
	return len(m.Chapters) > 0
}

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Filename string
	Duration string
	Size     string
	BitRate  string `json:"bit_rate"`
	Tags     map[string]string
}

func (f Format) Ext() string {
	if f.Filename != "" {
		return strings.TrimPrefix(path.Ext(f.Filename), ".")
	}
	return ""
}

func (f Format) DurationSecs(timebase float64) int {
	seconds := decimalSecsToFloat(f.Duration) * timebase
	return int(seconds)
}

type Chapters []*Chapter

type Chapter struct {
	Timebase string            `json:"time_base",ini:"timebase"`
	Start    int               `json:"start",ini:"start"`
	End      int               `json:"end",ini:"end"`
	Tags     map[string]string `json:"tags"`
	Title    string            `ini:"title"`
}

func (c *Chapter) TimeBase() string {
	if c.Timebase == "" {
		return "1/1000"
	}
	return c.Timebase
}

func (c *Chapter) StartToIntString() string {
	result := float64(c.Start) * c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 0, 64)
}

func (c *Chapter) CueStamp() string {
	sec := float64(c.Start) / c.TimebaseFloat()
	m := int(sec) / 60
	s := int(sec) % 60
	return fmt.Sprintf("%02d:%02d:00", m, s)
}

func (c *Chapter) StartToSeconds() string {
	if c.Start == 0 {
		return "0"
	}
	result := float64(c.Start) / c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 3, 64)
}

func (c *Chapter) EndToIntString() string {
	result := float64(c.End) * c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 0, 64)
}

func (c *Chapter) EndToSeconds() string {
	if c.End == 0 {
		return "0"
	}
	result := float64(c.End) / c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 3, 64)
}

func (c Chapter) TimebaseFloat() float64 {
	base := "1000"
	if tb := c.Timebase; tb != "" {
		base = strings.ReplaceAll(tb, "1/", "")
	}
	baseFloat, _ := strconv.ParseFloat(base, 64)
	return baseFloat
}
