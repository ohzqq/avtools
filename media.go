package avtools

import (
	"fmt"

	"github.com/ohzqq/avtools/timestamp"
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
	start timestamp.Time
	end   timestamp.Time
	base  timestamp.Timebase
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
			ch := NewChapter(
				chap.Start(),
				chap.End(),
				chap.Timebase(),
			)
			ch.title = chap.Title()
			m.Chapters = append(m.Chapters, ch)
		}
	}
	//if streams := meta.Streams(); len(streams) > 0 {
	m.Streams = meta.Streams()
	//}
	return m
}

func NewChapter[N Number](times ...N) *Chapter {
	var base float64 = 1

	var chapter Chapter

	switch t := len(times); t {
	case 3:
		base = float64(times[2])
		chapter.base = timestamp.Timebase(base)
		fallthrough
	case 2:
		chapter.end = timestamp.NewerTimeStamp(times[1], base)
		fallthrough
	case 1:
		chapter.start = timestamp.NewerTimeStamp(times[0], base)
	}

	return &chapter
}

func (ch Chapter) Start() timestamp.Time {
	return ch.start
}

func (ch Chapter) End() timestamp.Time {
	return ch.end
}

func (ch Chapter) Timebase() timestamp.Timebase {
	return ch.base
}

func (ch Chapter) Title() string {
	return ch.title
}

func (ch Chapter) Dur() (timestamp.Time, error) {
	if ch.end.Duration == 0 {
		return ch.end, fmt.Errorf("end time is needed to calculate duration")
	}
	t := ch.end.Duration - ch.start.Duration
	stamp := timestamp.NewerTimeStamp(t, float64(ch.base))
	return stamp, nil
}
