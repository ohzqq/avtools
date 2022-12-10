package chap

import (
	"encoding/json"
	"log"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/ohzqq/avtools/cue"
)

type Chapters struct {
	Chapters []*Chapter
	File     string
	ext      string
}

func NewChapters() Chapters {
	return Chapters{}
}

func (c Chapters) FromCue(name string) Chapters {
	sheet := cue.Load(name)
	for _, t := range sheet.Tracks {
		ch := NewChapter()
		ch.SetTimebase(1000)
		ch.SetStart(NewChapterTime(t.Start() * 1000))
		ch.SetEnd(NewChapterTime(t.End() * 1000))
		ch.SetTitle(t.Title())
		c.Chapters = append(c.Chapters, ch)
	}
	return c
}

func (c Chapters) ToCue() *cue.Sheet {
	name := c.File
	if c.File == "" {
		name = "tmp"
	}
	sheet := cue.NewCueSheet(name)
	for _, ch := range c.Each() {
		track := cue.NewTrack()
		track.SetTitle(ch.Title)
		track.SetStart(ch.Start().Secs())
		sheet.Tracks = append(sheet.Tracks, track)
	}

	return sheet
}

func (c Chapters) ToJson() []byte {
	var chaps []map[string]string
	for _, ch := range c.Chapters {
		chap := map[string]string{
			"start": ch.Start().HHMMSS(),
			"end":   ch.End().HHMMSS(),
			"title": ch.Title,
		}
		chaps = append(chaps, chap)
	}

	data, err := json.Marshal(chaps)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (c Chapters) LastChapter() *Chapter {
	t := len(c.Chapters)
	if t > 0 {
		return c.Chapters[t-1]
	}
	return NewChapter()
}

func (c *Chapters) SetExt(ext string) *Chapters {
	c.ext = ext
	return c
}

func (c Chapters) Ext() string {
	var ext string
	ext = strings.TrimPrefix(filepath.Ext(c.File), ".")
	if c.ext != "" {
		ext = c.ext
	}
	return strings.ToUpper(ext)
}

func (c Chapters) Print() {
	println(string(c.ToCue().Dump()))
}

func (c Chapters) Write() {
	sheet := c.ToCue()
	name := slug.Make(c.File)
	err := sheet.SaveAs(name)
	if err != nil {
		log.Fatal(err)
	}
}

func (c Chapters) Each() []*Chapter {
	return c.Chapters
}
