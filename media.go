package avtools

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ohzqq/dur"
	"github.com/ohzqq/fidi"
)

type Media struct {
	fidi.File
	Filename string
	tags     map[string]string
	Tagz     map[string]string
	chapters []*Chapter
	Chaps    []*Chapter
	streams  []map[string]string
	Streamz  []Stream
	Input    fidi.File
	Output   fidi.File
	Ini      fidi.File
	Cue      fidi.File
	Cover    fidi.File
	HasCover bool
}

type Stream struct {
	CodecType string
	CodecName string
	Index     string
	IsCover   bool
}

type Chapter struct {
	StartTime  Time
	EndTime    Time
	StartStamp dur.Timestamp
	EndStamp   dur.Timestamp
	ChapTitle  string
	Tags       map[string]string
}

type ChapterMeta interface {
	Start() time.Duration
	End() time.Duration
	Title() string
}

type Metaz interface {
	Chapters() []ChapterMeta
	Tags() map[string]string
	Streams() []map[string]string
	Source() fidi.File
}

type Meta interface {
	Chapters() []*Chapter
	Tags() map[string]string
	Streams() []map[string]string
}

func NewMedia() *Media {
	media := Media{
		tags: make(map[string]string),
	}
	return &media
}

func NewChapter(chap ChapterMeta) *Chapter {
	ss, err := dur.New(chap.Start())
	if err != nil {
		log.Fatal(err)
	}
	to, err := dur.New(chap.End())
	if err != nil {
		log.Fatal(err)
	}
	return &Chapter{
		ChapTitle:  chap.Title(),
		StartStamp: ss,
		EndStamp:   to,
		Tags:       make(map[string]string),
	}
}

func NewChapters(chaps []ChapterMeta) []*Chapter {
	var ch []*Chapter
	for _, c := range chaps {
		ch = append(ch, NewChapter(c))
	}
	return ch
}

func (m *Media) SetMeta(meta Meta) *Media {
	for key, val := range meta.Tags() {
		m.tags[key] = val
	}
	m.chapters = meta.Chapters()
	m.streams = meta.Streams()
	return m
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

func (m Media) Streams() []map[string]string {

	return m.streams
}

func (m Media) GetTag(key string) string {
	if val, ok := m.Tagz[key]; ok {
		return val
	}

	return ""
}

func (ch Chapter) Timebase() string {
	return "1/1000"
}

func (ch *Chapter) SS(ss string) *Chapter {
	dur := ParseStamp(ss)
	ch.StartTime = Timestamp(dur)
	return ch
}

func (ch *Chapter) To(to string) *Chapter {
	dur := ParseStamp(to)
	ch.EndTime = Timestamp(dur)
	return ch
}

func (ch Chapter) Start() time.Duration {
	return ch.StartStamp.Dur
}

func (ch Chapter) End() time.Duration {
	return ch.StartStamp.Dur
}

func (ch Chapter) Title() string {
	return ch.ChapTitle
}

func IsPlainText(mtype string) error {
	if strings.Contains(mtype, "text/plain") {
		return nil
	}
	return fmt.Errorf("needs to be plain text file")
}

//func (ch Chapter) Dur() (Time, error) {
//  if ch.end.Duration == 0 {
//    return ch.end, fmt.Errorf("end time is needed to calculate duration")
//  }
//  t := ch.end.Duration - ch.start.Duration
//  stamp := NewTime(t, float64(ch.base))
//  return stamp, nil
//}
