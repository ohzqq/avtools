package ffmeta

import (
	"bytes"
	"log"
	"path/filepath"
	"strconv"

	"github.com/ohzqq/avtools"
	"gopkg.in/ini.v1"
)

const FFmetaComment = ";FFMETADATA1\n"

type FFMeta struct {
	tags     map[string]string
	chapters []avtools.ChapterMeta
}

func Load(input string) *FFMeta {
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
			var c FFMetaChapter
			err := ch.MapTo(&c)
			if err != nil {
				log.Fatal(err)
			}
			ffmeta.chapters = append(ffmeta.chapters, c)
		}
	}

	return ffmeta
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
