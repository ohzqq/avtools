package avtools

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

const ffProbeMeta = `format=filename,start_time,duration,size,bit_rate:stream=codec_type,codec_name:format_tags`

type MediaMeta struct {
	Chapters Chapters
	Streams  []*Stream
	Format   *Format
	Tags     Tags
	//Tags     map[string]string
}

func (m *MediaMeta) SetChapters(ch Chapters) {
	m.Chapters = ch
}

func (m *MediaMeta) LastChapterEnd() {
	if m.Format.Duration != "" && len(m.Chapters) > 0 {
		lastCh := m.Chapters[len(m.Chapters)-1]
		lastCh.End = m.Format.DurationSecs(lastCh.TimebaseFloat())
	}
}

func (m *MediaMeta) HasChapters() bool {
	if len(m.Chapters) != 0 {
		return true
	}
	return false
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
	Tags            //map[string]string
}

func (f Format) Ext() string {
	if f.Filename != "" {
		return strings.TrimPrefix(path.Ext(f.Filename), ".")
	}
	return ""
}

func (f Format) DurationSecs(timebase float64) int {
	seconds := decimalSecsToFloat(f.Duration) * timebase
	return int(seconds)
}

type Tags struct {
	Title    string `json:"title",ini:"title"`
	Artist   string `json:"artist",ini:"artist"`
	Composer string `json:"composer",ini:"composer"`
	Album    string `json:"album",ini:"album"`
	Comment  string `json:"comment",ini:"comment"`
	Genre    string `json:"genre",ini:"genre"`
}

type Chapters []*Chapter

type Chapter struct {
	Timebase string `json:"time_base",ini:"timebase"`
	Start    int    `json:"start",ini:"start"`
	End      int    `json:"end",ini:"end"`
	//Tags     map[string]string
	Title string `ini:"title"`
}

func (c *Chapter) StartToIntString() string {
	result := float64(c.Start) * c.TimebaseFloat()
	return strconv.FormatFloat(result, 'f', 0, 64)
}

func (c *Chapter) CueStamp() string {
	sec := float64(c.Start) / c.TimebaseFloat()
	m := int(sec) / 60
	s := int(sec) % 60
	return fmt.Sprintf("%02d:%02d:00", m, s)
	//return s
}

func (c *Chapter) StartToSeconds() string {
	if c.Start == 0 {
		return "0"
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
		return "0"
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

func (m *MediaMeta) ToFFmeta() {
	m.LastChapterEnd()
	tmpl, err := GetTmpl("ffmeta")
	if err != nil {
		log.Println("executing template:", err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, m)
	if err != nil {
		log.Println("executing template:", err)
	}
	println(buf.String())
	//if len(m.Chapters) != 0 {
	//}
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
		Format: &Format{},
	}
	err = f.Section("").MapTo(&media.Format.Tags)
	if err != nil {
		log.Fatal(err)
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

func LoadCueSheet(file string) *MediaMeta {
	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	var (
		titles     []string
		startTimes []int
		meta       = MediaMeta{Format: &Format{}}
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
		//t := new(Chapter)
		t.Title = titles[i]
		t.Start = startTimes[i]
		if e < len(titles) {
			t.End = startTimes[e]
		}
		e++
		meta.Chapters = append(meta.Chapters, &t)
	}

	return &meta
}
