package media

import (
	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/meta"
)

type Media struct {
	*avtools.Media
}

func New(input string) *Media {
	m := avtools.NewMedia(input)
	return &Media{Media: m}
}

func (m *Media) LoadIni(name string) *Media {
	ini := meta.LoadIni(name)
	m.SetMeta(ini)
	return m
}

func (m *Media) LoadCue(name string) *Media {
	cue := meta.LoadCueSheet(name)
	m.SetMeta(cue)
	return m
}

func (m *Media) Probe() *Media {
	p := meta.FFProbe(m.Filename)
	m.SetMeta(p)
	return m
}
