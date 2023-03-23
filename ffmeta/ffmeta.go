package ffmeta

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/fidi"
	"gopkg.in/ini.v1"
)

const FFmetaComment = ";FFMETADATA1\n"

type FFMeta struct {
	fidi.File
	tags     map[string]string
	chapters []avtools.ChapterMeta
}

func Load(input string) (avtools.Metaz, error) {
	opts := ini.LoadOptions{}
	opts.Insensitive = true
	opts.InsensitiveSections = true
	opts.IgnoreInlineComment = true
	opts.AllowNonUniqueSections = true

	ffmeta := &FFMeta{}
	ffmeta.File = fidi.NewFile(input)

	if !IsFFMeta(ffmeta.File) {
		return ffmeta, fmt.Errorf("not an ffmetadata file")
	}

	f, err := ini.LoadSources(opts, ffmeta.Path())
	if err != nil {
		return ffmeta, err
	}

	ffmeta.tags = f.Section("").KeysHash()

	if f.HasSection("chapter") {
		sections, err := f.SectionsByName("chapter")
		if err != nil {
			return ffmeta, err
		}

		for _, ch := range sections {
			var c FFMetaChapter
			err := ch.MapTo(&c)
			if err != nil {
				return ffmeta, err
			}
			ffmeta.chapters = append(ffmeta.chapters, c)
		}
	}

	return ffmeta, nil
}

func Dump(meta avtools.Meta) []byte {
	ini.PrettyFormat = false

	opts := ini.LoadOptions{
		IgnoreInlineComment:    true,
		AllowNonUniqueSections: true,
	}

	ffmeta := ini.Empty(opts)

	for k, v := range meta.Tags() {
		_, err := ffmeta.Section("").NewKey(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, chapter := range meta.Chapters() {
		sec, err := ffmeta.NewSection("CHAPTER")
		if err != nil {
			log.Fatal(err)
		}
		sec.NewKey("TIMEBASE", "1/1000")
		ss := strconv.Itoa(int(chapter.StartStamp.Dur.Milliseconds()))
		sec.NewKey("START", ss)

		to := strconv.Itoa(int(chapter.EndStamp.Dur.Milliseconds()))
		sec.NewKey("END", to)
		sec.NewKey("title", chapter.ChapTitle)
		for k, v := range chapter.Tags {
			sec.NewKey(k, v)
		}
	}

	var buf bytes.Buffer
	_, err := buf.WriteString(FFmetaComment)
	_, err = ffmeta.WriteTo(&buf)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (ff FFMeta) Chapters() []avtools.ChapterMeta {
	return ff.chapters
}

func (ff FFMeta) Tags() map[string]string {
	return ff.tags
}

func (ff FFMeta) Streams() []map[string]string {
	return []map[string]string{}
}

func (ff FFMeta) Source() fidi.File {
	return ff.File
}

func IsFFMeta(f fidi.File) bool {
	if err := avtools.IsPlainText(f.Mime); err != nil {
		return false
	}
	contents, err := os.Open(f.Path())
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
	return false
}
