package chap

import (
	"fmt"
)

type Chapter struct {
	start    Time
	end      Time
	Timebase float64
	Title    string
}

type ChapMeta interface {
	Start() int
	End() int
	Title() string
	Timebase() float64
}

func NewChapter() *Chapter {
	return &Chapter{Timebase: 1}
}

func (ch *Chapter) SetMeta(c ChapMeta) *Chapter {
	ch.start = NewChapterTime(c.Start())
	ch.end = NewChapterTime(c.End())
	ch.Title = c.Title()
	ch.Timebase = c.Timebase()
	return ch
}

func (ch Chapter) Start() Time {
	if t := ch.Timebase; t != 1 {
		ch.start.base = t
	}
	return ch.start
}

func (ch Chapter) End() Time {
	if t := ch.Timebase; t != 1 {
		ch.start.base = t
	}
	return ch.end
}

func (ch Chapter) Dur() (Time, error) {
	if ch.end.time == 0 {
		return ch.end, fmt.Errorf("end time is needed to calculate duration")
	}
	t := ch.end.time - ch.start.time
	return Time{time: t, base: ch.Timebase}, nil
}

func (ch *Chapter) SetTitle(t string) *Chapter {
	ch.Title = t
	return ch
}

func (ch *Chapter) SetStart(t Time) *Chapter {
	ch.start = t
	return ch
}

func (ch *Chapter) SetEnd(t Time) *Chapter {
	ch.end = t
	return ch
}

func (ch *Chapter) SetTimebase(t float64) *Chapter {
	ch.Timebase = t
	return ch
}

const cueTmpl = `FILE "{{.File}}" {{.Ext}}
{{- range $index, $ch := .Each}}
TRACK {{$index}} AUDIO
  TITLE "{{$ch.Title}}"
  INDEX 01 {{$ch.Start.MMSS}}{{end -}}`
