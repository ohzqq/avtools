package meta

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type CueSheet struct {
	file       string
	Audio      string
	Tracks     []CueTrack
	titles     []string
	startTimes []int
}

type CueTrack struct {
	title string
	start int
	end   int
}

func NewCueSheet(f string) *CueSheet {
	return &CueSheet{file: f}
}

func NewTrack() CueTrack {
	return CueTrack{}
}

func LoadCueSheet(file string) *CueSheet {
	var sheet CueSheet

	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	scanner := bufio.NewScanner(contents)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.Contains(line, "TITLE"):
			sheet.titles = append(sheet.titles, title(line))
		case strings.Contains(line, "INDEX 01"):
			sheet.startTimes = append(sheet.startTimes, start(line))
		}
	}

	e := 1
	for i := 0; i < len(sheet.titles); i++ {
		var t CueTrack
		t.title = sheet.titles[i]
		t.start = sheet.startTimes[i]
		if e < len(sheet.titles) {
			t.end = sheet.startTimes[e]
		}
		e++
		sheet.Tracks = append(sheet.Tracks, t)
	}

	return &sheet
}

func (s *CueSheet) SetAudio(name string) *CueSheet {
	s.Audio = name
	return s
}

func (s CueSheet) File() string {
	return s.Audio.Abs
}

func (s CueSheet) Ext() string {
	ext := strings.TrimPrefix(s.Audio.Ext, ".")
	return strings.ToUpper(ext)
}

func cueFile(line string) string {
	fileRegexp := regexp.MustCompile(`^(\w+ )('|")(?P<title>.*)("|')( .*)$`)
	matches := fileRegexp.FindStringSubmatch(line)
	title := matches[fileRegexp.SubexpIndex("title")]
	return title
}

func title(line string) string {
	t := strings.TrimPrefix(line, "TITLE ")
	t = strings.Trim(t, "'")
	t = strings.Trim(t, `"`)
	return t
}

func start(line string) int {
	stamp := strings.TrimPrefix(line, "INDEX 01 ")
	split := strings.Split(stamp, ":")
	dur, err := time.ParseDuration(split[0] + "m" + split[1] + "s")
	if err != nil {
		log.Fatal(err)
	}
	return int(dur.Seconds())
}

func (s CueSheet) Dump() []byte {
	var (
		tmpl = template.Must(template.New("cue").Funcs(tmplFuncs).Parse(cueTmpl))
		buf  bytes.Buffer
	)

	if s.Audio.Abs == "" {
		//s.Audio = file.New("tmp")
	}

	err := tmpl.Execute(&buf, s)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (s CueSheet) Write(wr io.Writer) error {
	_, err := wr.Write(s.Dump())
	if err != nil {
		return err
	}
	return nil
}

func (s CueSheet) Save() error {
	return s.SaveAs(s.Audio.Abs)
}

func (s CueSheet) SaveAs(name string) error {
	if name == "" || s.Audio.Abs == "" {
		name = "tmp"
	}

	file, err := os.Create(name + ".cue")
	if err != nil {
		return err
	}
	defer file.Close()

	err = s.Write(file)
	if err != nil {
		return err
	}

	return nil
}

// cue tracks
func (t CueTrack) Title() string {
	return t.title
}

func (t *CueTrack) SetTitle(title string) *CueTrack {
	t.title = title
	return t
}

func (t CueTrack) Start() int {
	return t.start
}

func (t CueTrack) Stamp() string {
	mm := t.start / 60
	ss := t.start % 60
	start := fmt.Sprintf("%02d:%02d:00", mm, ss)
	return start
}

func (t *CueTrack) SetStart(secs int) *CueTrack {
	t.start = secs
	return t
}

func (t CueTrack) End() int {
	return t.end
}

func (t CueTrack) Timebase() float64 {
	return 1
}

var tmplFuncs = template.FuncMap{
	"inc": Inc,
}

func Inc(n int) int {
	return n + 1
}

const cueTmpl = `FILE "{{.File}}" {{.Ext -}}
{{range $idx, $ch := .Tracks}}
TRACK {{inc $idx}} AUDIO
{{- if eq $ch.Title ""}}
  TITLE "Chapter {{inc $idx}}"
{{- else}}
  TITLE "{{$ch.Title}}"
{{- end}}
  INDEX 01 {{$ch.Stamp}}
{{- end -}}`
