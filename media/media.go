package media

import (
	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/ffmeta"
	"github.com/ohzqq/avtools/ffprobe"
)

type Media struct {
	input       string
	files       RelatedFiles
	meta        *ffmeta.FFmeta
	cueChapters chap.Chapters
}

type RelatedFiles map[string]string

func New(i string) *Media {
	media := Media{
		input: i,
		files: make(RelatedFiles),
	}
	return &media
}

func (m Media) Meta() *ffmeta.FFmeta {
	meta := m.ReadEmbeddedMeta()

	if m.HasFFmeta() {
		meta = m.ReadFFmeta()
	}

	if m.HasCue() {
		meta.SetChapters(m.ReadCueSheet())
	}

	return meta
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

func (m *Media) ReadEmbeddedMeta() *ffmeta.FFmeta {
	probe := ffprobe.New()
	probe.Input(m.input).
		Stream("a").
		FormatEntry("filename", "start_time", "duration", "size", "bit_rate").
		StreamEntry("codec_type", "codec_name").
		Entry("format_tags").
		ShowChapters().
		Json()

	data := probe.Run()
	return ffmeta.LoadJson(data)
}

func (m *Media) ReadCueSheet() chap.Chapters {
	var ch chap.Chapters
	if m.HasCue() {
		ch = chap.NewChapters().FromCue(m.files["cue"])
	}
	return ch
}

func (m *Media) ReadFFmeta() *ffmeta.FFmeta {
	var ff *ffmeta.FFmeta
	if m.HasFFmeta() {
		ff = ffmeta.LoadIni(m.files["ffmeta"])
	}
	return ff
}
