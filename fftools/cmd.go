package fftools

import (
	//"os"
	"strings"
	"os/exec"
)

type FFmpegCmd struct {
	Cmd *exec.Cmd
	Input string
	VideoFilter string
	AudioFilter string
	CodecVideo string
	CodecAudio string
	VideoParams string
	AudioParams string
}

func New() *FFmpegCmd {
	ff := FFmpegCmd{}
	ff.Cmd = exec.Command("ffmpeg", "-hide_banner")
	return &ff
}

func (ff *FFmpegCmd) Args(args ...string) {
	ff.Cmd.Args = append(cmd.Args, args...)
}

func (ff *FFmpegCmd) Files(input ...string) string {
	return strings.Join(input, " ")
}

func (ff *FFmpegCmd) VF(filters ...string) string {
	return strings.Join(filters , " ")
}

func (ff *FFmpegCmd) AF(filters ...string) string {
	return strings.Join(filters , " ")
}

func (ff *FFmpegCmd) FilterComplex(filter string) string {
	return filter
}

func (ff *FFmpegCmd) CV(codec string) string {
	return codec
}

func (ff *FFmpegCmd) CA(codec string) string {
	return codec
}

func (ff *FFmpegCmd) VParams(params ...string) string {
	return strings.Join(params, " ")
}

func (ff *FFmpegCmd) AParams(params ...string) string {
	return strings.Join(params, " ")
}

