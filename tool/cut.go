package tool

import (
	"github.com/ohzqq/avtools/ffmpeg"
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

func (c *CutCmd) Start(ss string) *CutCmd {
	c.ss = ss
	return c
}

func (c *CutCmd) End(to string) *CutCmd {
	c.to = to
	return c
}

func (c *CutCmd) Parse() *Cmd {
	if c.flag.Args.HasStart() {
		c.ss = c.flag.Args.Start
	}

	if c.flag.Args.HasEnd() {
		c.to = c.flag.Args.End
	}

	ff := c.FFmpegCmd()

	c.Add(ff)

	return c.Cmd
}

func (c *CutCmd) FFmpegCmd() *ffmpeg.Cmd {
	ff := c.FFmpeg()
	ff.SS(c.ss)
	ff.To(c.to)

	return ff
}
