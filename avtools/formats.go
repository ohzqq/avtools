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
	name   string
	file   string
	ext    string
	meta   *MediaMeta
	from   string
	to     string
	tmpl   *template.Template
	parse  func(file string) *MediaMeta
	render func(meta *MediaMeta) []byte
	data   []byte
}

func NewFormat(kind string) *FileFormat {
	switch kind {
	case "ffmeta":
		return NewFFmeta()
	case "cue":
		return NewCueSheet()
	case "audio", "input":
		return NewMediaFile()
	default:
		return &FileFormat{}
	}
}

func NewMediaFile() *FileFormat {
	return &FileFormat{
		parse:  EmbeddedJsonMeta,
		render: MarshalJson,
	}
}

func NewCueSheet() *FileFormat {
	return &FileFormat{
		ext:    ".cue",
		parse:  LoadCueSheet,
		render: RenderCueTmpl,
	}
}

func NewFFmeta() *FileFormat {
	return &FileFormat{
		ext:    ".ini",
		parse:  LoadFFmetadataIni,
		render: RenderFFmetaTmpl,
	}
}

func (f *FileFormat) HasFile() bool {
	return f.file != ""
}

func (f *FileFormat) Ext() string {
	if f.HasFile() {
		return filepath.Ext(f.file)
	}
	return f.ext
}

func (f *FileFormat) Meta() *MediaMeta {
	return f.meta
}

func (f *FileFormat) Name() string {
	if f.HasFile() {
		return strings.TrimSuffix(filepath.Base(f.file), filepath.Ext(f.file))
	}
	if f.meta != nil && f.meta.GetTag("title") != "" {
		return f.meta.GetTag("title")
	}
	if f.name != "" {
		return f.name
	}
	return "tmp"
}

func (f *FileFormat) Path() string {
	abs, err := filepath.Abs(f.file)
	if err != nil {
		log.Fatal(err)
	}
	return abs
}

func (f *FileFormat) Mimetype() string {
	return mime.TypeByExtension(f.Ext())
}

func (f *FileFormat) Dir() string {
	return filepath.Dir(f.file)
}

func (f *FileFormat) Parse() *FileFormat {
	f.SetMeta(f.parse(f.Path()))
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

func (f *FileFormat) SetName(n string) *FileFormat {
	f.name = n
	return f
}

func (f *FileFormat) SetFile(input string) *FileFormat {
	f.file = input
	return f
}

func (f *FileFormat) Print() {
	println(f.String())
}

func (f *FileFormat) Write() {
	WriteFile(f.Name(), f.Ext(), f.Bytes())
}

func (f *FileFormat) String() string {
	return string(f.Bytes())
}

func (f *FileFormat) Bytes() []byte {
	return f.render(f.meta)
}

func (f FileFormat) IsImage() bool {
	if strings.Contains(f.Mimetype(), "image") {
		return true
	}
	return false
}

func (f FileFormat) IsAudio() bool {
	if strings.Contains(f.Mimetype(), "audio") {
		return true
	} else {
		fmt.Println("not an audio file")
	}
	return false
}

func (f FileFormat) IsPlainText() bool {
	if strings.Contains(f.Mimetype(), "text/plain") {
		return true
	} else {
		log.Fatalln("needs to be plain text file")
	}
	return false
}

func (f FileFormat) IsFFmeta() bool {
	if f.IsPlainText() {
		contents, err := os.Open(f.Path())
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
