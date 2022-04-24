package fftools

import (
	"fmt"
	//"io/fs"
	"os"
	"os/exec"
	"log"
	"bytes"
	"strings"
	"path/filepath"
)

func AddAlbumArt(m *Media, cover string) *FFmpegCmd {
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

func RmAlbumArt(m *Media) *FFmpegCmd {
	cmd := NewCmd().In(m)
	cmd.Args().Out("tmp-nocover").VCodec("vn")
	return cmd
}

func AddFFmeta(m *Media, meta string) *FFmpegCmd {
	path, err := filepath.Abs(meta)
	if err != nil {
		log.Fatal(err)
	}
	cmd := NewCmd().In(m)
	cmd.Args().Meta(path)
	return cmd
}

func RmFFmeta(m *Media) *FFmpegCmd {
	arg := newFlagArg("map_metadata", "-1")
	cmd := NewCmd().In(m)
	cmd.Args().Post(arg).Out("no-meta")
	return cmd
}

func AddChapters(m *Media) {
}

func RmChapters(m *Media) *FFmpegCmd {
	arg := newFlagArg("map_chapters", "-1")
	cmd := NewCmd().In(m)
	cmd.Args().Post(arg).Out("no-chaps")
	return cmd
}

func ConvertFFmetaChapToCue(m *Media) {
}

func ConvertCueToFFmetaChap(m *Media) {
}

func Split(m *Media) {
	if m.HasChapters() {
		for i, chap := range *m.Meta.Chapters {
			cmd := Cut(m, chap.Start, chap.End, i)
			cmd.Run()
		}
	}
}

func Cut(m *Media, ss, to string, no int) *FFmpegCmd {
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
	return cmd
}

func Join(ext string) *FFmpegCmd {
	ff := NewCmd()
	pre := flagArgs{"f": "concat", "safe": "0"}
	ff.Args().Pre(pre).Ext(ext)

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


