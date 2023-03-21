package meta

import (
	"bufio"
	"bytes"
	"html/template"
	"io"
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
	Tracks []*avtools.Chapter
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

func LoadCueSheet(file string) *CueSheet {
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
		track := &avtools.Chapter{}
		track.Title = titles[i]
		ss := avtools.Timestamp(times[i])
		track.Start = ss
		if e < len(titles) {
			track.End = avtools.Timestamp(times[e])
		}
		e++
		sheet.Tracks = append(sheet.Tracks, track)
	}

	return &sheet
}

func DumpCueSheet(file string, meta avtools.Meta) []byte {
	var (
		tmpl = template.Must(template.New("cue").Funcs(tmplFuncs).Parse(cueTmpl))
		buf  bytes.Buffer
	)

	cue := NewCueSheet(file)
	cue.Tracks = meta.Chapters()

	err := tmpl.Execute(&buf, cue)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (cue CueSheet) Chapters() []*avtools.Chapter {
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

func (cue CueSheet) Dump() []byte {
	var (
		tmpl = template.Must(template.New("cue").Funcs(tmplFuncs).Parse(cueTmpl))
		buf  bytes.Buffer
	)

	err := tmpl.Execute(&buf, cue)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (cue CueSheet) Write(wr io.Writer) error {
	_, err := wr.Write(cue.Dump())
	if err != nil {
		return err
	}
	return nil
}

func (cue CueSheet) Save() error {
	return cue.SaveAs(cue.File)
}

func (cue CueSheet) SaveAs(name string) error {
	if name == "" || cue.File == "" {
		name = "tmp"
	}

	file, err := os.Create(name + ".cue")
	if err != nil {
		return err
	}
	defer file.Close()

	err = cue.Write(file)
	if err != nil {
		return err
	}

	return nil
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
