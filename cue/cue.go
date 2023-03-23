package cue

import (
	"bufio"
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/dur"
)

type CueSheet struct {
	File   string
	Ext    string
	Tracks []avtools.ChapterMeta
}

type Track struct {
	StartStamp time.Duration
	EndStamp   time.Duration
	ChapTitle  string
}

func Load(file string) *CueSheet {
	var sheet CueSheet

	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	var times []time.Duration
	var titles []string
	scanner := bufio.NewScanner(contents)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.Contains(line, "TITLE"):
			title := strings.TrimPrefix(line, "TITLE ")
			title = strings.Trim(title, "'")
			title = strings.Trim(title, `"`)
			titles = append(titles, title)
		case strings.Contains(line, "INDEX 01"):
			stamp := strings.TrimPrefix(line, "INDEX 01 ")
			start, err := dur.Parse(stamp)
			if err != nil {
				log.Fatal(err)
			}
			times = append(times, start.Dur)
		}
	}

	e := 1
	for i := 0; i < len(titles); i++ {
		var track Track
		track.ChapTitle = titles[i]
		track.StartStamp = times[i]
		if e < len(titles) {
			track.EndStamp = times[e]
		}
		e++
		sheet.Tracks = append(sheet.Tracks, track)
	}

	return &sheet
}

func NewCueSheet(f string) *CueSheet {
	cue := &CueSheet{
		File: f,
		Ext:  filepath.Ext(f),
	}
	cue.Ext = strings.ToUpper(cue.Ext)
	cue.Ext = strings.TrimPrefix(cue.Ext, ".")
	return cue
}

func Dump(file string, meta avtools.Meta) []byte {
	var (
		tmpl = template.Must(template.New("cue").Funcs(tmplFuncs).Parse(cueTmpl))
		buf  bytes.Buffer
	)

	cue := NewCueSheet(file)
	//cue.Tracks = meta.Chapters()

	err := tmpl.Execute(&buf, cue)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (cue CueSheet) Chapters() []avtools.ChapterMeta {
	return cue.Tracks
}

func (cue CueSheet) Tags() map[string]string {
	return map[string]string{
		"filename": cue.File,
	}
}

func (cue CueSheet) Streams() []map[string]string {
	return []map[string]string{}
}

var tmplFuncs = template.FuncMap{
	"inc": Inc,
}

func Inc(n int) int {
	return n + 1
}

const cueTmpl = `FILE "{{.File}}" {{.Ext -}}
{{range $idx, $ch := .Tracks}}
TRACK {{inc $idx}} AUDIO
{{- if eq $ch.Title ""}}
  TITLE "Chapter {{inc $idx}}"
{{- else}}
  TITLE "{{$ch.Title}}"
{{- end}}
	INDEX 01 {{$ch.Start.MMSS}}:00
{{- end -}}`
