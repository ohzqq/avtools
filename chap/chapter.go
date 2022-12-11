package chap

import (
	"fmt"

	"github.com/ohzqq/avtools/timestamp"
)

type Chapter struct {
	start timestamp.Time
	end   timestamp.Time
	base  timestamp.Timebase
	Title string
}

type ChapMeta interface {
	Start() int
	End() int
	Title() string
	Timebase() float64
}

func NewChapter() *Chapter {
	return &Chapter{base: timestamp.Timebase(1)}
}

func (ch *Chapter) SetMeta(c ChapMeta) *Chapter {
	ch.start = timestamp.NewChapterTime(c.Start())
	ch.end = timestamp.NewChapterTime(c.End())
	ch.Title = c.Title()
	ch.base = timestamp.Timebase(c.Timebase())
	return ch
}

func (ch Chapter) Start() timestamp.Time {
	if t := ch.base; t != 1 {
		ch.start.Base = float64(t)
	}
	return ch.start
}

func (ch Chapter) End() timestamp.Time {
	if t := ch.base; t != 1 {
		ch.end.Base = float64(t)
	}
	return ch.end
}

func (ch Chapter) Timebase() timestamp.Timebase {
	return ch.base
}

func (ch Chapter) Dur() (timestamp.Time, error) {
	if ch.end.Time == 0 {
		return ch.end, fmt.Errorf("end time is needed to calculate duration")
	}
	t := ch.end.Time - ch.start.Time
	return timestamp.Time{Time: t, Base: float64(ch.base)}, nil
}

func (ch *Chapter) SetTitle(t string) *Chapter {
	ch.Title = t
	return ch
}

func (ch *Chapter) SetStart(t timestamp.Time) *Chapter {
	ch.start = t
	return ch
}

func (ch *Chapter) SetEnd(t timestamp.Time) *Chapter {
	ch.end = t
	return ch
}

func (ch *Chapter) SetTimebase(t float64) *Chapter {
	ch.base = timestamp.Timebase(t)
	return ch
}

const cueTmpl = `FILE "{{.File}}" {{.Ext}}
{{- range $index, $ch := .Each}}
TRACK {{$index}} AUDIO
  TITLE "{{$ch.Title}}"
  INDEX 01 {{$ch.Start.MMSS}}{{end -}}`
