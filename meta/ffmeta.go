package meta

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ohzqq/avtools"
	"gopkg.in/ini.v1"
)

const ffmetaComment = ";FFMETADATA1\n"

type FFMeta struct {
	tags     map[string]string
	chapters []avtools.ChapterMeta
}

type FFMetaChapter struct {
	Base         string  `ini:"TIMEBASE"`
	StartTime    float64 `ini:"START"`
	EndTime      float64 `ini:"END"`
	ChapterTitle string  `ini:"title"`
}

//func (ff FFmeta) Each() []FF

func (ff FFMeta) Chapters() []avtools.ChapterMeta {
	return ff.chapters
}

func (ff FFMeta) Tags() map[string]string {
	return ff.tags
}

func (ff FFMeta) Streams() []map[string]string {
	return []map[string]string{}
}

func (ch FFMetaChapter) Start() float64 {
	return ch.StartTime
}

func (ch FFMetaChapter) End() float64 {
	return ch.EndTime
}

func (ch FFMetaChapter) Title() string {
	return ch.ChapterTitle
}

func (ch FFMetaChapter) Timebase() float64 {
	if tb := ch.Base; tb != "" {
		ch.Base = strings.TrimPrefix(tb, "1/")
	}
	baseFloat, _ := strconv.ParseFloat(ch.Base, 64)
	return baseFloat
}

func LoadIni(input string) *FFMeta {
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

	ffmeta := &FFMeta{}
	ffmeta.tags = f.Section("").KeysHash()

	if f.HasSection("chapter") {
		sections, err := f.SectionsByName("chapter")
		if err != nil {
			log.Fatal(err)
		}

		for _, ch := range sections {
			var chap FFMetaChapter
			ch.MapTo(&chap)
			ffmeta.chapters = append(ffmeta.chapters, chap)
		}
	}

	return ffmeta
}

func (ff FFMeta) Dump() []byte {
	ini.PrettyFormat = false

	opts := ini.LoadOptions{
		IgnoreInlineComment:    true,
		AllowNonUniqueSections: true,
	}

	ffmeta := ini.Empty(opts)

	for k, v := range ff.tags {
		_, err := ffmeta.Section("").NewKey(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	var buf bytes.Buffer
	_, err := buf.WriteString(ffmetaComment)
	_, err = ffmeta.WriteTo(&buf)
	if err != nil {
		log.Fatal(err)
	}

	//_, err = buf.Write(ff.IniChaps())
	//if err != nil {
	//  log.Fatal(err)
	//}

	return buf.Bytes()
}

func (ff FFMeta) Write(wr io.Writer) error {
	_, err := wr.Write(ff.Dump())
	if err != nil {
		return err
	}
	return nil
}

func (ff FFMeta) Save() error {
	//return ff.SaveAs(ff.name)
	return nil
}

func (ff FFMeta) SaveAs(name string) error {
	if name == "" {
		name = "tmp"
	}

	file, err := os.Create(name + ".ini")
	if err != nil {
		return err
	}
	defer file.Close()

	err = ff.Write(file)
	if err != nil {
		return err
	}

	return nil
}

var metaTags = []string{
	"title",
	"album",
	"artist",
	"composer",
	"date",
	"year",
	"genre",
	"comment",
	"album_artist",
	"track",
	"language",
	"lyrics",
}
