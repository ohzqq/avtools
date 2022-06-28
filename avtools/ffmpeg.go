package avtools

import (
	"os"
	//"os/exec"
	//"io/fs"
	//"log"
	//"fmt"
	//"bytes"
	//"strings"
	//"strconv"
	//"path/filepath"

	//"cs.opensource.google/go/x/exp/slices"
	//"github.com/alessio/shellescape"
)

type FFmpegCmd struct {
	//cmd *exec.Cmd
	args []string
	input []string
	media []string
	cover string
	meta string
	tmpFile *os.File
}

func NewFFmpegCmd() *FFmpegCmd {
	return &FFmpegCmd{}
}

