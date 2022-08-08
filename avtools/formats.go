package avtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"text/template"
)

type FileFormat struct {
	name   string
	file   string
	meta   *MediaMeta
	from   string
	to     string
	ext    string
	render func(to string) *FileFormat
	data   []byte
}

func NewFileFormat(file string) *FileFormat {
	switch ext := path.Ext(file); ext {
	case ".cue":
		fmt, err := GetFormat("cue")
		if err != nil {
			log.Fatal(err)
		}
		fmt.meta = LoadCueSheet(file)
		return fmt
	case ".ini":
		fmt, err := GetFormat("ffmeta")
		if err != nil {
			log.Fatal(err)
		}
		fmt.meta = LoadFFmetadataIni(file)
		return fmt
	}
	return &FileFormat{}
}

func GetFormat(name string) (*FileFormat, error) {
	var formats = map[string]*FileFormat{
		"cue": &FileFormat{
			name: "cue",
			ext:  ".cue",
			meta: &MediaMeta{},
		},
		"ffmeta": &FileFormat{
			name: "ffmeta",
			ext:  ".ini",
			meta: &MediaMeta{},
		},
		"json": &FileFormat{
			name: "json",
			ext:  ".json",
			meta: &MediaMeta{},
		},
	}

	for n, _ := range formats {
		if n == name {
			return formats[n], nil
		}
	}
	return &FileFormat{}, fmt.Errorf("%v is not a recognized format", name)
}

func (fmt *FileFormat) ConvertTo(kind string) *FileFormat {
	fmt.to = kind
	switch fmt.to {
	case "json":
		return fmt.MarshalJson()
	case "ffmeta":
		return fmt.Render("cueToFFchaps")
	case "cue":
		if len(fmt.meta.Chapters) == 0 {
			log.Fatal("No chapters")
		}
		return fmt.Render("cue")
	}
	return fmt
}

func (f *FileFormat) MarshalJson() *FileFormat {
	data, err := json.Marshal(f.meta)
	if err != nil {
		fmt.Println("help")
	}
	f.data = data
	return f
}

func (f *FileFormat) Render(name string) *FileFormat {
	tmpl, err := GetTmpl(name)
	if err != nil {
		log.Println("executing template:", err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, f.meta)
	if err != nil {
		log.Println("executing template:", err)
	}
	f.data = buf.Bytes()
	return f
}

func GetTmpl(name string) (*template.Template, error) {
	var metaTmpl = map[string]*template.Template{
		"cue":          template.Must(template.New("cue").Funcs(funcs).Parse(cueTmpl)),
		"ffchaps":      template.Must(template.New("ffchaps").Funcs(funcs).Parse(ffChapTmpl)),
		"cueToFFchaps": template.Must(template.New("cueToFFchaps").Funcs(funcs).Parse(cueToChapTmpl)),
	}

	for n, _ := range metaTmpl {
		if n == name {
			return metaTmpl[n], nil
		}
	}
	return &template.Template{}, fmt.Errorf("%v is not a template", name)
}

func (f *FileFormat) SetMeta(m *MediaMeta) *FileFormat {
	f.meta = m
	return f
}

func (f *FileFormat) SetFileName(n string) *FileFormat {
	f.name = n
	return f
}

func (f *FileFormat) SetExt(ext string) *FileFormat {
	f.ext = ext
	return f
}

func (f *FileFormat) Print() {
	println(f.String())
}

func (f *FileFormat) String() string {
	return string(f.data)
}

type metaTemplates struct {
	cue          *template.Template
	ffchaps      *template.Template
	cueToFFchaps *template.Template
}

var funcs = template.FuncMap{
	"cueStamp": secsToCueStamp,
}

const cueTmpl = `FILE "{{.Format.Filename}}" {{.Format.Ext}}
{{- range $index, $ch := .Chapters}}
TRACK {{$index}} AUDIO
  TITLE {{if ne $ch.Title ""}}{{$ch.Title}}{{else}}Chapter {{$index}}{{end}}
  INDEX 01 {{$ch.CueStamp}}
{{- end}}`

const cueToChapTmpl = `;FFMETADATA1
{{range $index, $ch := .Chapters -}}
[CHAPTER]
TIMEBASE=1/1000
START={{$ch.Start}}
END={{$ch.End}}
title={{if ne $ch.Title ""}}{{$ch.Title}}{{else}}Chapter {{$index}}{{end}}
{{end}}`

const ffChapTmpl = `;FFMETADATA1
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
