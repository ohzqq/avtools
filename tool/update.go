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

	if u.HasCue() {
		u.FFmpeg.HasChapters()
		tmp := u.MkTmp()
		err := u.Media.FFmeta.Write(tmp)
		if err != nil {
			log.Fatal(err)
		}
		u.FFmpeg.FFmeta(tmp.Name())
	}

	var cover *Cmd
	if u.HasCover() {
		switch u.Cmd.Input.Ext {
		case ".m4b", ".m4a":
			u.FFmpeg.VN()
			cover = NewCmd().
				Bin("AtomicParsley").
				SetArgs(u.Output.Abs, "--artwork", u.Cover.Abs, "--overWrite")
		case ".mp3":
			u.FFmpeg.Input(u.Cover.Abs)
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
