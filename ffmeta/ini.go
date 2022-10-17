package ffmeta

import (
	"log"
	"path/filepath"

	"github.com/ohzqq/avtools/chap"
	"gopkg.in/ini.v1"
)

const ffmetaComment = ";FFMETADATA1\n"

func LoadIni(input string) *FFmeta {
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

	ffmeta := NewFFmeta()
	ffmeta.Tags = f.Section("").KeysHash()

	if f.HasSection("chapter") {
		sec, _ := f.SectionsByName("chapter")
		for _, chapter := range sec {
			c := Chapter{}
			err := chapter.MapTo(&c)
			if err != nil {
				log.Fatal(err)
			}
			ch := chap.NewChapter().SetMeta(c)
			ffmeta.Chapters.Chapters = append(ffmeta.Chapters.Chapters, ch)
		}
	}

	return ffmeta
}
