package avtools

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-ini/ini"
)

const ffProbeMeta = `format=filename,start_time,duration,size,bit_rate:stream=codec_type,codec_name:format_tags`

type MediaMeta struct {
	Chapters []*Chapter
	Streams  []*Stream
	Format   *Format
	Tags     map[string]string
}

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Filename string
	Duration string
	Size     string
	BitRate  string `json:"bit_rate"`
	Tags     map[string]string
}

type Tags struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Composer string `json:"composer"`
	Album    string `json:"album"`
	Comment  string `json:"comment"`
	Genre    string `json:"genre"`
}

type Chapter struct {
	Timebase string `json:"time_base"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Tags     map[string]string
	Title    string `ini:"title"`
}

func (c *Chapter) StartToIntString() string {
	result := float64(c.Start) * c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 0, 64)
}

func (c *Chapter) StartToSeconds() string {
	if c.Start == 0 {
		return ""
	}
	result := float64(c.Start) / c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 3, 64)
}

func (c *Chapter) EndToIntString() string {
	result := float64(c.End) * c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 0, 64)
}

func (c *Chapter) EndToSeconds() string {
	if c.End == 0 {
		return ""
	}
	result := float64(c.End) / c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 3, 64)
}

func (c Chapter) TimebaseFloat() float64 {
	base := "1000"
	if tb := c.Timebase; tb != "" {
		base = strings.ReplaceAll(tb, "1/", "")
	}
	baseFloat, _ := strconv.ParseFloat(base, 64)
	return baseFloat
}

func LoadFFmetadataIni(input string) *MediaMeta {
	opts := ini.LoadOptions{}
	opts.Insensitive = true
	opts.InsensitiveSections = true
	opts.IgnoreInlineComment = true
	opts.AllowNonUniqueSections = true

	f, err := ini.LoadSources(opts, input)
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

func LoadCueSheet(file string) []*Chapter {
	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	var (
		titles     []string
		startTimes []int
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
			start := cueStampToFFmpegTime(strings.TrimPrefix(s, "INDEX 01 "))
			startTimes = append(startTimes, start)
		}
	}

	var tracks []*Chapter
	e := 1
	for i := 0; i < len(titles); i++ {
		t := new(Chapter)
		t.Title = titles[i]
		t.Start = startTimes[i]
		if e < len(titles) {
			t.End = startTimes[e]
		}
		e++
		tracks = append(tracks, t)
	}

	return tracks
}

type metaTemplates struct {
	cue     *template.Template
	ffchaps *template.Template
}

var funcs = template.FuncMap{
	"cueStamp": secsToCueStamp,
}

var metaTmpl = metaTemplates{
	cue:     template.Must(template.New("cue").Funcs(funcs).Parse(cueTmpl)),
	ffchaps: template.Must(template.New("ffchaps").Funcs(funcs).Parse(ffChapTmpl)),
}

const cueTmpl = `FILE '{{.File}}' {{.Ext}}
{{- range $index, $ch := .Meta.Chapters}}
TRACK {{$index}} AUDIO
  TITLE "Chapter {{$index}}"
  INDEX 01 {{cueStamp $ch.StartToSeconds}}{{end}}`

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
