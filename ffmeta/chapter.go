package ffmeta

import (
	"bytes"
	"html/template"
	"log"
	"strconv"
	"strings"

	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/timestamp"
)

type Chapter struct {
	Base         string            `json:"time_base",ini:"timebase"`
	StartTime    int               `json:"start",ini:"start"`
	EndTime      int               `json:"end",ini:"end"`
	ChapterTitle string            `json:"title", ini:"title"`
	Tags         map[string]string `json:"tags"`
}

func (c Chapter) Title() string {
	if t, ok := c.Tags["title"]; ok {
		return t
	}
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

func (ff *Meta) SetChapters(c chap.Chapters) *Meta {
	ff.Chapters = c
	return ff
}

func (ff Meta) LastChapterEnd() *chap.Chapter {
	ch := ff.Chapters.LastChapter()
	if ch.End().Secs() == 0 && ff.Duration().Secs() != 0 {
		to := ff.Duration().Float() * 1000
		end := timestamp.NewChapterTime(to)
		ch.SetEnd(end)
	}
	return ch
}

var tmplFuncs = template.FuncMap{
	"inc": Inc,
}

func Inc(n int) int {
	return n + 1
}

func (ff Meta) IniChaps() []byte {
	var (
		tmpl = template.Must(template.New("ffmeta").Funcs(tmplFuncs).Parse(ffmetaTmpl))
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
{{- range $idx, $ch := .Each}}
[CHAPTER]
TIMEBASE={{$ch.Timebase.String}}
START={{$ch.Start.String}}
END={{$ch.End.String}}
{{- if eq $ch.Title ""}}
title=Chapter {{inc $idx}}
{{- else}}
title={{$ch.Title}}
{{- end}}
{{- end -}}
`
