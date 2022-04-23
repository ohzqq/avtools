package fftools

import (
	"log"
	"os"
	"bufio"
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"time"
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

func (c *Chapter) timeBaseFloat() float64 {
	tb := strings.ReplaceAll(c.Timebase, "1/", "")
	baseint, _ := strconv.ParseFloat(tb, 64)
	return baseint
}

func (c *Chapter) toSeconds() () {
	tb := c.timeBaseFloat()
	ss, _ := strconv.ParseFloat(c.Start, 64)
	to, _ := strconv.ParseFloat(c.End, 64)
	c.Start = strconv.FormatFloat(ss / tb, 'f', 6, 64)
	c.End = strconv.FormatFloat(to / tb, 'f', 6, 64)
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
	params := newFlagArg("f", "ffmetadata")
	cmd.Args().PostInput(params).ACodec("none").VCodec("none").Ext("ini")
	fmt.Printf("%v", cmd.Cmd().String())
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
			c.toSeconds()
			if err != nil { log.Fatal(err) }
			chapters = append(chapters, &c)
		}
	}

	return &MediaMeta{
		Chapters: &chapters,
		Tags: &tags,
	}
}

func ReadCueSheet(file string) *MediaMeta {
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
			start := strings.TrimPrefix(s, "INDEX 01 ")
			rmFrames := regexp.MustCompile(`:\d\d$`)
			start = rmFrames.ReplaceAllString(start, "s")
			start = strings.ReplaceAll(start, ":", "m")
			dur, _ := time.ParseDuration(start)
			durS := strconv.Itoa(int(dur.Seconds()))
			indices = append(indices, durS)
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

	return &MediaMeta{Chapters: &tracks}
}
