package media

import (
	"log"
	"path/filepath"

	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/ffmeta"
	"github.com/ohzqq/avtools/ffprobe"
)

type Media struct {
	Input string
	files RelatedFiles
	Meta
}

type Meta struct {
	*ffmeta.FFmeta
}

type RelatedFiles map[string]string

func NewMedia(i string) *Media {
	media := Media{
		Input: i,
		files: make(RelatedFiles),
	}
	media.Meta = media.ReadEmbeddedMeta()
	return &media
}

func (m *Media) AddFile(name, path string) *Media {
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	m.files[name] = abs
	return m
}

func (m Media) HasFFmeta() bool {
	if ff, ok := m.files["ffmeta"]; ok && ff != "" {
		return true
	}
	return false
}

func (m *Media) SetFFmeta(ff string) *Media {
	m.files["ffmeta"] = ff
	return m
}

func (m Media) HasCue() bool {
	if c, ok := m.files["cue"]; ok && c != "" {
		return true
	}
	return false
}

func (m *Media) SetCue(c string) *Media {
	m.files["cue"] = c
	return m
}

func (m *Media) SetMeta() *Media {
	if m.HasFFmeta() {
		m.Meta = m.ReadFFmeta()
	}

	if m.HasCue() {
		m.Meta.SetChapters(m.ReadCueSheet())
	}

	return m
}

func (m *Media) ReadEmbeddedMeta() Meta {
	probe := ffprobe.New()
	probe.Input(m.Input).
		Stream("a").
		FormatEntry("filename", "start_time", "duration", "size", "bit_rate").
		StreamEntry("codec_type", "codec_name").
		Entry("format_tags").
		ShowChapters().
		Json()

	data := probe.Run()
	return Meta{FFmeta: ffmeta.LoadJson(data)}
}

func (m *Media) ReadCueSheet() chap.Chapters {
	var ch chap.Chapters
	if m.HasCue() {
		ch = chap.NewChapters().FromCue(m.files["cue"])
	}
	return ch
}

func (m *Media) ReadFFmeta() Meta {
	var ff *ffmeta.FFmeta
	if m.HasFFmeta() {
		ff = ffmeta.LoadIni(m.files["ffmeta"])
	}
	return Meta{FFmeta: ff}
}
