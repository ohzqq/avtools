package avtools

import (
	"log"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"text/template"
	//"regexp"
	"encoding/json"
	//"reflect"

	"github.com/go-ini/ini"
)
var _ = fmt.Printf

type MediaMeta struct {
	Chapters []*Chapter
	Streams []*Stream
	Format *Format
}

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Filename string
	Duration string
	Size string
	BitRate string `json:"bit_rate"`
	Tags *Tags `json:"tags"`
}

type Tags struct {
	Title string `json:"title"`
	Artist string `json:"artist"`
	Composer string `json:"composer"`
	Album string `json:"album"`
	Comment string `json:"comment"`
	Genre string `json:"genre"`
}

type Chapter struct {
	Timebase string `json:"time_base"`
	Start string `json:"start_time"`
	End string `json:"end_time"`
	Tags *Tags `json:"tags"`
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

const ffProbeMeta = `format=filename,start_time,duration,size,bit_rate:stream=codec_type,codec_name:format_tags`

func ReadEmbeddedMeta(input string) *MediaMeta {
	ff := NewFFprobeCmd()
	ff.In(input)
	ff.Args().Entries(ffProbeMeta).Chapters().Verbosity("error").Format("json")

	m := ff.Run()

	media := MediaMeta{}

	err := json.Unmarshal(m, &media)
	if err != nil {
		fmt.Println("help")
	}

	return &media
}

func WriteFFmetadata(input string) {
	cmd := NewFFmpegCmd()
	cmd.In(NewMedia().Input(input))
	cmd.Args().Post("f", "ffmetadata").ACodec("none").VCodec("none").Ext("ini")
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

	media := &MediaMeta{}

	meta, err := f.GetSection("")
	if err != nil {
		log.Fatal(err)
	}
	meta.MapTo(media.Format.Tags)

	if f.HasSection("chapter") {
		sec, _ := f.SectionsByName("chapter")
		for _, chap := range sec {
			c := Chapter{}
			err := chap.MapTo(&c)
			ss, to := ffChapstoSeconds(c.Timebase, c.Start, c.End)
			c.Start = ss
			c.End = to
			if err != nil { log.Fatal(err) }
			media.Chapters = append(media.Chapters, &c)
		}
	}

	return media
}

func ReadCueSheet(file string) []*Chapter {
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

	var tracks []*Chapter
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

	return tracks
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

