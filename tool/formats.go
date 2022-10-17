package tool

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/cue"
	"gopkg.in/ini.v1"
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
	case "audio", "input", "json":
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
	if title := f.meta.GetTag("title"); f.meta != nil && title != "" {
		return title
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

func EmbeddedMeta(file string) []byte {
	return NewFFprobeCmd(file).EmbeddedMeta()
}

func MarshalJson(meta *MediaMeta) []byte {
	for _, ch := range meta.Chapters {
		if ch.Timebase == "" {
			ch.Timebase = "1/1000"
		}
	}

	data, err := json.Marshal(meta)
	if err != nil {
		fmt.Println("help")
	}
	return data
}

func LoadCue(file string) *MediaMeta {
	cue := cue.Load(file)
	var chapters chap.Chapters
	for _, track := range cue.Tracks {
		ch := chap.NewChapter().SetMeta(track)
		chapters.Chapters = append(chapters.Chapters, ch)
	}

	return &MediaMeta{
		Ch: chapters,
	}
}

func LoadCueSheet(file string) *MediaMeta {
	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	var (
		titles     []string
		startTimes []int
		meta       = MediaMeta{
			Format: &Format{},
			Ch:     chap.NewChapters(),
		}
		fileRegexp = regexp.MustCompile(`^(\w+ )('|")(?P<title>.*)("|')( .*)$`)
	)

	scanner := bufio.NewScanner(contents)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if strings.Contains(s, "FILE") {
			matches := fileRegexp.FindStringSubmatch(s)
			meta.Format.Filename = matches[fileRegexp.SubexpIndex("title")]
		}
		if strings.Contains(s, "TITLE") {
			t := strings.TrimPrefix(s, "TITLE ")
			t = strings.Trim(t, "'")
			t = strings.Trim(t, `"`)
			titles = append(titles, t)
		} else if strings.Contains(s, "INDEX") {
			start := cueStampToFFmpegTime(strings.TrimPrefix(s, "INDEX 01 "))
			startTimes = append(startTimes, start)
		}
	}

	e := 1
	for i := 0; i < len(titles); i++ {
		t := Chapter{}
		ch := chap.NewChapter().SetTimebase(1000)
		t.Title = titles[i]
		ch.Title = titles[i]
		ss := chap.NewChapterTime(startTimes[i])
		ch.SetStart(ss)
		t.Start = startTimes[i]
		if e < len(titles) {
			t.End = startTimes[e]
			to := chap.NewChapterTime(startTimes[e])
			ch.SetEnd(to)
		}
		e++
		meta.Ch.Chapters = append(meta.Ch.Chapters, ch)
		meta.Chapters = append(meta.Chapters, &t)
	}

	return &meta
}

func RenderCueTmpl(meta *MediaMeta) []byte {
	const cueTmpl = `FILE "{{.Format.Filename}}" {{.Format.Ext}}
{{- range $index, $ch := .Chapters}}
TRACK {{$index}} AUDIO
  TITLE {{if ne $ch.Title ""}}{{$ch.Title}}{{else}}Chapter {{$index}}{{end}}
  INDEX 01 {{$ch.CueStamp}}
{{- end}}`

	var (
		buf  bytes.Buffer
		tmpl = template.Must(template.New("cue").Parse(cueTmpl))
	)

	err := tmpl.Execute(&buf, meta)
	if err != nil {
		log.Println("executing template:", err)
	}

	return buf.Bytes()
}

func LoadFFmetadataIni(input string) *MediaMeta {
	opts := ini.LoadOptions{}
	opts.Insensitive = true
	opts.InsensitiveSections = true
	opts.IgnoreInlineComment = true
	opts.AllowNonUniqueSections = true

	abs, _ := filepath.Abs(input)
	f, err := ini.LoadSources(opts, abs)
	if err != nil {
		log.Fatal(err)
	}

	media := MediaMeta{
		Format: &Format{
			Tags: f.Section("").KeysHash(),
		},
	}

	if f.HasSection("chapter") {
		sec, _ := f.SectionsByName("chapter")
		for _, chap := range sec {
			c := Chapter{}
			err := chap.MapTo(&c)
			if err != nil {
				log.Fatal(err)
			}
			media.Chapters = append(media.Chapters, &c)
		}
	}
	return &media
}

func RenderFFmetaTmpl(meta *MediaMeta) []byte {
	const ffmetaTmpl = `;FFMETADATA1
{{range $key, $val := .Format.Tags -}}
	{{$key}}={{$val}}
{{end -}}
{{range $index, $ch := .Chapters -}}
[CHAPTER]
TIMEBASE={{$ch.TimeBase}}
START={{$ch.Start}}
END={{$ch.End}}
title={{with $ch.Title}}{{$ch.Title}}{{else}}Chapter {{$index}}{{end}}
{{end}}`

	var (
		buf  bytes.Buffer
		tmpl = template.Must(template.New("ffmeta").Parse(ffmetaTmpl))
	)

	meta.LastChapterEnd()

	err := tmpl.Execute(&buf, meta)
	if err != nil {
		log.Println("executing template:", err)
	}

	return buf.Bytes()
}
