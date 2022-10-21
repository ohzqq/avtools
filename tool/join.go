package tool

import (
	"fmt"
	"log"

	"github.com/ohzqq/avtools/chap"
	"github.com/ohzqq/avtools/file"
	"github.com/ohzqq/avtools/media"
)

type JoinCmd struct {
	*Cmd
	ext   string
	dir   string
	files []file.File
}

func Join(ext, dir string) *JoinCmd {
	return &JoinCmd{
		Cmd:   NewCmd(),
		ext:   ext,
		dir:   dir,
		files: FindFilesByExt(ext, dir),
	}
}

func (j *JoinCmd) Parse() *Cmd {
	in := j.MkTmp()
	for _, f := range j.files {
		if _, err := in.WriteString("file '" + f.Abs + "'\n"); err != nil {
			log.Fatal(err)
		}
	}

	ff := j.FFmpeg()

	chaps := j.CalculateChapters()
	med := media.NewMedia(j.files[0].Abs)
	med.SetChapters(chaps)
	med.Meta.SaveAs("ffmeta")

	ff.Input(in.Name())
	ff.VN()
	ff.FFmeta("ffmeta.ini")
	ff.AppendPreInput("f", "concat")
	ff.AppendPreInput("safe", "0")

	if !j.flag.Args.HasOutput() {
		ff.Output("tmp" + j.ext)
	}

	j.Add(ff)
	return j.Cmd
}

func (j JoinCmd) CalculateChapters() chap.Chapters {
	var (
		start = []float64{0}
		end   []float64
		chaps = chap.NewChapters()
	)
	for idx, f := range j.files {
		m := media.NewMedia(f.Abs)
		d := m.Duration().Float()
		e := start[idx] + d
		end = append(end, e)
		if idx < len(j.files)-1 {
			s := end[idx]
			start = append(start, s)
		}
		ch := chap.NewChapter()
		ch.SetTitle(fmt.Sprintf("Chapter %d", idx))
		ch.SetTimebase(1000)
		ss := chap.NewChapterTime(start[idx] * 1000)
		ch.SetStart(ss)
		to := chap.NewChapterTime(end[idx] * 1000)
		ch.SetEnd(to)
		chaps.Chapters = append(chaps.Chapters, ch)
	}
	return chaps
}
