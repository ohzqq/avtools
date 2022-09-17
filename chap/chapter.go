package chap

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
)

type Number interface {
	int | int32 | int64 | float32 | float64
}

type Chapters struct {
	Chapters []*Chapter
	File     string
	ext      string
}

func NewChapters() Chapters {
	return Chapters{}
}

func (c Chapters) ToCue() []byte {
	var (
		tmpl = template.Must(template.New("cue").Parse(cueTmpl))
		buf  bytes.Buffer
	)
	err := tmpl.Execute(&buf, c)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (c *Chapters) SetExt(ext string) *Chapters {
	c.ext = ext
	return c
}

func (c Chapters) Ext() string {
	var ext string
	if c.ext != "" {
		ext = c.ext
	}
	ext = strings.TrimPrefix(filepath.Ext(c.File), ".")
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

type Chapter struct {
	start    Time
	end      Time
	Timebase float64
	Title    string
}

func NewChapter() *Chapter {
	return &Chapter{Timebase: 1}
}

func (ch Chapter) Start() Time {
	if t := ch.Timebase; t != 1 {
		ch.start.base = t
	}
	return ch.start
}

func (ch Chapter) End() Time {
	if t := ch.Timebase; t != 1 {
		ch.start.base = t
	}
	return ch.end
}

func (ch Chapter) Dur() (Time, error) {
	if ch.end.time == 0 {
		return ch.end, fmt.Errorf("end time is needed to calculate duration")
	}
	t := ch.end.time - ch.start.time
	return Time{time: t, base: ch.Timebase}, nil
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
	ch.Timebase = t
	return ch
}

const cueTmpl = `FILE "{{.File}}" {{.Ext}}
{{- range $index, $ch := .Each}}
TRACK {{$index}} AUDIO
  TITLE "{{$ch.Title}}"
  INDEX 01 {{$ch.Start.MMSS}}{{end -}}`
