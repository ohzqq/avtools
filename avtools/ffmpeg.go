package avtools

import (
	"os"
	"log"
	"fmt"
	"bytes"
	"strings"
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
	profile bool
	tmpFile *os.File
	args *CmdArgs
}

func NewFFmpegCmd() *FFmpegCmd {
	ff := FFmpegCmd{
		args: NewArgs(),
		cmd: exec.Command("ffmpeg", "-hide_banner"),
	}
	return &ff
}

func (ff *FFmpegCmd) In(input *Media) *FFmpegCmd {
	ff.MediaInput = append(ff.MediaInput, input)
	return ff
}

func (ff *FFmpegCmd) Profile(p string) *FFmpegCmd {
	ff.profile = true
	profile := Cfg().GetProfile(p)

	if ff.args.PreInput == nil {
		for _, arg := range profile.PreInput {
			ff.args.PreInput = append(ff.args.PreInput, arg)
		}
	}

	if ff.args.PostInput == nil {
		for _, arg := range profile.PostInput {
			ff.args.PostInput = append(ff.args.PostInput, arg)
		}
	}

	if ff.args.VideoParams == nil {
		for _, arg := range profile.VideoParams {
			ff.args.VideoParams = append(ff.args.VideoParams, arg)
		}
	}

	if ff.args.VideoCodec == "" {
		ff.args.VCodec(profile.VideoCodec)
	}

	if ff.args.VideoFilters == "" {
		ff.args.VFilters(profile.VideoFilters)
	}

	if ff.args.AudioParams == nil {
		for _, arg := range profile.AudioParams {
			ff.args.AudioParams = append(ff.args.AudioParams, arg)
		}
	}

	if ff.args.AudioCodec == "" {
		ff.args.ACodec(profile.AudioCodec)
	}

	if ff.args.AudioFilters == "" {
		ff.args.AFilters(profile.AudioFilters)
	}

	if ff.args.Padding == "" {
		ff.args.Pad(Cfg().defaults.Padding)
	}

	if len(ff.args.FilterComplex) == 0 {
		ff.args.Filters(profile.FilterComplex)
	}

	if ff.args.Verbosity == "" {
		ff.args.LogLevel(Cfg().defaults.Verbosity)
	}

	if profile.Extension != "" {
		ff.args.Ext(profile.Extension)
	}

	if ff.args.Output == "" {
		ff.args.Out(Cfg().defaults.Output)
	}

	if ff.args.Overwrite == true {
		ff.args.OverWrite()
	}

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

		fmt.Println(cmd.String())
	if ff.Args().verbose {
		//fmt.Printf("%+v\n", ff.Args())
		fmt.Println(cmd.String())
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
			if Cfg().defaults.Verbosity != "" {
				ff.push("-loglevel")
				ff.push(Cfg().defaults.Verbosity)
			}
		case "Overwrite":
			if Cfg().defaults.Overwrite {
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
			case "":
			case "none":
				fallthrough
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
			if len(ff.args.FilterComplex) > 0 {
				ff.push("-vf")
				ff.push(strings.Join(ff.args.FilterComplex, ","))
			}
		case "MiscParams":
			if params := ff.args.MiscParams; len(params) > 0 {
				for _, p := range params {
					ff.push(p)
				}
			}
		case "AudioCodec":
			switch ac := ff.args.AudioCodec; ac {
			case "":
			case "none":
				fallthrough
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

func (ff *FFmpegCmd) Output() {
	var (
		name string
		ext string
	)

	if out := ff.args.Output; out != "" {
		name = out
	}

	if p := ff.args.Padding; p != "" {
		name = name + fmt.Sprintf(p, ff.args.num)
	}

	switch e := ff.args.Extension; e != "" {
	case true:
		ext = e
	default:
		ext = ff.MediaInput[0].Ext
	}

	ff.push(name + ext)
}

