package avtools

import (
	"fmt"
	"log"
	"os"
)

type ffmpegCmd struct {
	media *Media
	args  cmdArgs
	opts  *Options
	*Args
}

func NewFFmpegCmd(i string) *ffmpegCmd {
	ff := ffmpegCmd{}

	if i != "" {
		ff.media = NewMedia(i)
	}
	return &ff
}

func (cmd *ffmpegCmd) Options(f *Options) *ffmpegCmd {
	cmd.opts = f
	return cmd
}

func (cmd *ffmpegCmd) ShowMeta() {
	cmd.ParseOptions()
	args := cmd.Args
	args.Input = cmd.media.GetFile("input").Path()
	args.Name = cmd.media.SafeName()
	args.Options = *cmd.opts
	fmt.Printf("%+V\n", args.Parse())
	if cmd.opts.CueSwitch {
		cmd.media.Meta().MarshalTo("cue").Render().Print()
	}
	if cmd.opts.JsonSwitch {
		cmd.media.Meta().MarshalTo("json").Render().Print()
	}
	if cmd.opts.MetaSwitch {
		cmd.media.Meta().MarshalTo("ffmeta").Render().Print()
	}
}

func (cmd *ffmpegCmd) Extract() {
	cmd.ParseOptions()

	switch {
	case cmd.opts.CueSwitch:
		cmd.media.Meta().MarshalTo("cue").Render().Write()
		return
	case cmd.opts.CoverSwitch:
		cmd.AudioCodec = "an"
		cmd.VideoCodec = "copy"
		cmd.Ext = ".jpg"
	case cmd.opts.MetaSwitch:
		cmd.AppendMapArg("post", "f", "ffmetadata")
		cmd.AudioCodec = "none"
		cmd.VideoCodec = "none"
		cmd.Ext = ".ini"
	}

	cmd.Output = cmd.media.SafeName()
	cmd.Padding = ""

	command := cmd.ParseArgs()
	command.Run()
}

func (cmd *ffmpegCmd) Join(ext string) {
	cmd.ParseOptions()

	tmp, err := os.CreateTemp("", "audiofiles")
	if err != nil {
		log.Fatal(err)
	}

	files := find(ext)
	for _, f := range files {
		if _, err := tmp.WriteString("file '" + f + "'\n"); err != nil {
			log.Fatal(err)
		}
	}

	cmd.AppendMapArg("pre", "f", "concat")
	cmd.AppendMapArg("pre", "safe", "0")
	cmd.Input = tmp.Name()
	cmd.VideoCodec = "vn"
	cmd.Ext = ext

	c := cmd.ParseArgs()
	c.tmpFile = tmp

	c.Run()
}

func (c *ffmpegCmd) Remove() {
	c.ParseOptions()

	if c.opts.ChapSwitch {
		c.AppendMapArg("post", "map_chapters", "-1")
	}

	if c.opts.CoverSwitch {
		c.VideoCodec = "vn"
	}

	if c.opts.MetaSwitch {
		c.AppendMapArg("post", "map_metadata", "-1")
	}

	cmd := c.ParseArgs()
	cmd.Run()
}

func (cmd *ffmpegCmd) Split() error {
	chaps := cmd.media.Meta().Chapters

	m := cmd.media.GetFile("input")
	for i, ch := range chaps {
		NewFFmpegCmd(m.Path()).Options(cmd.opts).Cut(ch.StartToSeconds(), ch.EndToSeconds(), i)
	}
	return nil
}

func (cmd *ffmpegCmd) Cut(ss, to string, no int) {
	cmd.ParseOptions()

	var (
		count = no + 1
		start = ss
		end   = to
	)

	if cmd.opts.ChapNo != 0 {
		chaps := cmd.media.Meta().Chapters
		ch := chaps[cmd.opts.ChapNo-1]
		count = cmd.opts.ChapNo
		start = ch.StartToSeconds()
		end = ch.EndToSeconds()
	}

	cmd.PreInput = mapArgs{}
	cmd.num = count

	if start != "" {
		cmd.AppendMapArg("pre", "ss", start)
	}

	if end != "" {
		cmd.AppendMapArg("pre", "to", end)
	}

	c := cmd.ParseArgs()
	c.Run()
}
