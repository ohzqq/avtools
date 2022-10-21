package tool

import (
	"log"

	"github.com/ohzqq/avtools/ffmpeg"
	"github.com/ohzqq/avtools/file"
)

type CutCmd struct {
	*Cmd
	ss string
	to string
}

func Cut() *CutCmd {
	return &CutCmd{
		Cmd: NewCmd(),
	}
}

func (c CutCmd) Chap(no int) *ffmpeg.Cmd {
	var (
		num   = no - 1
		chaps = c.Media().Meta.Chapters.Chapters
		ff    *ffmpeg.Cmd
	)

	if num < len(chaps) {
		ch := chaps[num]
		ss := ch.Start().SecsString()
		to := ch.End().SecsString()
		in := file.New(c.Media().Input.String())
		c.Start(ss).End(to)
		ff = c.FFmpegCmd()
		ff.Output(in.Pad(no))
	} else {
		log.Fatalf("there are only %v chapters", len(chaps))
	}

	return ff
}

func (c *CutCmd) Start(ss string) *CutCmd {
	c.ss = ss
	return c
}

func (c *CutCmd) End(to string) *CutCmd {
	c.to = to
	return c
}

func (c *CutCmd) Parse() *Cmd {
	var ff *ffmpeg.Cmd

	if c.flag.Args.HasChapNo() {
		ff = c.Chap(c.Args.ChapNo)
	} else {
		if c.flag.Args.HasStart() {
			c.ss = c.Args.Start
		}

		if c.flag.Args.HasEnd() {
			c.to = c.Args.End
		}

		o := c.Args.Input.AddSuffix(c.ss + "-" + c.to)

		ff = c.FFmpegCmd()
		ff.Output(o)
	}

	c.Add(ff)

	return c.Cmd
}

func (c *CutCmd) FFmpegCmd() *ffmpeg.Cmd {
	ff := c.FFmpeg()
	ff.SS(c.ss)
	ff.To(c.to)

	return ff
}
