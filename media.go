package avtools

import (
	"fmt"
)

type Media struct {
	Filename string
	Dur      string
	Size     string
	BitRate  string
	Tags     map[string]string
	Chapters []*Chapter
	Streams  []map[string]string
}

type Chapter struct {
	start Time
	end   Time
	base  Timebase
	title string
}

type ChapterMeta interface {
	Start() float64
	End() float64
	Timebase() float64
	Title() string
}

type Meta interface {
	Chapters() []ChapterMeta
	Tags() map[string]string
	Streams() []map[string]string
}

func NewMedia(input string) *Media {
	media := Media{
		Filename: input,
	}
	return &media
}

func (m *Media) SetMeta(meta Meta) *Media {
	if tags := meta.Tags(); tags != nil {
		m.Tags = tags
	}
	if chaps := meta.Chapters(); len(chaps) > 0 {
		for _, chap := range chaps {
			m.Chapters = append(m.Chapters, NewChapter(chap))
		}
	}
	//if streams := meta.Streams(); len(streams) > 0 {
	m.Streams = meta.Streams()
	//}
	return m
}

func NewChapter(chap ChapterMeta) *Chapter {
	return &Chapter{
		title: chap.Title(),
		start: Timestamp(chap.Start(), chap.Timebase()),
		end:   Timestamp(chap.End(), chap.Timebase()),
		base:  Timebase(chap.Timebase()),
	}
}

func (ch Chapter) Start() Time {
	return ch.start
}

func (ch Chapter) End() Time {
	return ch.end
}

func (ch Chapter) Timebase() Timebase {
	return ch.base
}

func (ch Chapter) Title() string {
	return ch.title
}

func (ch Chapter) Dur() (Time, error) {
	if ch.end.Duration == 0 {
		return ch.end, fmt.Errorf("end time is needed to calculate duration")
	}
	t := ch.end.Duration - ch.start.Duration
	stamp := Timestamp(t, float64(ch.base))
	return stamp, nil
}
