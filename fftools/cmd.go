package media

import (
	//"os"
	"os/exec"
)

type FFmpegCmd struct {
	exec.Cmd
	Input string
	VF string
	AF string
	CV string
	CA string
	VParams string
	CParams string
}

func FFcmd() *FFmpegCmd (
	return exec.Command("ffmpeg")
)

func (ff *FFmpegCmd) FFArgs(args ...string) {
	ff.Args = args
}

func (ff *FFmpegCmd) Input(input ...string) string {
}

func (ff *FFmpegCmd) VF(filters ...string) string {
}

func (ff *FFmpegCmd) AF(filters ...string) string {
}

func (ff *FFmpegCmd) FilterComplex(filter string) string {
}

func (ff *FFmpegCmd) CV(codec string) string {
}

func (ff *FFmpegCmd) CA(codec string) string {
}

func (ff *FFmpegCmd) VParams(params ...string) string {
}

func (ff *FFmpegCmd) AParams(params ...string) string {
}

