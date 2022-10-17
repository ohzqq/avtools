package ffmeta

import (
	"bytes"
	"log"

	"gopkg.in/ini.v1"
)

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
