package media

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/ff"
	"github.com/ohzqq/avtools/meta"
)

type Media struct {
	*avtools.Media
	streams     []Stream
	Input       File
	Output      File
	Ini         File
	Cue         File
	Cover       File
	HasCover    bool
	MetaChanged bool
}

type Stream struct {
	CodecType string
	CodecName string
	Index     string
	IsCover   bool
}

func New(input string) *Media {
	m := avtools.NewMedia(input)
	med := &Media{
		Media: m,
		Input: NewFile(input),
	}
	med.Output = File{FileName: med.Input.NewName()}
	med.Probe()

	return med
}

func (m *Media) LoadMeta(name string) *Media {
	file := NewFile(name)
	switch {
	case file.IsFFMeta():
		ini := meta.LoadIni(file.Abs)
		m.Media.SetMeta(ini)
		m.Ini = file
	case file.IsCue():
		cue := meta.LoadCueSheet(file.Abs)
		m.Media.SetMeta(cue)
		dur := m.GetTag("duration")
		last := m.Chapters()[len(m.Chapters())-1]
		last.End = avtools.Timestamp(avtools.ParseStamp(dur))
		m.Cue = file
	}
	m.MetaChanged = true
	return m
}

func (m Media) Command() ff.Cmd {
	cmd := ff.New()
	cmd.In(m.Input.Abs)
	return cmd
}

func (m Media) HasChapters() bool {
	return len(m.Chapters()) > 0
}

func (m Media) AudioStreams() []Stream {
	var streams []Stream
	for _, stream := range m.streams {
		if stream.CodecType == "audio" {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (m Media) VideoStreams() []Stream {
	var streams []Stream
	for _, stream := range m.streams {
		if stream.CodecType == "video" {
			streams = append(streams, stream)
		}
	}
	return streams
}

func IsPlainText(mtype string) bool {
	if strings.Contains(mtype, "text/plain") {
		return true
	} else {
		log.Fatalln("needs to be plain text file")
	}
	return false
}

type File struct {
	*FileName
	Abs      string
	Base     string
	File     string
	Mimetype string
}

type FileName struct {
	Path    string
	Ext     string
	Name    string
	Padding string
	data    []byte
	file    *os.File
}

func NewFile(n string) File {
	abs, err := filepath.Abs(n)
	if err != nil {
		log.Fatal(err)
	}

	f := File{
		Base:     filepath.Base(abs),
		Abs:      abs,
		FileName: NewFileName(),
	}

	f.Ext = filepath.Ext(abs)
	f.Mimetype = mime.TypeByExtension(f.Ext)
	f.Name = strings.TrimSuffix(f.Base, f.Ext)

	f.Path, f.File = filepath.Split(abs)

	return f
}

func NewFileName() *FileName {
	name := &FileName{
		Padding: "%03d",
	}
	return name
}

func (f File) NewName() *FileName {
	name := &FileName{
		Name:    f.Name,
		Path:    f.Path,
		Padding: f.Padding,
		//Ext:     f.Ext,
	}
	return name
}

func (f File) IsFFMeta() bool {
	if IsPlainText(f.Mimetype) {
		contents, err := os.Open(f.Abs)
		if err != nil {
			log.Fatal(err)
		}
		defer contents.Close()

		scanner := bufio.NewScanner(contents)
		line := 0
		for scanner.Scan() {
			if line == 0 && scanner.Text() == ";FFMETADATA1" {
				return true
				break
			}
		}
	}
	return false
}

func (f File) IsCue() bool {
	if IsPlainText(f.Mimetype) {
		return f.Ext == ".cue"
	}
	return false
}

func (f File) IsImage() bool {
	return strings.Contains(f.Mimetype, "image")
}

func (f *FileName) WithExt(e string) *FileName {
	f.Ext = e
	return f
}

func (f *FileName) Suffix(suf string) *FileName {
	f.Name = f.Name + suf
	return f
}

func (f *FileName) Prefix(pre string) *FileName {
	f.Name = pre + f.Name
	return f
}

func (f *FileName) Pad(i int) *FileName {
	p := fmt.Sprintf(f.Padding, i)
	return f.Suffix(p)
}

func (f FileName) Join() string {
	return filepath.Join(f.Path, f.Name+f.Ext)
}

func (f FileName) Write(wr io.Writer) error {
	_, err := wr.Write(f.data)
	if err != nil {
		return err
	}
	return nil
}

func (f FileName) Run() error {
	if f.file != nil {
		defer f.file.Close()
	}

	err := f.Write(f.file)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileName) Tmp(data []byte) {
	file, err := os.CreateTemp("", f.Name)
	if err != nil {
		log.Fatal(err)
	}
	f.file = file
	f.data = data
}

func (f *FileName) Save(data []byte) {
	file, err := os.Create(f.Join())
	if err != nil {
		log.Fatal(err)
	}
	f.file = file
	f.data = data
}
