package ffmeta

import (
	"bytes"
	"io"
	"log"
	"os"
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

func (ff FFmeta) Dump() []byte {
	ini.PrettyFormat = false
	opts := ini.LoadOptions{
		IgnoreInlineComment:    true,
		AllowNonUniqueSections: true,
	}
	ffmeta := ini.Empty(opts)
	for k, v := range ff.Tags {
		_, err := ffmeta.Section("").NewKey(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(ff.Chapters.Each()) > 0 {
		for _, c := range ff.Chapters.Each() {
			sec, err := ffmeta.NewSection("CHAPTER")
			if err != nil {
				log.Fatal(err)
			}

			_, err = sec.NewKey("title", c.Title)
			if err != nil {
				log.Fatal(err)
			}

			_, err = sec.NewKey("START", c.Start().String())
			if err != nil {
				log.Fatal(err)
			}

			_, err = sec.NewKey("END", c.End().String())
			if err != nil {
				log.Fatal(err)
			}

			_, err = sec.NewKey("TIMEBASE", c.Timebase().String())
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	var buf bytes.Buffer
	_, err := ffmeta.WriteTo(&buf)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
func (ff FFmeta) Write(wr io.Writer) error {
	_, err := io.WriteString(wr, ffmetaComment)
	_, err = wr.Write(ff.Dump())
	if err != nil {
		return err
	}
	return nil
}

func (ff FFmeta) Save() error {
	return ff.SaveAs(ff.name)
}

func (ff FFmeta) SaveAs(name string) error {
	if name == "" || ff.name == "" {
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
