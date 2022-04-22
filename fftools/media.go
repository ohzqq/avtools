package fftools

import (
	"path/filepath"
	"log"
	"fmt"
	//"os"
)
var _ = fmt.Printf

type Media struct {
	File string
	Path string
	Dir string
	Ext string
	FFmeta string
	Cue string
	Cover string
	Meta *MediaMeta
}

type MediaMeta struct {
	Chapters *Chapters
	Streams *Streams
	Format *Format
	Tags *Tags
}

type Streams []*Stream

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Filename string
	StartTime string `json:"start_time"`
	Duration string
	Size string
	BitRate string `json:"bit_rate"`
}

type Tags struct {
	Title string `json:"title"`
	Artist string `json:"artist"`
	Composer string `json:"composer"`
	Album string `json:"album"`
	Comment string `json:"comment"`
	Genre string `json:"genre"`
}

type Chapters []*Chapter

type Chapter struct {
	Timebase string `json:"time_base"`
	Start string `json:"start_time"`
	End string `json:"end_time"`
	Title string
}

func (c *Chapter) timeBaseFloat() float64 {
	tb := strings.ReplaceAll(c.Timebase, "1/", "")
	baseint, _ := strconv.ParseFloat(tb, 64)
	return baseint
}

func (c *Chapter) toSeconds() () {
	tb := c.timeBaseFloat()
	ss, _ := strconv.ParseFloat(c.Start, 64)
	to, _ := strconv.ParseFloat(c.End, 64)
	c.Start = strconv.FormatFloat(ss / tb, 'f', 6, 64)
	c.End = strconv.FormatFloat(to / tb, 'f', 6, 64)
}

func NewMedia(input string) *Media {
	media := new(Media)

	abs, err := filepath.Abs(input)
	if err != nil { log.Fatal(err) }

	media.Path = abs
	media.File = filepath.Base(input)
	media.Dir = filepath.Dir(input)
	media.Ext = filepath.Ext(input)

	return media
}

func (m *Media) Cut(ss, to string, no int) {
	count := fmt.Sprintf("%06d", no + 1)
	cmd := NewCmd().In(m)
	timestamps := make(map[string]string)
	if ss != "" {
		timestamps["ss"] = ss
	}
	if to != "" {
		timestamps["to"] = to
	}
	cmd.Args().PostInput(timestamps).Out("tmp" + count).Ext(m.Ext)
	cmd.Run()
}

func (m *Media) WithMeta() *Media {
	m.Meta = m.ReadMeta()
	return m
}

func (m *Media) ReadMeta() *MediaMeta {
	return ReadEmbeddedMeta(m.Path)
}

func (m *Media) WriteMeta() {
	WriteFFmetadata(m.Path)
}

func (m *Media) HasChapters() bool {
	if m.Meta != nil {
		if len(*m.Meta.Chapters) != 0 {
			return true
		}
	}
	return false
}
