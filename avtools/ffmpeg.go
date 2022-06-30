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
	exec *exec.Cmd
	flags *Flags
	*Args
}

func NewFFmpegCmd(i string) *ffmpegCmd {
	return &ffmpegCmd{media: NewMedia(i)}
}

func(c *ffmpegCmd) Extract() {
	switch {
	case c.flags.CueSwitch:
		c.media.FFmetaChapsToCue()
		return
	case c.flags.CoverSwitch:
		c.AudioCodec = "an"
		c.VideoCodec = "copy"
		c.Output = "cover"
		c.Ext = ".jpg"
	case c.flags.MetaSwitch:
		c.PostInput = append(c.PostInput, newMapArg("f", "ffmetadata"))
		c.AudioCodec = "none"
		c.VideoCodec = "none"
		c.Output = "ffmeta"
		c.Ext = ".ini"
	}
	cmd := c.Parse()
	cmd.Run()
}

func(cmd *ffmpegCmd) ParseFlags() *ffmpegCmd {
	cmd.media.JsonMeta().Unmarshal()

	if meta := cmd.flags.MetaFile; meta != "" {
		cmd.media.SetMeta(LoadFFmetadataIni(meta))
	}

	if cue := cmd.flags.CueFile; cue != "" {
		cmd.media.SetChapters(LoadCueSheet(cue))
	}

	if y := cmd.flags.Overwrite; y {
		cmd.Overwrite = y
	}

	if o := cmd.flags.Output; o != "" {
		cmd.Name = o
	}

	if c := cmd.flags.ChapNo; c  != 0 {
		cmd.num = c
	}

	return cmd
}

func(cmd *ffmpegCmd) Parse() Cmd {
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
	cover := cmd.Flags.CoverFile
	meta := cmd.Flags.MetaFile
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

	return Cmd{exec: exec.Command("ffmpeg", cmd.args.args...), flags: cmd.flags}
}
