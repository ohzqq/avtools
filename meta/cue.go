package meta

import (
	"bufio"
	"bytes"
	"html/template"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ohzqq/avtools"
	"github.com/samber/lo"
)

type CueSheet struct {
	file   string
	Audio  string
	Tracks []*avtools.Chapter
}

func LoadCueSheet(file string) *CueSheet {
	var sheet CueSheet

	contents, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	var times []time.Duration
	var titles []string
	scanner := bufio.NewScanner(contents)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.Contains(line, "TITLE"):
			titles = append(titles, title(line))
		case strings.Contains(line, "INDEX 01"):
			times = append(times, start(line))
		}
	}

	e := 1
	for i := 0; i < len(titles); i++ {
		track := &avtools.Chapter{}
		track.Title = titles[i]
		track.Start = avtools.Timestamp(times[i])
		if e < len(titles) {
			track.End = avtools.Timestamp(times[e])
		}
		e++
		sheet.Tracks = append(sheet.Tracks, track)
	}

	return &sheet
}

func (cue CueSheet) Chapters() []*avtools.Chapter {
	return cue.Tracks
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

func start(line string) time.Duration {
	stamp := strings.TrimPrefix(line, "INDEX 01 ")
	split := strings.Split(stamp, ":")
	split = lo.DropRight(split, 1)
	stamp = strings.Join(split, ":")
	s := avtools.ParseStamp(stamp)
	return s
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
