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

type FileFormats struct {
	formats []*FileFormat
	from    *FileFormat
	to      *FileFormat
}

func (f FileFormats) ListFormats() []string {
	var formats []string
	for _, f := range f.formats {
		formats = append(formats, f.Ext)
	}
	return formats
}

func (f FileFormats) GetFormat(format string) *FileFormat {
	var ext string
	switch format {
	case "ffmeta", "ini":
		ext = ".ini"
	case "cue", "cuesheet":
		ext = ".cue"
	default:
		ext = format
	}
	for _, fmt := range f.formats {
		switch ext {
		case "audio":
			if fmt.IsAudio() {
				return fmt
			}
		case fmt.Ext:
			return fmt
		}
	}
	return nil
}

func (f *FileFormats) AddFormat(file string) *FileFormats {
	format := NewFormat(file)
	format.Parse()
	f.formats = append(f.formats, format)
	return f
}

func (f FileFormats) HasCue() bool {
	return f.GetFormat(".cue") != nil
}

func (f FileFormats) HasFFmeta() bool {
	return f.GetFormat(".ini") != nil
}

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
	render   func(f *FileFormat) []byte
	data     []byte
}

func NewFormat(input string) *FileFormat {
	mime.AddExtensionType(".ini", "text/plain")
	mime.AddExtensionType(".cue", "text/plain")
	mime.AddExtensionType(".m4b", "audio/mp4")

	abs, err := filepath.Abs(input)
	if err != nil {
		log.Fatal(err)
	}

	f := FileFormat{
		Path: abs,
		File: filepath.Base(input),
		Dir:  filepath.Dir(input),
		Ext:  filepath.Ext(input),
	}
	f.Mimetype = mime.TypeByExtension(f.Ext)

	switch f.Ext {
	case ".ini", "ini", "ffmeta":
		f.parse = LoadFFmetadataIni
		f.render = RenderTmpl
		f.tmpl = template.Must(template.New("ffmeta").Funcs(funcs).Parse(ffmetaTmpl))
	case ".cue", "cue":
		f.parse = LoadCueSheet
		f.render = RenderTmpl
		f.tmpl = template.Must(template.New("cue").Funcs(funcs).Parse(cueTmpl))
	case ".jpg", "jpg", ".png", "png":
	default:
		if f.IsAudio() {
			f.parse = EmbeddedJsonMeta
			f.render = MarshalJson
		}
	}
	//fmt.Printf("%+V\n", f.Path)
	return &f

}

func (f *FileFormat) Parse() *FileFormat {
	f.meta = f.parse(f.Path)
	return f
}

func (f *FileFormat) Render() *FileFormat {
	f.data = f.render(f)
	return f
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
	f.Ext = ext
	return f
}

func (f *FileFormat) Print() {
	println(f.String())
}

func (f *FileFormat) String() string {
	return string(f.data)
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
