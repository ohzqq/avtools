package fftools

import (
	"fmt"
	//"io/fs"
	"os"
	"log"
	"strings"
)

func AddAlbumArt(m *Media, cover string) {
	switch codec := m.AudioCodec(); codec {
	case "aac":
	case "mp3":
		cmd := NewCmd().In(m).Cover(cover)
		cmd.Args().Out("tmp-cover").Params(Mp3CoverArgs())
		fmt.Printf("%v", cmd.String())
	}
}

func RmAlbumArt(m *Media) {
}

func AddFFmeta(m *Media) {
}

func RmFFmeta(m *Media) {
}

func AddChapters(m *Media) {
}

func RmChapters(m *Media) {
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


