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

func (m *Media) SetMeta(meta *MediaMeta) *Media {
	m.Meta = meta
	return m
}

func (m *Media) HasChapters() bool {
	if m.Meta != nil {
		if *m.Meta.Chapters != nil {
			return true
		}
	}
	return false
}

func (m *Media) SetChapters(ch *Chapters) {
	if !m.HasChapters() {
		m.WithMeta()
	}
	m.Meta.Chapters = ch
}

func (m *Media) HasVideo() bool {
	if !m.HasStreams() {
		m.WithMeta()
	}
	for _, stream := range *m.Meta.Streams {
		if stream.CodecType == "video" {
			return true
		}
	}
	return false
}

func (m *Media) HasAudio() bool {
	if !m.HasStreams() {
		m.WithMeta()
	}
	for _, stream := range *m.Meta.Streams {
		if stream.CodecType == "audio" {
			return true
		}
	}
	return false
}

func (m *Media) VideoCodec() string {
	if !m.HasStreams() {
		m.WithMeta()
	}
	for _, stream := range *m.Meta.Streams {
		if stream.CodecType == "video" {
			return stream.CodecName
		}
	}
	return ""
}

func (m *Media) AudioCodec() string {
	if !m.HasStreams() {
		m.WithMeta()
	}
	for _, stream := range *m.Meta.Streams {
		if stream.CodecType == "audio" {
			return stream.CodecName
		}
	}
	return ""
}

func (m *Media) HasMeta() bool {
	if m.Meta != nil {
		return true
	}
	return false
}

func (m *Media) HasStreams() bool {
	if m.HasMeta() {
		if m.Meta.Streams != nil {
			return true
		}
	}
	return false
}

func (m *Media) hasFormat() bool {
	if m.Meta != nil {
		if m.Meta.Format != nil {
			return true
		}
	}
	return false
}

func (m *Media) Duration() string {
	if !m.hasFormat() {
		m.WithMeta()
	}
	return secsToHHMMSS(m.Meta.Format.Duration)
}

