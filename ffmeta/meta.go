package ffmeta

import (
	"bytes"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/ohzqq/avtools/chap"
	"gopkg.in/ini.v1"
)

type FFmeta struct {
	chap.Chapters
	name    string
	Streams []*Stream
	Format  `json:"format"`
	Chaps   []Chapter `json:"chapters"`
}

type Stream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
}

type Format struct {
	Filename string
	Dur      duration `json:"duration"`
	Size     string
	BitRate  string `json:"bit_rate"`
	Tags     map[string]string
}

type duration string

func NewFFmeta() *FFmeta {
	return &FFmeta{Chapters: chap.NewChapters()}
}

func (ff FFmeta) Duration() duration {
	return ff.Dur
}

func (d duration) String() string {
	return string(d)
}

func (d duration) Int() int {
	return int(math.Round(d.Float()))
}

func (d duration) Float() float64 {
	f, err := strconv.ParseFloat(d.String(), 64)
	if err != nil {
		return 0
	}
	return f
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

	var buf bytes.Buffer
	_, err := ffmeta.WriteTo(&buf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = buf.Write(ff.IniChaps())
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
