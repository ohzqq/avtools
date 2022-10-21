package chap

import (
	"fmt"
)

type Chapter struct {
	start Time
	end   Time
	base  timebase
	Title string
}

type ChapMeta interface {
	Start() int
	End() int
	Title() string
	Timebase() float64
}

func NewChapter() *Chapter {
	return &Chapter{base: timebase(1)}
}

func (ch *Chapter) SetMeta(c ChapMeta) *Chapter {
	ch.start = NewChapterTime(c.Start())
	ch.end = NewChapterTime(c.End())
	ch.Title = c.Title()
	ch.base = timebase(c.Timebase())
	return ch
}

func (ch Chapter) Start() Time {
	if t := ch.base; t != 1 {
		ch.start.base = float64(t)
	}
	return ch.start
}

func (ch Chapter) End() Time {
	if t := ch.base; t != 1 {
		ch.end.base = float64(t)
	}
	return ch.end
}

func (ch Chapter) Timebase() timebase {
	return ch.base
}

func (ch Chapter) Dur() (Time, error) {
	if ch.end.time == 0 {
		return ch.end, fmt.Errorf("end time is needed to calculate duration")
	}
	t := ch.end.time - ch.start.time
	return Time{time: t, base: float64(ch.base)}, nil
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
	ch.base = timebase(t)
	return ch
}

const cueTmpl = `FILE "{{.File}}" {{.Ext}}
{{- range $index, $ch := .Each}}
TRACK {{$index}} AUDIO
  TITLE "{{$ch.Title}}"
  INDEX 01 {{$ch.Start.MMSS}}{{end -}}`
