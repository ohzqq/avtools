package ffmeta

import (
	"bytes"
	"html/template"
	"log"
	"strconv"
	"strings"

	"github.com/ohzqq/avtools/chap"
)

type Chapter struct {
	Base         string            `json:"time_base",ini:"timebase"`
	StartTime    int               `json:"start",ini:"start"`
	EndTime      int               `json:"end",ini:"end"`
	ChapterTitle string            `ini:"title"`
	Tags         map[string]string `json:"tags"`
}

func (c Chapter) Title() string {
	return c.ChapterTitle
}

func (c Chapter) Start() int {
	return c.StartTime
}

func (c Chapter) End() int {
	return c.EndTime
}

func (c Chapter) Timebase() float64 {
	if tb := c.Base; tb != "" {
		c.Base = strings.TrimPrefix(tb, "1/")
	}
	baseFloat, _ := strconv.ParseFloat(c.Base, 64)
	return baseFloat
}

func (ff *FFmeta) SetChapters(c chap.Chapters) *FFmeta {
	ff.Chapters = c
	return ff
}

func (ff FFmeta) LastChapterEnd() *chap.Chapter {
	ch := ff.LastChapter()
	if ch.End().Secs() == 0 && ff.Duration().Int() != 0 {
		ch.SetEnd(chap.NewChapterTime(ff.Duration().Float() * 1000))
	}
	return ch
}

func (ff FFmeta) IniChaps() []byte {
	var (
		tmpl = template.Must(template.New("ffmeta").Parse(ffmetaTmpl))
		buf  bytes.Buffer
	)

	ff.LastChapterEnd()

	err := tmpl.Execute(&buf, ff.Chapters)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

const ffmetaTmpl = `
{{- range .Each}}
[CHAPTER]
TIMEBASE={{.Timebase.String}}
START={{.Start.String}}
END={{.End.String}}
title={{.Title}}
{{- end -}}
`