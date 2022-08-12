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
	//fmt.Printf("%+V\n", cmd.media.ListFormats())
	fmt.Printf("%+V\n", cmd.media.Meta())
	//fmt.Printf("%+V\n", cmd.media.GetFormat("audio"))
}

func (c *ffmpegCmd) getChapters() (Chapters, error) {
	//if len(c.media.json) == 0 {
	//  c.media.JsonMeta().Unmarshal()
	//}

	switch {
	case c.opts.CueFile != "":
		return LoadCueSheet(c.opts.CueFile).Chapters, nil
	case c.opts.MetaFile != "":
		return LoadFFmetadataIni(c.opts.MetaFile).Chapters, nil
	//case c.media.HasChapters():
	//return c.media.Meta().Chapters, nil
	default:
		return nil, fmt.Errorf("There are no chapters!")
	}
}

func (c *ffmpegCmd) Extract() {
	//c.media.JsonMeta().Unmarshal()
	c.ParseOptions()

	switch {
	case c.opts.CueSwitch:
		c.media.FFmetaChapsToCue()
		return
	case c.opts.CoverSwitch:
		c.AudioCodec = "an"
		c.VideoCodec = "copy"
		c.Output = "cover"
		c.Ext = ".jpg"
	case c.opts.MetaSwitch:
		c.AppendMapArg("post", "f", "ffmetadata")
		c.AudioCodec = "none"
		c.VideoCodec = "none"
		c.Output = "ffmeta"
		c.Ext = ".ini"
	}

	cmd := c.ParseArgs()
	cmd.Run()
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
	//c.media.JsonMeta().Unmarshal()
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
	chaps, err := cmd.getChapters()
	if err != nil {
		return err
	}

	m := cmd.media.Input
	for i, ch := range chaps {
		NewFFmpegCmd(m.Path).Options(cmd.opts).Cut(ch.StartToSeconds(), ch.EndToSeconds(), i)
	}
	return nil
}

func (cmd *ffmpegCmd) Cut(ss, to string, no int) {
	//cmd.media.JsonMeta().Unmarshal()
	cmd.ParseOptions()

	var (
		count = no + 1
		start = ss
		end   = to
	)

	if cmd.opts.ChapNo != 0 {
		chaps, err := cmd.getChapters()
		if err != nil {
			log.Fatal(err)
		}
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
