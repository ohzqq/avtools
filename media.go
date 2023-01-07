package avtools

import (
	"time"
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
	StartTime Time
	EndTime   Time
	ChTitle   string
}

type ChapterMeta interface {
	Start() time.Duration
	End() time.Duration
	Title() string
}

type Meta interface {
	Chapters() []*Chapter
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
	m.Chapters = meta.Chapters()

	//if streams := meta.Streams(); len(streams) > 0 {
	m.Streams = meta.Streams()
	//}
	return m
}

func NewChapter(chap ChapterMeta) *Chapter {
	return &Chapter{
		ChTitle:   chap.Title(),
		StartTime: Timestamp(chap.Start()),
		EndTime:   Timestamp(chap.End()),
	}
}

func (ch Chapter) Start() Time {
	return ch.StartTime
}

func (ch Chapter) End() Time {
	return ch.EndTime
}

func (ch Chapter) Timebase() string {
	return "1/1000"
}

func (ch Chapter) Title() string {
	return ch.ChTitle
}

//func (ch Chapter) Dur() (Time, error) {
//  if ch.end.Duration == 0 {
//    return ch.end, fmt.Errorf("end time is needed to calculate duration")
//  }
//  t := ch.end.Duration - ch.start.Duration
//  stamp := NewTime(t, float64(ch.base))
//  return stamp, nil
//}
