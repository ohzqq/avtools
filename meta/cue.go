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

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/timestamp"
)

type CueSheet struct {
	file       string
	Audio      string
	Tracks     []CueTrack
	titles     []string
	startTimes []float64
}

type CueTrack struct {
	title string
	start float64
	end   float64
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

func (cue CueSheet) Chapters() []avtools.ChapterMeta {
	var chaps []avtools.ChapterMeta
	for _, track := range cue.Tracks {
		chaps = append(chaps, track)
	}
	return chaps
}

func (cue CueSheet) Tags() map[string]string {
	return map[string]string{
		"filename": cue.Audio,
	}
}

func (cue CueSheet) Streams() []map[string]string {
	return []map[string]string{}
}

func NewCueSheet(f string) *CueSheet {
	return &CueSheet{file: f}
}

func NewTrack() CueTrack {
	return CueTrack{}
}

func (s *CueSheet) SetAudio(name string) *CueSheet {
	s.Audio = name
	return s
}

func (s CueSheet) File() string {
	return s.Audio
}

func (s CueSheet) Ext() string {
	//ext := strings.TrimPrefix(s.Audio.Ext, ".")
	//return strings.ToUpper(ext)
	return ""
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

func start(line string) float64 {
	stamp := strings.TrimPrefix(line, "INDEX 01 ")
	return timestamp.ParseHHSS(stamp)
}

func (s CueSheet) Dump() []byte {
	var (
		tmpl = template.Must(template.New("cue").Funcs(tmplFuncs).Parse(cueTmpl))
		buf  bytes.Buffer
	)

	if s.Audio == "" {
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
	return s.SaveAs(s.Audio)
}

func (s CueSheet) SaveAs(name string) error {
	if name == "" || s.Audio == "" {
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

func (t CueTrack) Start() float64 {
	return t.start
}

func (t CueTrack) Stamp() string {
	mm := int(t.start) / 60
	ss := int(t.start) % 60
	start := fmt.Sprintf("%02d:%02d:00", mm, ss)
	return start
}

func (t *CueTrack) SetStart(secs float64) *CueTrack {
	t.start = secs
	return t
}

func (t CueTrack) End() float64 {
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
