package av

import (
	"log"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/cue"
	"github.com/ohzqq/avtools/probe"
	"github.com/ohzqq/fidi"
)

type Media struct {
	*avtools.Media
	Profile     string
	HasCover    bool
	MetaChanged bool
}

func New(input string) *Media {
	media := avtools.Media{
		File: fidi.NewFile(input),
	}
	return &Media{
		Media: &media,
	}
}

func (m *Media) Cue(file string) *Media {
	meta, err := cue.Load(file)
	if err != nil {
		log.Fatal(err)
	}
	m.Chaps = avtools.NewChapters(meta.Chapters())
	return m
}

func (m *Media) Probe() *Media {
	meta, err := probe.Load(m.Path())
	if err != nil {
		log.Fatal(err)
	}
	m.Chaps = avtools.NewChapters(meta.Chapters())
	m.SetStreams(meta.Streams())
	m.SetTags(meta.Tags())
	return m
}
