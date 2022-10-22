package tool

type Flag struct {
	Args ArgFlag
	Bool BoolFlag
}

type BoolFlag struct {
	Meta      bool
	Cover     bool
	Cue       bool
	Chap      bool
	Json      bool
	Overwrite bool
	Verbose   bool
}

type ArgFlag struct {
	Profile string
	Start   string
	End     string
	Output  string
	ChapNo  int
	Input   string
	Cover   string
	Meta    string
	Cue     string
}
