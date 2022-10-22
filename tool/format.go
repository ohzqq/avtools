package tool

import (
	"bufio"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type RelatedFiles map[string]FileFormat

type FileFormat string

func (r RelatedFiles) Get(name string) string {
	if r.Has(name) {
		return r[name].String()
	}
	return ""
}

func (r RelatedFiles) Has(name string) bool {
	if f, ok := r[name]; ok && f != "" {
		return true
	}
	return false
}

func (f FileFormat) String() string {
	return string(f)
}

func (f FileFormat) Ext() string {
	ext := filepath.Ext(f.String())
	return ext
}

func (f FileFormat) Mimetype() string {
	return mime.TypeByExtension(f.Ext())
}

func (f FileFormat) IsImage() bool {
	if strings.Contains(f.Mimetype(), "image") {
		return true
	}
	return false
}

func (f FileFormat) IsAudio() bool {
	if strings.Contains(f.Mimetype(), "audio") {
		return true
	} else {
		fmt.Println("not an audio file")
	}
	return false
}

func (f FileFormat) IsPlainText() bool {
	if strings.Contains(f.Mimetype(), "text/plain") {
		return true
	} else {
		log.Fatalln("needs to be plain text file")
	}
	return false
}

func (f FileFormat) IsFFmeta() bool {
	if f.IsPlainText() {
		contents, err := os.Open(f.String())
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
