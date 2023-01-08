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
	return &Media{
		Media: m,
		Input: NewFile(input),
	}
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
	Abs      string
	AbsName  string
	Path     string
	Base     string
	Ext      string
	Name     string
	File     string
	Padding  string
	Mimetype string
	name     string
}

func NewFile(n string) File {
	abs, err := filepath.Abs(n)
	if err != nil {
		log.Fatal(err)
	}

	f := File{
		Base:    filepath.Base(abs),
		Ext:     filepath.Ext(abs),
		Abs:     abs,
		Padding: "%03d",
	}
	f.Mimetype = mime.TypeByExtension(f.Ext)
	f.Name = strings.TrimSuffix(f.Base, f.Ext)
	f.AbsName = strings.TrimSuffix(f.Abs, f.Ext)

	f.Path, f.File = filepath.Split(abs)

	return f
}

func (f File) WithExt(e string) string {
	return filepath.Join(f.Path, f.name+e)
}

func (f File) AddSuffix(s string) string {
	name := f.name + s + f.Ext
	return filepath.Join(f.Path, name)
}

func (f File) Pad(i int) string {
	p := fmt.Sprintf(f.Padding, i)
	return f.AddSuffix(p)
}

func (s File) Write(wr io.Writer, data []byte) error {
	_, err := wr.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (s File) Save(data []byte) error {
	file, err := os.Create(s.AbsName + s.Ext)
	if err != nil {
		return err
	}
	defer file.Close()

	err = s.Write(file, data)
	if err != nil {
		return err
	}

	return nil
}
