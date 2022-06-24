package avtools

import (
	"os"
	"log"
	"fmt"
	"bytes"
	//"strings"
	"os/exec"
	"strconv"
	//"io/fs"
	//"path/filepath"

	//"cs.opensource.google/go/x/exp/slices"
	//"github.com/alessio/shellescape"
)

type FFmpegCmd struct {
	cmd *exec.Cmd
	input []string
	MediaInput []*Media
	cover string
	ffmeta string
	padding bool
	profile bool
	tmpFile *os.File
	args *CmdArgs
}

func NewFFmpegCmd() *FFmpegCmd {
	ff := FFmpegCmd{
		args: NewArgs(),
		padding: false,
		cmd: exec.Command("ffmpeg", "-hide_banner"),
	}
	if !ff.profile {
		ff.Profile(Cfg().DefaultProfile())
	}
	return &ff
}

func (ff *FFmpegCmd) In(input *Media) *FFmpegCmd {
	ff.MediaInput = append(ff.MediaInput, input)
	return ff
}

func (ff *FFmpegCmd) Profile(p string) *FFmpegCmd {
	if ff.args.PreInput == nil {
		ff.args.Pre(Cfg().GetProfile(p).PreInput)
	}
	if ff.args.PostInput == nil {
		for k, v := range Cfg().GetProfile(p).PostInput {
			ff.args.Post(k, v)
		}
	}
	if ff.args.VideoParams == nil {
		ff.args.VParams(Cfg().GetProfile(p).VideoParams)
	}
	if ff.args.VideoCodec == "" {
		ff.args.VCodec(Cfg().GetProfile(p).VideoCodec)
	}
	if ff.args.VideoFilters == "" {
		ff.args.VFilters(Cfg().GetProfile(p).VideoFilters)
	}
	if ff.args.AudioParams == nil {
		ff.args.AParams(Cfg().GetProfile(p).AudioParams)
	}
	if ff.args.AudioCodec == "" {
		ff.args.ACodec(Cfg().GetProfile(p).AudioCodec)
	}
	if ff.args.AudioFilters == "" {
		ff.args.AFilters(Cfg().GetProfile(p).AudioFilters)
	}
	if ff.args.FilterComplex == "" {
		ff.args.Filter(Cfg().GetProfile(p).FilterComplex)
	}
	if ff.args.Verbosity == "" {
		ff.args.LogLevel(Cfg().Defaults.Verbosity)
	}
	if ff.args.Output == "" {
		ff.args.Out(Cfg().Defaults.Output)
	}
	if ff.args.Overwrite == true {
		ff.args.OverWrite()
	}
	return ff
}

func (ff *FFmpegCmd) Cover(cover string) *FFmpegCmd {
	ff.cover = cover
	return ff
}

func (ff *FFmpegCmd) FFmeta(meta string) *FFmpegCmd {
	ff.ffmeta = meta
	return ff
}

func (ff *FFmpegCmd) Args() *CmdArgs {
	return ff.args
}

func(ff *FFmpegCmd) SetArgs(a *CmdArgs) *FFmpegCmd {
	ff.args = a
	return ff
}

func (ff *FFmpegCmd) Run() {
	if ff.tmpFile != nil {
		defer os.Remove(ff.tmpFile.Name())
	}

	cmd := ff.buildCmd()

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		log.Printf("%v finished with error:\n%v", cmd.String(), err)
		fmt.Printf("%v\n", stderr.String())
	} else {
		fmt.Println("Success!")
	}

	if stdout.String() != "" {
		fmt.Printf("%v\n", stdout.String())
	}
}

func (ff *FFmpegCmd) String() string {
	return ff.buildCmd().String()
}

func (ff *FFmpegCmd) buildCmd() *exec.Cmd {
	var argOrder = []string{
		"Verbosity",
		"Overwrite",
		"Pre",
		"Input",
		"Meta",
		"Post",
		"MiscParams",
		"VideoCodec",
		"VideoParams",
		"VideoFilters",
		"FilterComplex",
		"AudioCodec",
		"AudioParams",
		"AudioFilters",
		"Output",
		"Ext",
	}

	for _, arg := range argOrder {
		switch arg {
		case "Verbosity":
			if Cfg().Defaults.Verbosity != "" {
				ff.push("-loglevel")
				ff.push(Cfg().Defaults.Verbosity)
			}
		case "Overwrite":
			if Cfg().Defaults.Overwrite {
				ff.push("-y")
			}

			if ff.args.Overwrite {
				ff.push("-y")
			}
		case "Pre":
			for _, arg := range ff.args.PreInput.Split() {
				ff.push(arg)
			}
		case "Input":
			ff.processInput()
			ff.mapInput()
		case "Post":
			for _, arg := range ff.args.PostInput.Split() {
				ff.push(arg)
			}
		case "VideoCodec":
			switch vc := ff.args.VideoCodec; vc {
			case "none":
			case "":
			case "vn":
				ff.push("-vn")
			default:
				ff.push("-c:v")
				ff.push(ff.args.VideoCodec)
			}
		case "VideoParams":
			for _, arg := range ff.args.VideoParams.Split() {
				ff.push(arg)
			}
		case "VideoFilters":
			if ff.args.VideoFilters != "" {
				ff.push("-vf")
				ff.push(ff.args.VideoFilters)
			}
		case "FilterComplex":
			if ff.args.FilterComplex != "" {
				ff.push("-filter")
				ff.push(ff.args.FilterComplex)
			}
		case "MiscParams":
			if params := ff.args.MiscParams; len(params) > 0 {
				for _, p := range params {
					ff.push(p)
				}
			}
		case "AudioCodec":
			switch ac := ff.args.AudioCodec; ac {
			case "none":
			case "":
			case "an":
				ff.push("-an")
			default:
				ff.push("-c:a")
				ff.push(ff.args.AudioCodec)
			}
		case "AudioParams":
			for _, arg := range ff.args.AudioParams.Split() {
				ff.push(arg)
			}
		case "AudioFilters":
			if ff.args.AudioFilters != "" {
				ff.push("-af")
				ff.push(ff.args.AudioFilters)
			}
		case "Output":
			ff.Output()
		}
	}
	return ff.cmd
}

func (ff *FFmpegCmd) push(arg string) {
	ff.cmd.Args = append(ff.cmd.Args, arg)
}

func (ff *FFmpegCmd) processInput() {
	if len(ff.MediaInput) > 0 {
		for _, i := range ff.MediaInput {
			ff.pushInput(i.Path)
		}
	} else {
		log.Fatal("No input specified")
	}

	if ff.args.AlbumArt != "" {
		ff.pushInput(ff.args.AlbumArt)
	}

	if ff.args.Metadata != "" {
		ff.pushInput(ff.args.Metadata)
	}
}

func (ff *FFmpegCmd) pushInput(input string) {
	ff.push("-i")
	ff.push(input)
}

func (ff *FFmpegCmd) mapInput() {
	if ff.args.AlbumArt != "" || ff.args.Metadata != "" {
		for idx, _ := range ff.MediaInput {
			ff.push("-map")
			ff.push(strconv.Itoa(idx) + ":0")
		}
	}

	idx := len(ff.MediaInput)
	if ff.args.AlbumArt != "" {
		ff.push("-map")
		ff.push(strconv.Itoa(idx) + ":0")
		idx++
	}

	if ff.args.Metadata != "" {
		ff.push("-map_metadata")
		ff.push(strconv.Itoa(idx))
		idx++
	}
}

//func (ff *FFmpegCmd) metadata(meta string) {
//  ff.push("-i")
//  ff.push(meta)
//}

func (ff *FFmpegCmd) Output() {
	var o string
	var pad string
	if Cfg().Defaults.Output != "" {
		o = Cfg().Defaults.Output
		if ff.padding {
			pad = "%06d"
		} else {
			pad = ""
		}
	}

	var ext string
	if ff.args.Extension != "" {
		ext = ff.args.Extension
	} else {
		ext = ff.MediaInput[0].Ext
	}

	if ff.args.Output == "" {
		ff.push(o + pad + ext)
	} else {
		ff.push(ff.args.Output + pad + ext)
	}
}

