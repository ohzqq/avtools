package tool

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ohzqq/avtools/file"
)

type RelatedFiles map[string]MediaFile

type MediaFile struct {
	file.File
}

func (r RelatedFiles) Get(name string) MediaFile {
	if r.Has(name) {
		return r[name]
	}
	return MediaFile{}
}

func (r RelatedFiles) Has(name string) bool {
	if f, ok := r[name]; ok && f != (MediaFile{}) {
		return true
	}
	return false
}

func (f MediaFile) IsImage() bool {
	if strings.Contains(f.Mimetype, "image") {
		return true
	}
	return false
}

func (f MediaFile) IsAudio() bool {
	if strings.Contains(f.Mimetype, "audio") {
		return true
	} else {
		fmt.Println("not an audio file")
	}
	return false
}

func (f MediaFile) IsPlainText() bool {
	if strings.Contains(f.Mimetype, "text/plain") {
		return true
	} else {
		log.Fatalln("needs to be plain text file")
	}
	return false
}

func (f MediaFile) IsFFmeta() bool {
	if f.IsPlainText() {
		contents, err := os.Open(f.Abs)
		if err != nil {
			log.Fatal(err)
		}
		defer contents.Close()

		scanner := bufio.NewScanner(contents)
		line := 0
		for scanner.Scan() {
			if line == 0 && scanner.Text() == ";FFMETADATA1" {
				return true
			} else {
				log.Fatalln("ffmpeg metadata files need to have ';FFMETADATA1' as the first line")
			}
		}
	}
	return false
}
