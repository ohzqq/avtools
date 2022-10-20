package tool

import (
	"log"

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

	if u.flag.Args.HasCue() {
		u.FFmpeg.HasChapters()
		tmp := u.MkTmp()
		err := u.Args.Media.Meta.Write(tmp)
		if err != nil {
			log.Fatal(err)
		}
		u.FFmpeg.FFmeta(tmp.Name())
	}

	var cover *Cmd
	if u.flag.Args.HasCover() {
		switch u.Cmd.Args.Input.Ext {
		case ".m4b", ".m4a":
			u.FFmpeg.VN()
			cover = NewCmd().
				Bin("AtomicParsley").
				SetArgs(u.Args.Output(), "--artwork", u.Args.Cover.Abs, "--overWrite")
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
