package media

import (
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/ff"
)

type Media struct {
	*avtools.Media
	Streams  []Stream
	Input    File
	Output   File
	Ini      File
	Cue      File
	Cover    File
	HasCover bool
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

func (m *Media) SetMeta(name string) *Media {
	file := NewFile(name)
	switch file.Ext {
	case ".cue":
		m.LoadCue(name)
	case ".ini":
		m.LoadIni(name)
	}
	return m
}

func (m Media) Command() ff.Cmd {
	cmd := ff.New()
	cmd.In(m.Input.Abs)
	return cmd
}

func (m Media) HasChapters() bool {
	return len(m.Chapters) > 0
}

func (m Media) AudioStreams() []Stream {
	var streams []Stream
	for _, stream := range m.Streams {
		if stream.CodecType == "audio" {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (m Media) VideoStreams() []Stream {
	var streams []Stream
	for _, stream := range m.Streams {
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
}

func NewFile(n string) File {
	abs, err := filepath.Abs(n)
	if err != nil {
		log.Fatal(err)
	}

	f := File{
		Base:     filepath.Base(abs),
		Abs:      abs,
		FileName: &FileName{},
	}

	f.Padding = "%03d"
	f.Ext = filepath.Ext(abs)
	f.Mimetype = mime.TypeByExtension(f.Ext)
	f.Name = strings.TrimSuffix(f.Base, f.Ext)

	f.Path, f.File = filepath.Split(abs)

	return f
}

func (f File) NewName() *FileName {
	name := &FileName{
		Name:    f.Name,
		Path:    f.Path,
		Padding: f.Padding,
	}
	return name
}

func (f *FileName) WithExt(e string) *FileName {
	//name := filepath.Join(f.Path, f.Name+e)
	f.Ext = e
	return f
}

func (f *FileName) Suffix(suf string) *FileName {
	//name := f.Name + suf + f.Ext
	//return filepath.Join(f.Path, name)
	f.Name = f.Name + suf
	return f
}

func (f *FileName) Prefix(pre string) *FileName {
	f.Name = pre + f.Name
	//return filepath.Join(f.Path, name)
	return f
}

func (f *FileName) Pad(i int) *FileName {
	p := fmt.Sprintf(f.Padding, i)
	return f.Suffix(p)
}

func (f FileName) Join() string {
	return filepath.Join(f.Path, f.Name+f.Ext)
}

func (f FileName) Write(wr io.Writer, data []byte) error {
	_, err := wr.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (f FileName) Run() error {
	file, err := os.Create(f.Join())
	if err != nil {
		return err
	}
	defer file.Close()

	err = f.Write(file, f.data)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileName) Save(data []byte) {
	f.data = data
}
