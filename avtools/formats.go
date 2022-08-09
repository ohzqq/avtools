package avtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	"gopkg.in/ini.v1"
)

type FileFormat struct {
	name   string
	file   string
	meta   *MediaMeta
	from   string
	to     string
	ext    string
	ini    *ini.File
	tmpl   *template.Template
	parse  func(file string) *MediaMeta
	render func(f *FileFormat) []byte
	data   []byte
}

func (f *FileFormat) Parse() *FileFormat {
	f.meta = f.parse(f.file)
	return f
}

func (f *FileFormat) Render() *FileFormat {
	d := f.render(f)
	f.data = d
	return f
}

//func NewFileFormat(file string) *FileFormat {
//  switch ext := path.Ext(file); ext {
//  case ".cue", "cue":
//    fmt, err := GetFormat("cue")
//    if err != nil {
//      log.Fatal(err)
//    }
//    fmt.meta = LoadCueSheet(file)
//    return fmt
//  case ".ini", "ini":
//    fmt, err := GetFormat("ffmeta")
//    if err != nil {
//      log.Fatal(err)
//    }
//    fmt.meta = LoadFFmetadataIni(file)
//    return fmt
//  }
//  return &FileFormat{}
//}

func (f *FileFormat) ConvertTo(kind string) *FileFormat {
	//  fmt.to = kind
	switch kind {
	case "json":
		f.data = MarshalJson(f)
	case "ffmeta":
		fmt.Printf("%+V\n", f.tmpl)
		f.data = RenderTmpl(f)
	case "cue":
		if len(f.meta.Chapters) == 0 {
			log.Fatal("No chapters")
		}
		f.data = RenderTmpl(f)
		//    return fmt.Render("cue")
	}
	return f
}

func JsonMeta(file string) *MediaMeta {
	cmd := NewFFprobeCmd(file)
	cmd.entries = ffProbeMeta
	cmd.showChaps = true
	cmd.format = "json"
	data := cmd.Parse().Run()

	meta := MediaMeta{}
	err := json.Unmarshal(data, &meta)
	if err != nil {
		fmt.Println("help")
	}
	meta.Tags = meta.Format.Tags

	return &meta
}

func MarshalJson(f *FileFormat) []byte {
	data, err := json.Marshal(f.meta)
	if err != nil {
		fmt.Println("help")
	}
	return data
}

func RenderTmpl(f *FileFormat) []byte {
	var buf bytes.Buffer
	err := f.tmpl.Execute(&buf, f.meta)
	if err != nil {
		log.Println("executing template:", err)
	}
	return buf.Bytes()
}

//func (f *FileFormat) Render(name string) *FileFormat {
//  tmpl, err := GetTmpl(name)
//  if err != nil {
//    log.Println("executing template:", err)
//  }
//  var buf bytes.Buffer
//  err = tmpl.Execute(&buf, f.meta)
//  if err != nil {
//    log.Println("executing template:", err)
//  }
//  f.data = buf.Bytes()
//  return f
//}

func GetTmpl(name string) (*template.Template, error) {
	var metaTmpl = map[string]*template.Template{
		"cue":          template.Must(template.New("cue").Funcs(funcs).Parse(cueTmpl)),
		"ffchaps":      template.Must(template.New("ffchaps").Funcs(funcs).Parse(ffmetaTmpl)),
		"cueToFFchaps": template.Must(template.New("cueToFFchaps").Funcs(funcs).Parse(ffmetaTmpl)),
		"ffmeta":       template.Must(template.New("ffmeta").Funcs(funcs).Parse(ffmetaTmpl)),
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

const ffmetaTmpl = `;FFMETADATA1
{{with .Format.Tags.Title -}}
	title={{.}}
{{end -}}
{{with .Format.Tags.Album -}}
	album={{.}}
{{end -}}
{{with .Format.Tags.Artist -}}
	artist={{.}}
{{end -}}
{{with .Format.Tags.Composer -}}
	composer={{.}}
{{end -}}
{{with .Format.Tags.Genre -}}
	genre={{.}}
{{end -}}
{{with .Format.Tags.Comment -}}
	comment={{.}}
{{end -}}
{{range $index, $ch := .Chapters -}}
[CHAPTER]
TIMEBASE=1/1000
START={{$ch.Start}}
END={{$ch.End}}
title={{if ne $ch.Title ""}}{{$ch.Title}}{{else}}Chapter {{$index}}{{end}}
{{end}}`
