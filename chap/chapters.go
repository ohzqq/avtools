package chap

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
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
