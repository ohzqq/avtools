package fftools

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
	profile bool
	padding bool
	tmpFile *os.File
	args CmdArgs
}

func NewCmd() *FFmpegCmd {
	ff := FFmpegCmd{}
	ff.padding = false
	ff.cmd = exec.Command("ffmpeg", "-hide_banner")
	return &ff
}

func (ff *FFmpegCmd) In(input *Media) *FFmpegCmd {
	ff.MediaInput = append(ff.MediaInput, input)
	return ff
}

func (ff *FFmpegCmd) Profile(p string) *FFmpegCmd {
	ff.profile = true
	if ff.args.PreInput == nil {
		ff.args.Pre(Cfg.Profiles[p].PreInput)
	}
	if ff.args.PostInput == nil {
		ff.args.Post(Cfg.Profiles[p].PostInput)
	}
	if ff.args.VideoParams == nil {
		ff.args.VParams(Cfg.Profiles[p].VideoParams)
	}
	if ff.args.VideoCodec == "" {
		ff.args.VCodec(Cfg.Profiles[p].VideoCodec)
	}
	if ff.args.VideoFilters == "" {
		ff.args.VFilters(Cfg.Profiles[p].VideoFilters)
	}
	if ff.args.AudioParams == nil {
		ff.args.AParams(Cfg.Profiles[p].AudioParams)
	}
	if ff.args.AudioCodec == "" {
		ff.args.ACodec(Cfg.Profiles[p].AudioCodec)
	}
	if ff.args.AudioFilters == "" {
		ff.args.AFilters(Cfg.Profiles[p].AudioFilters)
	}
	if ff.args.FilterComplex == "" {
		ff.args.Filter(Cfg.Profiles[p].FilterComplex)
	}
	if ff.args.Verbosity == "" {
		ff.args.LogLevel(Cfg.Defaults.Verbosity)
	}
	if ff.args.Output == "" {
		ff.args.Out(Cfg.Defaults.Output)
	}
	if ff.args.Overwrite == false {
		ff.args.OverWrite(Cfg.Defaults.Overwrite)
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
	if !ff.profile {
		ff.args = CmdArgs{}
		ff.args.VCodec("copy")
		ff.args.ACodec("copy")
	}
	return &ff.args
}

func (ff *FFmpegCmd) Run() {
	if ff.tmpFile != nil {
		defer os.Remove(ff.tmpFile.Name())
	}

	cmd := ff.Cmd()

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	fmt.Printf("%v\n", cmd.String())

	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		fmt.Printf("%v\n", stderr.String())
	}
	fmt.Printf("%v\n", stdout.String())
}

func (ff *FFmpegCmd) String() string {
	return ff.Cmd().String()
}

func (ff *FFmpegCmd) Cmd() *exec.Cmd {
	for _, arg := range argOrder {
		switch arg {
		case "Verbosity":
			ff.Verbosity()
		case "Overwrite":
			ff.Overwrite()
		case "Pre":
			ff.Pre()
		case "Input":
			ff.processInput()
			ff.mapInput()
		case "Post":
			ff.Post()
		case "VideoCodec":
			ff.VideoCodec()
		case "VideoParams":
			ff.VideoParams()
		case "VideoFilters":
			ff.VideoFilters()
		case "FilterComplex":
			ff.FilterComplex()
		case "MiscParams":
			ff.MiscParams()
		case "AudioCodec":
			ff.AudioCodec()
		case "AudioParams":
			ff.AudioParams()
		case "AudioFilters":
			ff.AudioFilters()
		case "Output":
			ff.Output()
		}
	}
	return ff.cmd
}

func (ff *FFmpegCmd) push(arg string) {
	ff.cmd.Args = append(ff.cmd.Args, arg)
}

func (ff *FFmpegCmd) Verbosity() {
	if Cfg.Defaults.Verbosity != "" {
		ff.push("-loglevel")
		ff.push(Cfg.Defaults.Verbosity)
	}
}

func (ff *FFmpegCmd) processInput() {
	if len(ff.MediaInput) > 0 {
		for _, i := range ff.MediaInput {
			ff.pushInput(i.Path)
		}
	} else {
		log.Fatal("No input specified")
	}

	if ff.cover != "" {
		ff.pushInput(ff.cover)
	}

	if ff.ffmeta != "" {
		ff.pushInput(ff.ffmeta)
	}
}

func (ff *FFmpegCmd) pushInput(input string) {
	ff.push("-i")
	ff.push(input)
}

func (ff *FFmpegCmd) mapInput() {
	if ff.cover != "" || ff.ffmeta != "" {
		for idx, _ := range ff.MediaInput {
			ff.push("-map " + strconv.Itoa(idx) + ":0")
		}
	}

	idx := len(ff.MediaInput)
	if ff.cover != "" {
		ff.push("-map " + strconv.Itoa(idx) + ":0")
		idx++
	}

	if ff.ffmeta != "" {
		ff.push("-map_metadata " + strconv.Itoa(idx))
		idx++
	}
}

func (ff *FFmpegCmd) metadata(meta string) {
	ff.push("-i")
	ff.push(meta)
}
func (ff *FFmpegCmd) Pre() {
	for _, arg := range ff.args.PreInput.Split() {
		ff.push(arg)
	}
}

func (ff *FFmpegCmd) Overwrite() {
	if Cfg.Defaults.Overwrite {
		ff.push("-y")
	}
}

func (ff *FFmpegCmd) Post() {
	for _, arg := range ff.args.PostInput.Split() {
		ff.push(arg)
	}
}

func (ff *FFmpegCmd) VideoCodec() {
	switch vc := ff.args.VideoCodec; vc {
	case "none":
	case "":
	case "vn":
		ff.push("-vn")
	default:
		ff.push("-c:v")
		ff.push(ff.args.VideoCodec)
	}
}

func (ff *FFmpegCmd) VideoParams() {
	for _, arg := range ff.args.VideoParams.Split() {
		ff.push(arg)
	}
}

func (ff *FFmpegCmd) VideoFilters() {
	if ff.args.VideoFilters != "" {
		ff.push("-vf")
		ff.push(ff.args.VideoFilters)
	}
}

func (ff *FFmpegCmd) FilterComplex() {
	if ff.args.FilterComplex != "" {
		ff.push("-filter")
		ff.push(ff.args.FilterComplex)
	}
}

func (ff *FFmpegCmd) MiscParams() {
	if params := ff.args.MiscParams; len(params) > 0 {
		for _, p := range params {
			ff.push(p)
		}
	}
}

func (ff *FFmpegCmd) AudioCodec() {
	switch ac := ff.args.AudioCodec; ac {
	case "none":
	case "":
	case "an":
		ff.push("-an")
	default:
		ff.push("-c:a")
		ff.push(ff.args.AudioCodec)
	}
}

func (ff *FFmpegCmd) AudioParams() {
	for _, arg := range ff.args.AudioParams.Split() {
		ff.push(arg)
	}
}

func (ff *FFmpegCmd) AudioFilters() {
	if ff.args.AudioFilters != "" {
		ff.push("-af")
		ff.push(ff.args.AudioFilters)
	}
}

func (ff *FFmpegCmd) Output() {
	var o string
	var pad string
	if Cfg.Defaults.Output != "" {
		o = Cfg.Defaults.Output
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
