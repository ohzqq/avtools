package fftools

import (
	"fmt"
	"io/fs"
	"os"
	"log"
	"strings"
	"path/filepath"
)

func (m *Media) Split() {
	cmd := NewCmd().In(m)
	ch := cmd.GetChapters()
	for i, chap := range *ch {
		m.Cut(chap.Start, chap.End, i)
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
	cmd.Args().PostInput(timestamps).Out("tmp" + count).Ext(m.Ext)
	cmd.Run()
}

func Join(ext string) *FFmpegCmd {
	ff := NewCmd()
	files := find(ext)
	ff.tmpFile = concatFile(files)
	ff.In(NewMedia(ff.tmpFile.Name()))
	pre := flagArgs{"f": "concat", "safe": "0"}
	ff.Args().PreInput(pre).Ext(ext)
	return ff
}

func concatFile(files []string) *os.File {
	file, err := os.CreateTemp("", "audiofiles")
	if err != nil { log.Fatal(err) }

	var cat strings.Builder
	for _, f := range files {
		abs, err := filepath.Abs(f)
		if err != nil { log.Fatal(err) }
		cat.WriteString("file ")
		cat.WriteString("'")
		cat.WriteString(abs)
		cat.WriteString("'")
		cat.WriteString("\n")
	}

	if _, err := file.WriteString(cat.String()); err != nil {
		log.Fatal(err)
	}

	return file
}

func find(ext string) []string {
	var files []string
	filepath.WalkDir(".", func(file string, dir fs.DirEntry, e error) error {
		if e != nil { return e }
		if filepath.Ext(dir.Name()) == ext {
			files = append(files, file)
		}
		return nil
	})
	return files
}

