package tool

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/ohzqq/avtools/tool/ffmpeg"
)

type Update struct {
	*Cmd
	FFmpeg *ffmpeg.Cmd
}

func NewUpdateCmd() *Update {
	return &Update{Cmd: NewerCmd()}
}

//func (u *Update) SetFlags(f Flag) *Update {
//  u.Flag = f
//  return u
//}

func (u *Update) ParseArgs() *Update {
	u.Cmd.SetFlags(u.Flag)
	u.FFmpeg = u.FFmpegCmd()

	u.FFmpeg.Stream()

	var coverCmd *exec.Cmd
	if u.Flag.Args.HasCover() {
		switch u.Cmd.Media.AudioCodec() {
		case "aac":
			coverCmd = AacCover(u.Flag.Args.Input, u.Flag.Args.Cover)
		case "mp3":
			u.FFmpeg.Input(u.Flag.Args.Cover)
			u.FFmpeg.AppendAudioParam("id3v2_version", "3")
			u.FFmpeg.AppendAudioParam("metadata:s:v", "title='Album cover'")
			u.FFmpeg.AppendAudioParam("metadata:s:v", "comment='Cover (front)'")
		}
	}

	if u.Flag.Args.HasCue() {
		u.FFmpeg.SetHasChapters()
		meta := u.Cmd.Media.Meta().MarshalTo("ffmeta").Bytes()
		u.Cmd.tmpFile = TmpFile(meta)
		u.FFmpeg.FFmeta(u.Cmd.tmpFile)
	}

	args, err := u.FFmpeg.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}
	u.Cmd.New("ffmpeg", args)

	if coverCmd != nil {
		u.AddCmd(coverCmd)
	}

	return u
}

func (cmd *ffmpegCmd) Update() {
	//cmd.media.JsonMeta().Unmarshal()
	cmd.ParseOptions()

	var cmdExec *Cmd

	if cmd.opts.CoverFile != "" {
		//switch codec := cmd.media.AudioCodec(); codec {
		//case "aac":
		//  cmdExec = addAacCover(cmd.media.File, cmd.opts.CoverFile, cmd.opts.Verbose)
		//case "mp3":
		//  cmdExec = cmd.addMp3Cover()
		//}

	}

	if cmd.opts.MetaFile != "" {
		cmdExec = cmd.ParseArgs()
	}

	if cmd.opts.CueFile != "" {
		meta := cmd.media.GetFile("cue").Render()
		tmp := TmpFile([]byte(meta.data))

		cmd.AppendMapArg("post", "i", tmp)
		cmd.AppendMapArg("post", "map_chapters", "1")
		cmdExec = cmd.ParseArgs()
	}

	if cmd.opts.CueFile == "" && cmd.opts.MetaFile == "" && cmd.opts.CoverFile == "" {
		log.Fatal("the update command requires *something* to update")
	}

	cmdExec.Run()
}

func AacCover(file, cover string) *exec.Cmd {
	cpath, err := filepath.Abs(cover)
	if err != nil {
		log.Fatal(err)
	}
	return exec.Command("AtomicParsley", file, "--artwork", cpath, "--overWrite")
}

func addAacCover(file, cover string, verbose bool) *Cmd {
	cpath, err := filepath.Abs(cover)
	if err != nil {
		log.Fatal(err)
	}
	return NewCmd(
		exec.Command("AtomicParsley", file, "--artwork", cpath, "--overWrite"),
		verbose,
	)
}

func (cmd *ffmpegCmd) addMp3Cover() *Cmd {
	//cmd := ffmpegCmd{}
	cmd.VideoCodec = ""
	cmd.AppendMapArg("audioParams", "id3v2_version", "3")
	cmd.AppendMapArg("audioParams", "metadata:s:v", "title='Album cover'")
	cmd.AppendMapArg("audioParams", "metadata:s:v", "comment='Cover (front)'")
	cmd.Output = "with-cover"
	return cmd.ParseArgs()
}
