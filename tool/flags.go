package tool

type Flags struct {
	Overwrite   bool
	Profile     string
	Start       string
	End         string
	Output      string
	ChapNo      int
	MetaSwitch  bool
	CoverSwitch bool
	CueSwitch   bool
	ChapSwitch  bool
	JsonSwitch  bool
	Verbose     bool
	Input       string
	CoverFile   string
	MetaFile    string
	CueFile     string
}
