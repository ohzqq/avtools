package media

import (
	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/ffmeta"
	"github.com/ohzqq/avtools/ffprobe"
)

type Media struct {
	input string
	files RelatedFiles
	meta  *ffmeta.FFmeta
}

type RelatedFiles map[string]string

func New(i string) *Media {
	media := Media{
		input: i,
		files: make(RelatedFiles),
	}
	media.ReadEmbeddedMeta()
	return &media
}

func (m *Media) ReadEmbeddedMeta() *Media {
	probe := ffprobe.New()
	probe.Input(m.input).
		Stream("a").
		FormatEntry("filename", "start_time", "duration", "size", "bit_rate").
		StreamEntry("codec_type", "codec_name").
		Entry("format_tags").
		ShowChapters().
		Json()

	data := probe.Run()
	m.meta = ffmeta.LoadJson(data)

	return m
}

func (m *Media) ReadCueSheet(c string) *Media {
	ch := chap.NewChapters().FromCue(c)
	m.meta.SetChapters(ch)
	return m
}

func (m *Media) ReadFFmeta(ff string) *Media {
	m.meta = ffmeta.LoadIni(ff)
	return m
}
