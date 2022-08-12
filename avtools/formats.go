package avtools

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type FileFormat struct {
	name     string
	meta     *MediaMeta
	from     string
	to       string
	File     string
	Path     string
	Dir      string
	Ext      string
	Mimetype string
	tmpl     *template.Template
	parse    func(file string) *MediaMeta
	render   func(meta *MediaMeta) []byte
	data     []byte
}

func NewFormat(input string) *FileFormat {
	f := FileFormat{}
	switch input {
	case "ffmeta":
	case "cue":
	default:
		abs, err := filepath.Abs(input)
		if err != nil {
			log.Fatal(err)
		}

		f.Path = abs
		f.File = filepath.Base(input)
		f.Dir = filepath.Dir(input)
		f.Ext = filepath.Ext(input)
		f.Mimetype = mime.TypeByExtension(f.Ext)
	}

	switch f.Ext {
	case "ini", ".ini":
		f.parse = LoadFFmetadataIni
		f.render = RenderFFmetaTmpl
		f.Parse()
	case "cue", ".cue":
		f.parse = LoadCueSheet
		f.render = RenderCueTmpl
		f.Parse()
	default:
		if f.IsAudio() {
			f.parse = EmbeddedJsonMeta
			f.render = MarshalJson
			f.Parse()
		}
	}
	//fmt.Printf("%+V\n", f.Path)
	return &f

}

func (f *FileFormat) Parse() *FileFormat {
	f.SetMeta(f.parse(f.Path))
	return f
}

func (f *FileFormat) Render() *FileFormat {
	f.data = f.render(f.meta)
	return f
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
	f.Ext = ext
	return f
}

func (f *FileFormat) Print() {
	println(f.String())
}

func (f *FileFormat) String() string {
	return string(f.data)
}

func (f FileFormat) IsImage() bool {
	if strings.Contains(f.Mimetype, "image") {
		return true
	}
	return false
}

func (f FileFormat) IsAudio() bool {
	if strings.Contains(f.Mimetype, "audio") {
		return true
	} else {
		fmt.Println("not an audio file")
	}
	return false
}

func (f FileFormat) IsPlainText() bool {
	if strings.Contains(f.Mimetype, "text/plain") {
		return true
	} else {
		log.Fatalln("needs to be plain text file")
	}
	return false
}

func (f FileFormat) IsFFmeta() bool {
	if f.IsPlainText() {
		contents, err := os.Open(f.Path)
		if err != nil {
			log.Fatal(err)
		}
		defer contents.Close()

		scanner := bufio.NewScanner(contents)
		line := 0
		for scanner.Scan() {
			if line == 0 && scanner.Text() == ";FFMETADATA1" {
				return true
			} else {
				log.Fatalln("ffmpeg metadata files need to have ';FFMETADATA1' as the first line")
			}
		}
	}
	return false
}

func EmbeddedJsonMeta(file string) *MediaMeta {
	data := NewFFprobeCmd(file).EmbeddedMeta()

	meta := MediaMeta{}
	err := json.Unmarshal(data, &meta)
	if err != nil {
		fmt.Println("help")
	}
	meta.Tags = meta.Format.Tags

	return &meta
}

func MarshalJson(meta *MediaMeta) []byte {
	data, err := json.Marshal(meta)
	if err != nil {
		fmt.Println("help")
	}
	return data
}

func RenderCueTmpl(meta *MediaMeta) []byte {
	var (
		buf  bytes.Buffer
		tmpl = template.Must(template.New("cue").Funcs(funcs).Parse(cueTmpl))
	)
	err := tmpl.Execute(&buf, meta)
	if err != nil {
		log.Println("executing template:", err)
	}
	return buf.Bytes()
}

func RenderFFmetaTmpl(meta *MediaMeta) []byte {
	meta.LastChapterEnd()
	var (
		buf  bytes.Buffer
		tmpl = template.Must(template.New("ffmeta").Funcs(funcs).Parse(ffmetaTmpl))
	)
	err := tmpl.Execute(&buf, meta)
	if err != nil {
		log.Println("executing template:", err)
	}
	return buf.Bytes()
}

//func GetTmpl(name string) (*template.Template, error) {
//  var metaTmpl = map[string]*template.Template{
//    "cue":          template.Must(template.New("cue").Funcs(funcs).Parse(cueTmpl)),
//    "ffchaps":      template.Must(template.New("ffchaps").Funcs(funcs).Parse(ffmetaTmpl)),
//    "cueToFFchaps": template.Must(template.New("cueToFFchaps").Funcs(funcs).Parse(ffmetaTmpl)),
//    "ffmeta":       template.Must(template.New("ffmeta").Funcs(funcs).Parse(ffmetaTmpl)),
//  }

//  for n, _ := range metaTmpl {
//    if n == name {
//      return metaTmpl[n], nil
//    }
//  }
//  return &template.Template{}, fmt.Errorf("%v is not a template", name)
//}

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
title={{with $ch.Title}}{{$ch.Title}}{{else}}Chapter {{$index}}{{end}}
{{end}}`
