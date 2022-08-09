package avtools

import (
	"log"
	"os/exec"
	"path/filepath"
)

func (cmd *ffmpegCmd) Update() {
	cmd.media.JsonMeta().Unmarshal()
	cmd.ParseOptions()

	var cmdExec *Cmd

	if cmd.opts.CoverFile != "" {
		switch codec := cmd.media.AudioCodec(); codec {
		case "aac":
			cmdExec = addAacCover(cmd.media.File, cmd.opts.CoverFile, cmd.opts.Verbose)
		case "mp3":
			cmdExec = cmd.addMp3Cover()
		}
	}

	if cmd.opts.MetaFile != "" {
		cmdExec = cmd.ParseArgs()
	}

	if cmd.opts.CueFile != "" {
		meta := cmd.media.GetFormat(cmd.opts.CueFile).Render()
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
