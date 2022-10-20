package media

import (
	"log"
	"path/filepath"

	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/ffmeta"
	"github.com/ohzqq/avtools/ffprobe"
)

type Media struct {
	Input FileFormat
	Files RelatedFiles
	Meta
}

type Meta struct {
	*ffmeta.FFmeta
}

func NewMedia(i string) *Media {
	media := Media{
		Input: FileFormat(i),
		Files: make(RelatedFiles),
	}
	media.Meta = media.ReadEmbeddedMeta()
	return &media
}

func (m *Media) AddFile(name, path string) *Media {
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	m.Files[name] = FileFormat(abs)
	return m
}

func (m Media) HasFFmeta() bool {
	return m.Files.Has("ffmeta")
}

func (m *Media) SetFFmeta(ff string) *Media {
	m.Files["ffmeta"] = FileFormat(ff)
	return m
}

func (m Media) HasCue() bool {
	return m.Files.Has("cue")
}

func (m *Media) SetCue(c string) *Media {
	m.Files["cue"] = FileFormat(c)
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
	probe.Input(m.Input.String()).
		Stream("a").
		FormatEntry("filename", "start_time", "duration", "size", "bit_rate").
		StreamEntry("codec_type", "codec_name").
		Entry("format_tags").
		ShowChapters().
		Json()

	data, err := probe.Run()
	if err != nil {
		log.Fatal(err)
	}
	return Meta{FFmeta: ffmeta.LoadJson(data)}
}

func (m *Media) ReadCueSheet() chap.Chapters {
	var ch chap.Chapters
	if m.HasCue() {
		ch = chap.NewChapters().FromCue(m.Files.Get("cue"))
	}
	return ch
}

func (m *Media) ReadFFmeta() Meta {
	var ff *ffmeta.FFmeta
	if m.HasFFmeta() {
		ff = ffmeta.LoadIni(m.Files.Get("ffmeta"))
	}
	return Meta{FFmeta: ff}
}

func (m Media) AudioCodec() string {
	if m.HasAudio() {
		a := m.AudioStreams()
		return a[0].CodecType
	}
	return ""
}
