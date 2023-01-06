package media

import (
	"bufio"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/meta"
)

type Media struct {
	*avtools.Media
	Input  File
	Output File
	FFmeta File
	Cue    File
	Cover  File
}

func New(input string) *Media {
	m := avtools.NewMedia(input)
	return &Media{
		Media: m,
		Input: NewFile(input),
	}
}

func (m *Media) LoadIni(name string) *Media {
	file := NewFile(name)
	if IsPlainText(file.Mimetype) {
		contents, err := os.Open(file.Abs)
		if err != nil {
			log.Fatal(err)
		}
		defer contents.Close()

		scanner := bufio.NewScanner(contents)
		line := 0
		for scanner.Scan() {
			if line == 0 && scanner.Text() == meta.FFmetaComment {
				ini := meta.LoadIni(file.Abs)
				m.SetMeta(ini)
				m.FFmeta = file
				break
			} else {
				log.Fatalln("ffmpeg metadata files need to have ';FFMETADATA1' as the first line")
			}
		}
	}
	return m
}

func (m *Media) LoadCue(name string) *Media {
	file := NewFile(name)
	if IsPlainText(file.Mimetype) {
		cue := meta.LoadCueSheet(file.Abs)
		m.SetMeta(cue)
	}
	return m
}

func (m *Media) Probe() *Media {
	p := meta.FFProbe(m.Input.Abs)
	m.SetMeta(p)
	return m
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
	f.Name = strings.TrimSuffix(abs, f.Ext)
	f.name = strings.TrimSuffix(f.Base, f.Ext)

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
