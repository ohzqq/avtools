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
	out := NewOutput(u.flag.Args.Output)
	u.FFmpeg = u.Cmd.FFmpeg()

	u.FFmpeg.Stream()

	var cover *Cmd
	if u.flag.Args.HasCover() {
		u.FFmpeg.VN()
		switch u.Cmd.Args.Input.Ext {
		case ".m4b", ".m4a":
			cpath, err := filepath.Abs(u.flag.Args.Cover)
			if err != nil {
				log.Fatal(err)
			}
			cover = NewCmd().
				Bin("AtomicParsley").
				SetArgs(out.String(), "--artwork", cpath, "--overWrite")
		case ".mp3":
			u.FFmpeg.Input(u.flag.Args.Cover)
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
