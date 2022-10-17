package chap

import (
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

func (c Chapters) ToCue() *cue.CueSheet {
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
