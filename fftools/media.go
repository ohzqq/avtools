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
