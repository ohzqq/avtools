package av

import (
	"log"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/cue"
	"github.com/ohzqq/avtools/ffmeta"
	"github.com/ohzqq/avtools/probe"
	"github.com/ohzqq/dur"
	"github.com/ohzqq/fidi"
)

type Media struct {
	*avtools.Media
	Profile     string
	MetaChanged bool
}

func New(input string) *Media {
	media := avtools.Media{
		File: fidi.NewFile(input),
	}
	m := &Media{
		Media: &media,
	}

	meta, err := probe.Load(m.Path())
	if err != nil {
		log.Fatal(err)
	}
	m.SetChapters(meta.Chapters())
	m.SetStreams(meta.Streams())
	m.SetTags(meta.Tags())

	return m
}

func (m *Media) Cue(file string) *Media {
	meta, err := cue.Load(file)
	if err != nil {
		log.Fatal(err)
	}
	m.SetChapters(meta.Chapters())
	m.Media.Cue = meta.Source()

	if d := m.GetTag("duration"); d != "" {
		last := m.Chaps[len(m.Chaps)-1]
		to, err := dur.Parse(d)
		if err != nil {
			log.Fatal(err)
		}
		last.EndStamp = to
	}
	return m
}

func (m *Media) FFMeta(file string) *Media {
	meta, err := ffmeta.Load(m.Path())
	if err != nil {
		log.Fatal(err)
	}
	m.Media.Ini = meta.Source()
	m.SetChapters(meta.Chapters())
	m.SetStreams(meta.Streams())
	m.SetTags(meta.Tags())
	return m
}

func (m *Media) SetChapters(chaps []avtools.ChapterMeta) {
	m.Media.Chaps = avtools.NewChapters(chaps)
}

func (m *Media) SetTags(tags map[string]string) {
	m.Media.Tagz = tags
}

func (m *Media) SetStreams(streams []map[string]string) {
	for _, stream := range streams {
		s := avtools.Stream{}
		for key, val := range stream {
			switch key {
			case "codec_type":
				s.CodecType = val
			case "codec_name":
				s.CodecName = val
			case "index":
				s.Index = val
			case "cover":
				if val == "true" {
					s.IsCover = true
					m.HasCover = true
				}
			}
		}
		m.Media.Streamz = append(m.Media.Streamz, s)
	}
}
