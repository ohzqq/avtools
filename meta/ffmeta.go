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

const FFmetaComment = ";FFMETADATA1\n"

type FFMeta struct {
	tags  map[string]string
	chaps []*avtools.Chapter
}

type FFMetaChapter struct {
	Base  string `ini:"TIMEBASE"`
	Start int    `ini:"START"`
	End   int    `ini:"END"`
	Title string `ini:"title"`
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

			c := &avtools.Chapter{
				Start: avtools.Timestamp(avtools.ParseStampDuration(chap.Start, chap.Timebase())),
				End:   avtools.Timestamp(avtools.ParseStampDuration(chap.End, chap.Timebase())),
				Title: chap.Title,
			}
			ffmeta.chaps = append(ffmeta.chaps, c)
		}
	}

	return ffmeta
}

func (ff FFMeta) Chapters() []*avtools.Chapter {
	return ff.chaps
}

func (ff FFMeta) Tags() map[string]string {
	return ff.tags
}

func (ff FFMeta) Streams() []map[string]string {
	return []map[string]string{}
}

func (ch FFMetaChapter) Timebase() int {
	if tb := ch.Base; tb != "" {
		ch.Base = strings.TrimPrefix(tb, "1/")
	}
	baseFloat, _ := strconv.Atoi(ch.Base)
	return baseFloat
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
	_, err := buf.WriteString(FFmetaComment)
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
