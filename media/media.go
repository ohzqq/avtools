package media

import (
	"bufio"
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
}

func New(input string) *Media {
	m := avtools.NewMedia(input)
	return &Media{Media: m}
}

func (m *Media) LoadIni(name string) *Media {
	abs, err := filepath.Abs(name)
	if err != nil {
		log.Fatal(err)
	}
	if IsPlainText(name) {
		contents, err := os.Open(abs)
		if err != nil {
			log.Fatal(err)
		}
		defer contents.Close()

		scanner := bufio.NewScanner(contents)
		line := 0
		for scanner.Scan() {
			if line == 0 && scanner.Text() == ";FFMETADATA1" {
				ini := meta.LoadIni(name)
				m.SetMeta(ini)
				break
			} else {
				log.Fatalln("ffmpeg metadata files need to have ';FFMETADATA1' as the first line")
			}
		}
	}
	return m
}

func (m *Media) LoadCue(name string) *Media {
	cue := meta.LoadCueSheet(name)
	m.SetMeta(cue)
	return m
}

func (m *Media) Probe() *Media {
	p := meta.FFProbe(m.Filename)
	m.SetMeta(p)
	return m
}

func IsPlainText(file string) bool {
	ext := filepath.Ext(file)
	mt := mime.TypeByExtension(ext)
	if strings.Contains(mt, "text/plain") {
		return true
	} else {
		log.Fatalln("needs to be plain text file")
	}
	return false
}
