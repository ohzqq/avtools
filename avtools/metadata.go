package avtools

import (
	"log"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	//"regexp"
	"encoding/json"
	//"reflect"

	"github.com/go-ini/ini"
)
var _ = fmt.Printf

type MediaMeta struct {
	Chapters *Chapters
	Streams *Streams
	Format *Format
	Tags *Tags
}

type Streams []*Stream

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Filename string
	Duration string
	Size string
	BitRate string `json:"bit_rate"`
}

type Tags struct {
	Title string `json:"title"`
	Artist string `json:"artist"`
	Composer string `json:"composer"`
	Album string `json:"album"`
	Comment string `json:"comment"`
	Genre string `json:"genre"`
}

type Chapters []*Chapter

type Chapter struct {
	Timebase string `json:"time_base"`
	Start string `json:"start_time"`
	End string `json:"end_time"`
	Title string
}

func(c Chapter) StartTimebaseProduct() string {
	ss, _ := strconv.ParseFloat(c.Start, 64)
	result := ss * c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 0, 64)
}

func(c Chapter) EndTimebaseProduct() string {
	to, _ := strconv.ParseFloat(c.End, 64)
	result := to * c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 0, 64)
}

func(c Chapter) TimebaseFloat() float64 {
	base := "1000"
	if tb := c.Timebase; tb != "" {
		base = strings.ReplaceAll(tb, "1/", "")
	}
	baseFloat, _ := strconv.ParseFloat(base, 64)
	return baseFloat
}

type jsonMeta struct {
	Chapters []jsonChapter
	Streams *Streams
	Format jsonFormat
	Tags *Tags
}

type jsonFormat struct {
	Format
	Tags *Tags
}

type jsonChapter struct {
	Chapter
	Tags *Tags `json:"tags"`
}

type metaTemplates struct {
	cue *template.Template
	ffchaps *template.Template
}

var funcs = template.FuncMap{
	"cueStamp": secsToCueStamp,
}

var metaTmpl = metaTemplates{
	cue: template.Must(template.New("cue").Funcs(funcs).Parse(cueTmpl)),
	ffchaps: template.Must(template.New("ffchaps").Funcs(funcs).Parse(ffChapTmpl)),
}

const cueTmpl = `FILE '{{.File}}' {{.Ext}}
{{- range $index, $ch := .Meta.Chapters}}
TRACK {{$index}} AUDIO
  TITLE "Chapter {{$index}}"
  INDEX 01 {{cueStamp $ch.Start}}{{end}}`

const ffChapTmpl = `
{{- range $index, $ch := .Meta.Chapters -}}
[CHAPTER]
TITLE={{if ne $ch.Title ""}}{{.}}{{else}}Chapter {{$index}}{{end}}
START={{$ch.StartTimebaseProduct}}
END={{$ch.EndTimebaseProduct}}
TIMEBASE=1/1000
{{end -}}`

func ReadEmbeddedMeta(input string) *MediaMeta {
	ff := NewFFProbeCmd()
	ff.In(input)
	ff.Args().
		Entries("format=filename,start_time,duration,size,bit_rate:stream=codec_type,codec_name:format_tags").
		Chapters().
		Verbosity("error").
		Format("json")

	m := ff.Run()

	var meta jsonMeta
	err := json.Unmarshal(m, &meta)
	if err != nil { fmt.Println("help")}

	media := MediaMeta{}
	media.Streams = meta.Streams
	media.Format = &Format{
		Filename: meta.Format.Filename,
		Duration: meta.Format.Duration,
		Size: meta.Format.Size,
		BitRate: meta.Format.BitRate,
	}
	media.Tags = meta.Format.Tags

	var chapters Chapters
	for _, ch := range meta.Chapters {
		c := new(Chapter)
		c.Title = ch.Tags.Title
		c.Timebase = ch.Timebase
		c.Start = ch.Start
		c.End = ch.End
		chapters = append(chapters, c)
	}
	media.Chapters = &chapters
	return &media
}

func WriteFFmetadata(input string) {
	cmd := NewCmd()
	cmd.In(NewMedia(input))
	cmd.Args().Post("f", "ffmetadata").ACodec("none").VCodec("none").Ext("ini")
	//fmt.Printf("%v", cmd.Cmd().String())
	cmd.Run()
}

func ReadFFmetadata(input string) *MediaMeta {
	opts := ini.LoadOptions{}
	opts.Insensitive = true
	opts.InsensitiveSections = true
	opts.IgnoreInlineComment = true
	opts.AllowNonUniqueSections = true

	f, err := ini.LoadSources(opts, input)
	if err != nil {
		log.Fatal(err)
	}

	meta, err := f.GetSection("")
	if err != nil {
		log.Fatal(err)
	}
	m := new(Tags)
	meta.MapTo(m)
	tags := *m

	var chapters Chapters
	if f.HasSection("chapter") {
		sec, _ := f.SectionsByName("chapter")
		for _, chap := range sec {
			c := Chapter{}
			err := chap.MapTo(&c)
			ss, to := ffChapstoSeconds(c.Timebase, c.Start, c.End)
			c.Start = ss
			c.End = to
			if err != nil { log.Fatal(err) }
			chapters = append(chapters, &c)
		}
	}

	return &MediaMeta{
		Chapters: &chapters,
		Tags: &tags,
	}
}

func ReadCueSheet(file string) *Chapters {
	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	var (
		titles []string
		indices []string
	)
	scanner := bufio.NewScanner(contents)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if strings.Contains(s, "TITLE") {
			t := strings.TrimPrefix(s, "TITLE ")
			t = strings.Trim(t, "'")
			t = strings.Trim(t, `"`)
			titles = append(titles, t)
		} else if strings.Contains(s, "INDEX") {
			start := cueStampToSecs(strings.TrimPrefix(s, "INDEX 01 "))
			indices = append(indices, start)
		}
	}

	var tracks Chapters
	e := 1
	for i := 0; i < len(titles); i++ {
		t := new(Chapter)
		t.Title = titles[i]
		t.Start = indices[i]
		if e < len(titles) {
			t.End = indices[e]
		}
		e++
		tracks = append(tracks, t)
		//fmt.Printf("%v", t)
	}

	return &tracks
}
