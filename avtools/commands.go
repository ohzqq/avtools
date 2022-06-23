package avtools

import (
	"fmt"
	//"io/fs"
	"os"
	"os/exec"
	"log"
	"text/template"
	"bytes"
	"strings"
	//"strconv"
	"path/filepath"
)

func(m *Media) AddAlbumArt(cover string) *FFmpegCmd {
	path, err := filepath.Abs(cover)
	if err != nil {
		log.Fatal(err)
	}
	switch codec := m.AudioCodec(); codec {
	case "aac":
		_, err := exec.LookPath("AtomicParsley")
		if err != nil {
			log.Fatal("embedding album art with aac requires AtomicParsley to be installed")
		}
		addAacCover(m.Path, path)
		return nil
	case "mp3":
		cmd := NewCmd().In(m)
		cmd.Args().Out("tmp-cover").Params(Mp3CoverArgs()).Cover(path)
		return cmd
	}
	return nil
}

func addAacCover(media, cover string) {
	cmd := exec.Command("AtomicParsley", media, "--artwork", cover, "--overWrite")
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		fmt.Printf("%v\n", stderr.String())
	}
}

func(m *Media) AddFFmeta(meta string) *FFmpegCmd {
	path, err := filepath.Abs(meta)
	if err != nil {
		log.Fatal(err)
	}
	cmd := NewCmd().In(m)
	cmd.Args().Meta(path)
	return cmd
}

func(m *Media) Remove(chaps, cover, meta bool) *FFmpegCmd {
	cmd := NewCmd().In(m)
	if chaps {
		cmd.Args().Post("map_chapters", "-1")
	}
	if cover {
		cmd.Args().VCodec("vn")
	}
	if meta {
		cmd.Args().Post("map_metadata", "-1")
	}
	return cmd
}

func(m *Media) Extract(chaps, cover, meta bool) {
	cmd := NewCmd().In(m)
	if chaps {
		m.FFmetaChapsToCue()
	}
	if cover {
		cmd.Args().ACodec("an").Out("cover").Ext(".jpg")
		cmd.Run()
	}
	if meta {
		cmd.Args().Post("f", "ffmetadata").ACodec("none").VCodec("none").Ext(".ini").Out("ffmeta")
		cmd.Run()
	}
}

func(m *Media) Split(cue string) {
	var chaps *Chapters
	switch {
	case cue != "":
		chaps = ReadCueSheet(cue)
	case m.HasChapters():
		chaps = m.Meta.Chapters
	}

	for i, ch := range *chaps {
		cmd := m.Cut(ch.Start, ch.End, i)
		cmd.Run()
	}
}

func(m *Media) Cut(ss, to string, no int) *FFmpegCmd {
	count := fmt.Sprintf("%06d", no + 1)
	cmd := NewCmd().In(m)
	if ss != "" {
		cmd.Args().Post("ss", ss)
	}
	if to != "" {
		cmd.Args().Post("to", to)
	}
	cmd.Args().Out("tmp" + count).Ext(m.Ext)
	return cmd
}

func Join(ext string) *FFmpegCmd {
	ff := NewCmd()
	pre := flagArgs{"f": "concat", "safe": "0"}
	ff.Args().Pre(pre).VCodec("vn").Ext(ext)

	files := find(ext)
	file, err := os.CreateTemp("", "audiofiles")
	if err != nil {
		log.Fatal(err)
	}

	var cat strings.Builder
	for _, f := range files {
		cat.WriteString("file ")
		cat.WriteString("'")
		cat.WriteString(f)
		cat.WriteString("'")
		cat.WriteString("\n")
	}

	if _, err := file.WriteString(cat.String()); err != nil {
		log.Fatal(err)
	}

	ff.tmpFile = file

	ff.In(NewMedia(ff.tmpFile.Name()))
	return ff
}


