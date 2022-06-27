package avtools

import (
	"fmt"
	"os"
	//"bytes"
	//"os/exec"
	"log"
	//"strings"
	//"path/filepath"
)
var _ = fmt.Printf

type Cmd struct {
	//Media *Media
	args *Args
	CliArgs *Args
	CmdArgs *Args
	Arg []string
	//FFmpegCmd *FFmpegCmd
	//FFprobeCmd *FFprobeCmd
	InputSlice []string
	Input string
	cwd string
}

func NewCmd() *Cmd {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &Cmd{
		cwd: cwd,
	}
	//cmd.args = Cfg().GetProfile(cmd.Profile)
	//cmd.FFmpegCmd = NewFFmpegCmd()
	//fmt.Printf("%+v\n", cmd.Args())

	return cmd
}
