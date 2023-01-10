package avtools

import (
	"time"
)

type Media struct {
	Filename string
	tags     map[string]string
	chapters []*Chapter
	streams  []map[string]string
}

type Chapter struct {
	Start Time
	End   Time
	Title string
	Tags  map[string]string
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
		tags:     make(map[string]string),
	}
	return &media
}

func NewChapter(chap ChapterMeta) *Chapter {
	return &Chapter{
		Title: chap.Title(),
		Start: Timestamp(chap.Start()),
		End:   Timestamp(chap.End()),
		Tags:  make(map[string]string),
	}
}

func (m Media) Chapters() []*Chapter {
	return m.chapters
}

func (m *Media) SetChapters(chaps []*Chapter) {
	m.chapters = chaps
}

func (m Media) Tags() map[string]string {
	return m.tags
}

func (m *Media) SetTags(tags map[string]string) {
	m.tags = tags
}

func (m Media) Streams() []map[string]string {
	return m.streams
}

func (m *Media) SetStreams(streams []map[string]string) {
	m.streams = streams
}

func (m Media) GetTag(key string) string {
	if val, ok := m.tags[key]; ok {
		return val
	}

	return ""
}

func (m *Media) SetMeta(meta Meta) *Media {
	for key, val := range meta.Tags() {
		m.tags[key] = val
	}

	m.chapters = meta.Chapters()

	m.streams = meta.Streams()
	return m
}

func (ch Chapter) Timebase() string {
	return "1/1000"
}

func (ch *Chapter) SS(ss string) *Chapter {
	dur := ParseStamp(ss)
	ch.Start = Timestamp(dur)
	return ch
}

func (ch *Chapter) To(to string) *Chapter {
	dur := ParseStamp(to)
	ch.End = Timestamp(dur)
	return ch
}

//func (ch Chapter) Dur() (Time, error) {
//  if ch.end.Duration == 0 {
//    return ch.end, fmt.Errorf("end time is needed to calculate duration")
//  }
//  t := ch.end.Duration - ch.start.Duration
//  stamp := NewTime(t, float64(ch.base))
//  return stamp, nil
//}
