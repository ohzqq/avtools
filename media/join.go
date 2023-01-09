package media

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/ff"
)

type JoinCmd struct {
	dir    string
	ext    string
	ffmeta string
}

func Join(ext string, dir ...string) Cmd {
	d := "."
	if len(dir) > 0 {
		d = dir[0]
	}

	path, err := filepath.Abs(d)
	if err != nil {
		log.Fatal(err)
	}

	files, err := filepath.Glob(path + "/*" + ext)
	if err != nil {
		log.Fatal(err)
	}

	tmp, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer tmp.Close()

	var media []*Media
	for _, f := range files {
		media = append(media, New(f))
		str := "file '" + f + "'\n"
		if _, err := tmp.WriteString(str); err != nil {
			log.Fatal(err)
		}
	}

	tmpMedia := media[0]
	tmpMedia.Chapters = GenerateChapters(media)
	s := tmpMedia.SaveMetaFmt("ini")
	s.Run()

	cmd := ff.New()
	cmd.In(tmp.Name())
	cmd.Input.Set("f", "concat")
	cmd.Input.Set("safe", "0")
	cmd.Input.Set("y", "")

	if media[0].HasCover {
		cmd.Output.Set("vn", "")
	}

	base := filepath.Base(d)
	name := filepath.Join(path, base)
	cmd.Output.Set("c", "copy").Ext(ext).Name(name)

	return cmd
}

func GenerateChapters(media []*Media) []*avtools.Chapter {
	var chapters []*avtools.Chapter

	var start = []int64{0}
	var end []int64
	for idx, m := range media {
		d := m.GetTag("duration")
		dur := avtools.ParseStamp(d)
		e := start[idx] + dur.Milliseconds()
		end = append(end, e)
		if idx < len(media)-1 {
			s := end[idx]
			start = append(start, s)
		}
		ss := avtools.ParseStampDuration(start[idx], 1000)
		to := avtools.ParseStampDuration(end[idx], 1000)
		chapter := &avtools.Chapter{
			Start: avtools.Timestamp(ss),
			End:   avtools.Timestamp(to),
			Title: "Chapter " + strconv.Itoa(idx+1),
		}
		chapters = append(chapters, chapter)
	}

	return chapters
}
