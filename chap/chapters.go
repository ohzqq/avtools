package chap

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/ohzqq/avtools/tool/cue"
)

type Chapters struct {
	Chapters []*Chapter
	File     string
	ext      string
}

func NewChapters() Chapters {
	return Chapters{}
}

func (c Chapters) ToCue() []byte {
	//var (
	//  tmpl = template.Must(template.New("cue").Parse(cueTmpl))
	//  buf  bytes.Buffer
	//)

	//err := tmpl.Execute(&buf, c)
	//if err != nil {
	//  log.Fatal(err)
	//}

	sheet := cue.NewCueSheet(c.File)
	for _, ch := range c.Each() {
		track := cue.NewTrack()
		track.SetTitle(ch.Title)
		track.SetStart(ch.Start().Secs())
		sheet.Tracks = append(sheet.Tracks, track)
	}

	//return buf.Bytes()
	return sheet.Dump()
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
	println(string(c.ToCue()))
}

func (c Chapters) Write() {
	file, err := os.Create(slug.Make(c.File) + ".cue")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(c.ToCue())
	if err != nil {
		log.Fatal(err)
	}
}

func (c Chapters) Each() []*Chapter {
	return c.Chapters
}
