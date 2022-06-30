package avtools

import (
	"fmt"
	//"os"
	"os/exec"
	//"bytes"
	//"log"
	"strconv"
	//"strings"
	//"path/filepath"
)
var _ = fmt.Printf

type ffmpegCmd struct {
	media *Media
	args cmdArgs
	opts *Options
	*Args
}

func NewFFmpegCmd(i string) *ffmpegCmd {
	ff := &ffmpegCmd{media: NewMedia(i)}
	return ff
}

func(cmd *ffmpegCmd) Options(f *Options) *ffmpegCmd {
	cmd.opts = f
	return cmd
}

func(c *ffmpegCmd) getChapters() ([]*Chapter, error) {
	if len(c.media.json) == 0 {
		c.media.JsonMeta().Unmarshal()
	}

	switch {
	case c.opts.CueFile != "":
		return LoadCueSheet(c.opts.CueFile), nil
	case c.opts.MetaFile != "":
		return LoadFFmetadataIni(c.opts.MetaFile).Chapters, nil
	case c.media.HasChapters():
		return c.media.Meta.Chapters, nil
	default:
		return nil, fmt.Errorf("There are no chapters!")
	}
}

func(c *ffmpegCmd) Extract() {
	c.media.JsonMeta().Unmarshal()
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
	cmd := c.Parse()
	cmd.Run()
}

func(c *ffmpegCmd) Remove() {
	c.media.JsonMeta().Unmarshal()
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

	cmd := c.Parse()
	cmd.Run()
}

//func(c *Cmd) Split() error {
//  chaps, err := c.getChapters()
//  if err != nil {
//    return err
//  }

//  for i, ch := range chaps {
//    cmd := c.Cut(ch.StartToSeconds(), ch.EndToSeconds(), i)
//    if c.Overwrite {
//      cmd.Args().OverWrite()
//    }
//    cmd.Run()
//  }
//  return nil
//}

//func(c *Cmd) Cut(ss, to string, no int) *FFmpegCmd {
//  var (
//    count = no + 1
//    start = ss
//    end = to
//  )

//  if c.Arg.ChapNo != 0 {
//    chaps, err := c.getChapters()
//    if err != nil {
//      log.Fatal(err)
//    }
//    ch := chaps[c.Arg.ChapNo - 1]
//    count = c.Arg.ChapNo
//    start = ch.StartToSeconds()
//    end = ch.EndToSeconds()
//  }

//  c.Args().
//    Pre("ss", start).
//    Pre("to", end).
//    Num(count).
//    Out("tmp")

//  cmd := c.Parse()
//  cmd.Run()
//}

func(cmd *ffmpegCmd) ParseOptions() *ffmpegCmd {
	cmd.Args = Cfg().GetProfile(cmd.opts.Profile)

	if meta := cmd.opts.MetaFile; meta != "" {
		cmd.media.SetMeta(LoadFFmetadataIni(meta))
	}

	if cue := cmd.opts.CueFile; cue != "" {
		cmd.media.SetChapters(LoadCueSheet(cue))
	}

	if y := cmd.opts.Overwrite; y {
		cmd.Overwrite = y
	}

	if o := cmd.opts.Output; o != "" {
		cmd.Name = o
	}

	if c := cmd.opts.ChapNo; c  != 0 {
		cmd.num = c
	}

	return cmd
}

func(cmd *ffmpegCmd) Parse() *Cmd {
	if log := cmd.LogLevel; log != "" {
		cmd.args.Append("-v", log)
	}

	if cmd.Overwrite {
		cmd.args.Append("-y")
	}

	// pre input
	if pre := cmd.PreInput; len(pre) > 0 {
		cmd.args.Append(pre.Split()...)
	}

	// input
	cmd.args.Append("-i", cmd.media.Path)
	cover := cmd.opts.CoverFile
	meta := cmd.opts.MetaFile
	if meta != "" {
		cmd.args.Append("-i", meta)
	}

	if cover != "" {
		cmd.args.Append("-i", cover)
	}

	//map input
	idx := 0
	if cover != "" || meta != "" {
		cmd.args.Append("-map", strconv.Itoa(idx) + ":0")
		idx++
	}

	if cover != "" {
		cmd.args.Append("-map", "0:" + strconv.Itoa(idx))
		idx++
	}

	if meta != "" {
		cmd.args.Append("-map_metadata", strconv.Itoa(idx))
		idx++
	}

	// post input
	if post := cmd.PostInput; len(post) > 0 {
		cmd.args.Append(post.Split()...)
	}

	//video codec
	if codec := cmd.VideoCodec; codec != "" {
		switch codec {
		case "":
		case "none", "vn":
			cmd.args.Append("-vn")
		default:
			cmd.args.Append("-c:v", codec)
			//video params
			if params := cmd.VideoParams.Split(); len(params) > 0 {
				cmd.args.Append(params...)
			}

			//video filters
			if filters := cmd.VideoFilters.Join(); len(filters) > 0 {
				cmd.args.Append("-vf", filters)
			}
		}
	}

	//filter complex
	if filters := cmd.FilterComplex.Join(); len(filters) > 0 {
		cmd.args.Append("-vf", filters)
	}

	//audio codec
	if codec := cmd.AudioCodec; codec != "" {
		switch codec {
		case "":
		case "none", "an":
			cmd.args.Append("-an")
		default:
			cmd.args.Append("-c:a", codec)
			//audio params
			if params := cmd.AudioParams.Split(); len(params) > 0 {
				cmd.args.Append(params...)
			}

			//audio filters
			if filters := cmd.AudioFilters.Join(); len(filters) > 0 {
				cmd.args.Append("-af", filters)
			}
		}
	}

	//output
	var (
		name string
		ext string
	)

	if out := cmd.Output; out != "" {
		name = out
	}

	if p := cmd.Padding; p != "" {
		name = name + fmt.Sprintf(p, cmd.num)
	}

	switch {
	//case cmd.Ext != "":
	//  ext = cmd.Ext
	case cmd.Ext != "":
		ext = cmd.Ext
	default:
		ext = cmd.media.Ext
	}
	cmd.args.Append(name + ext)

	return NewCmd(exec.Command("ffmpeg", cmd.args.args...), cmd.opts.Verbose)
}
