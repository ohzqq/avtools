package fftools

import (
	"os"
	"log"
	//"fmt"
	"strings"
	"io/fs"
	"path/filepath"

	//"github.com/alessio/shellescape"
)

func (ff *FFmpegCmd) Join(ext string) *FFmpegCmd {
	files := find("." + ext)
	fileList := concatFile(files)
	//defer os.Remove(fileList.Name())
	ff.In(fileList.Name())
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

