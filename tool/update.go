package tool

import (
	"log"
	"path/filepath"

	"github.com/ohzqq/avtools/ffmpeg"
)

type Update struct {
	*Cmd
	FFmpeg *ffmpeg.Cmd
	cover  *Cmd
}

func NewUpdateCmd() *Update {
	return &Update{Cmd: NewCmd()}
}

func (u *Update) Parse() *Cmd {
	u.FFmpeg = u.Cmd.FFmpeg()

	u.FFmpeg.Stream()

	var cover *Cmd
	if u.Flag.Args.HasCover() {
		switch u.Cmd.Media.Input.Ext() {
		case ".m4b", ".m4a":
			cpath, err := filepath.Abs(u.Flag.Args.Cover)
			if err != nil {
				log.Fatal(err)
			}
			cover = NewCmd().
				Bin("AtomicParsley").
				SetArgs(u.output.Abs, "--artwork", cpath, "--overWrite")
		case ".mp3":
			u.FFmpeg.Input(u.Flag.Args.Cover)
			u.FFmpeg.AppendAudioParam("id3v2_version", "3")
			u.FFmpeg.AppendAudioParam("metadata:s:v", "title='Album cover'")
			u.FFmpeg.AppendAudioParam("metadata:s:v", "comment='Cover (front)'")
		}
	}

	u.Add(u.FFmpeg)

	if cover != nil {
		u.Add(cover)
	}

	return u.Cmd
}
