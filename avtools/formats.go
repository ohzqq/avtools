package avtools

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/gosimple/slug"
)

type Format struct {
	name   string
	meta   *MediaMeta
	from   string
	to     string
	ext    string
	tmpl   metaTemplates
	buffer *bytes.Buffer
}

func GetFormat(name string) (*Format, error) {
	var Formats = map[string]*Format{
		"cue": &Format{
			name: "cue",
			ext:  ".cue",
			meta: &MediaMeta{},
			tmpl: metaTemplate,
		},
		"ffmeta": &Format{
			name: "ffmeta",
			ext:  ".ini",
			meta: &MediaMeta{},
			tmpl: metaTemplate,
		},
		"json": &Format{
			name: "json",
			ext:  ".json",
			meta: &MediaMeta{},
			tmpl: metaTemplate,
		},
	}

	for n, _ := range Formats {
		if n == name {
			return Formats[n], nil
		}
	}
	return &Format{}, fmt.Errorf("%v is not a recognized format")
}

func (f *Format) ConvertTo(f string) *Format {
	f.to = f
	return f
}

func (f *Format) Print() {
	println(f.buffer.String())
}

func (f *Format) Write() {
	file, err := os.Create(slug.Make(c.title) + ".cue")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(c.cue.render().Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func (f *Format) Render() []byte {
	err := metaTmpl.cueToFFchaps.ExecuteTemplate(f.buffer, f.to, f.meta)
	if err != nil {
		log.Println("executing template:", err)
	}
}

func CueToFFmeta(c Chapters) string {
	var chaps bytes.Buffer
	err := metaTmpl.cueToFFchaps.ExecuteTemplate(&chaps, "cueToFF", c)
	if err != nil {
		log.Println("executing template:", err)
	}
	return chaps.String()
}

func (m *Media) FFmetaChapsToCue() {
	if !m.HasChapters() {
		log.Fatal("No chapters")
	}

	f, err := os.Create("chapters.cue")
	if err != nil {
		log.Fatal(err)
	}

	err = metaTmpl.cue.ExecuteTemplate(f, "cue", m)
	if err != nil {
		log.Println("executing template:", err)
	}
}

type metaTemplates struct {
	cue          *template.Template
	ffchaps      *template.Template
	cueToFFchaps *template.Template
}

var funcs = template.FuncMap{
	"cueStamp": secsToCueStamp,
}

var metaTmpl = metaTemplates{
	cue:          template.Must(template.New("cue").Funcs(funcs).Parse(cueTmpl)),
	ffchaps:      template.Must(template.New("ffchaps").Funcs(funcs).Parse(ffChapTmpl)),
	cueToFFchaps: template.Must(template.New("cueToFF").Funcs(funcs).Parse(cueToChapTmpl)),
}

const cueTmpl = `FILE '{{.File}}' {{.Ext}}
{{- range $index, $ch := .Meta.Chapters}}
TRACK {{$index}} AUDIO
  TITLE "Chapter {{$index}}"
  INDEX 01 {{cueStamp $ch.StartToSeconds}}{{end}}`

const cueToChapTmpl = `
{{- range $index, $ch := . -}}
[CHAPTER]
TITLE={{if ne $ch.Title ""}}{{$ch.Title}}{{else}}Chapter {{$index}}{{end}}
START={{$ch.Start}}
END={{$ch.End}}
TIMEBASE=1/1000
{{end}}`

const ffChapTmpl = `
{{- $media := . -}}
{{- range $index, $ch := $media.Meta.Chapters -}}
[CHAPTER]
TITLE={{if ne $ch.Title ""}}{{$ch.Title}}{{else}}Chapter {{$index}}{{end}}
START=
{{- if $media.CueChaps -}}
	{{- $ch.StartToIntString -}}
{{- else -}}
	{{- $ch.Start -}}
{{- end}}
END=
{{- if $media.CueChaps -}}
	{{- $ch.EndToIntString -}}
{{- else -}}
	{{- $ch.End -}}
{{- end}}
TIMEBASE=1/1000
{{end -}}`
