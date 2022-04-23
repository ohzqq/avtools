package fftools

import (
	"path/filepath"
	"log"
	"fmt"
	//"strconv"
	//"strings"
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

func (m *Media) SetChapters(ch *Chapters) {
	m.Meta.Chapters = ch
}

func (m *Media) HasVideo() bool {
	for _, stream := range *m.Meta.Streams {
		if stream.CodecType == "video" {
			return true
		}
	}
	return false
}

func (m *Media) HasAudio() bool {
	for _, stream := range *m.Meta.Streams {
		if stream.CodecType == "audio" {
			return true
		}
	}
	return false
}

func (m *Media) VideoCodec() string {
	for _, stream := range *m.Meta.Streams {
		if stream.CodecType == "video" {
			return stream.CodecName
		}
	}
	return ""
}

func (m *Media) AudioCodec() string {
	for _, stream := range *m.Meta.Streams {
		if stream.CodecType == "audio" {
			return stream.CodecName
		}
	}
	return ""
}

func (m *Media) Duration() string {
	return secsToHHMMSS(m.Meta.Format.Duration)
}

