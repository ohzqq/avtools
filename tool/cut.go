package tool

import (
	"log"
	"strings"

	"github.com/ohzqq/avtools/ffmpeg"
	"github.com/ohzqq/avtools/file"
)

type CutCmd struct {
	*Cmd
}

func Cut() *CutCmd {
	return &CutCmd{
		Cmd: NewCmd(),
	}
}

func (c CutCmd) Chap(no int) *ffmpeg.Cmd {
	var (
		num   = no - 1
		chaps = c.Media.Meta.Chapters.Chapters
		ff    *ffmpeg.Cmd
	)

	if num < len(chaps) {
		ch := chaps[num]
		ss := ch.Start().SecsString()
		to := ch.End().SecsString()
		in := file.New(c.Media.Input.String())
		c.Start = ss
		c.End = to
		ff = c.FFmpegCmd()
		ff.Output(in.Pad(no))
	} else {
		log.Fatalf("there are only %v chapters", len(chaps))
	}

	return ff
}

func (c *CutCmd) FFmpegCmd() *ffmpeg.Cmd {
	ff := c.FFmpeg()
	ff.SS(c.Start)
	ff.To(c.End)

	return ff
}

func (c *CutCmd) Parse() *Cmd {
	var ff *ffmpeg.Cmd

	if c.HasChapNo() {
		ff = c.Chap(c.ChapNo)
	} else {
		var (
			ss string
			to string
		)
		if c.HasStart() {
			ss = strings.ReplaceAll(c.Start, ":", "")
		}

		if c.HasEnd() {
			to = strings.ReplaceAll(c.End, ":", "")
		}

		o := c.Input.AddSuffix(ss + "-" + to)

		ff = c.FFmpegCmd()
		ff.Output(o)
	}

	c.Add(ff)

	return c.Cmd
}
