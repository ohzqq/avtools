package fftools

import (
	"fmt"
	//"io/fs"
	"os"
	"log"
	"strings"
	"path/filepath"
)

func (m *Media) Split() {
	if m.HasChapters() {
		for i, chap := range *m.Meta.Chapters {
			m.Cut(chap.Start, chap.End, i)
		}
	}
}

func (m *Media) Cut(ss, to string, no int) {
	count := fmt.Sprintf("%06d", no + 1)
	cmd := NewCmd().In(m)
	timestamps := make(map[string]string)
	if ss != "" {
		timestamps["ss"] = ss
	}
	if to != "" {
		timestamps["to"] = to
	}
	cmd.Args().Post(timestamps).Out("tmp" + count).Ext(m.Ext)
	cmd.Run()
}

func Join(ext string) *FFmpegCmd {
	ff := NewCmd()
	pre := flagArgs{"f": "concat", "safe": "0"}
	ff.Args().Pre(pre).Ext(ext)
	files := find(ext)
	ff.tmpFile = concatFile(files)
	ff.In(NewMedia(ff.tmpFile.Name()))
	return ff
}

func concatFile(files []string) *os.File {
	file, err := os.CreateTemp("", "audiofiles")
	if err != nil { log.Fatal(err) }

	var cat strings.Builder
	for _, f := range files {
		cat.WriteString("file ")
		cat.WriteString("'")
		cat.WriteString(f)
		cat.WriteString("'")
		cat.WriteString("\n")
	}

		fmt.Println(cat.String())
	if _, err := file.WriteString(cat.String()); err != nil {
		log.Fatal(err)
	}

	return file
}

func find(ext string) []string {
	var files []string

	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range entries {
		if filepath.Ext(f.Name()) == ext {
			file, err := filepath.Abs(f.Name())
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, file)
		}
	}
	return files
}

