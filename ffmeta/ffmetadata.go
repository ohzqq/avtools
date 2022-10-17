package ffmeta

import (
	"io"
	"os"

	"github.com/ohzqq/avtools/chap"
)

const ffmetaComment = ";FFMETADATA1\n"

type FFmeta struct {
	chap.Chapters
	name string
	Tags map[string]string
}

func NewFFmeta() *FFmeta {
	return &FFmeta{Chapters: chap.NewChapters()}
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
