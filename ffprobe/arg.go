package ffprobe

import "strings"

const (
	ofJson  = "json=c=1"
	ofPlain = "default=nk=1:nw=1"
)

type Args struct {
	logLevel  []string
	pretty    bool
	streams   []string
	entries   Entries
	showChaps bool
	format    []string
	input     string
}

type Entries map[string][]string

func NewEntries() Entries {
	return Entries{
		"format": string{},
		"stream": string{},
	}
}

func (e Entries) AddFormat(vals ...string) {
	//e["format"] = append(e["format"], vals...)
	e.Add("format", vals...)
}

func (e Entries) AddStream(vals ...string) {
	//e["stream"] = append(e["stream"], vals...)
	e.Add("stream", vals...)
}

func (e Entries) Add(key, vals ...string) {
	e[key] = append(e[key], vals...)
}

func (e Entries) String() string {
	var entries []string
	for key, val := range e {
		entry := []string{key}

		if len(val) > 0 {
			v := strings.Join(val, ",")
			entry = append(entry, v)
		}

		entry = strings.Join(entry, "=")

		entries = append(entries, entry)
	}

	return strings.Join(entries, ":")
}

func NewArgs() *Args {
	return &Args{
		logLevel: []string{verbose},
		entries:  NewEntries(),
		format:   []string{writer},
	}
}

func (c Args) HasLogLevel() bool {
	return len(c.logLevel) > 1
}

func (c *Args) LogLevel(l string) *Args {
	c.logLevel = append(c.logLevel, l)
	return c
}

func (c *Args) Pretty() *Args {
	c.pretty = true
	return c
}

func (c Args) HasStreams() bool {
	return len(c.streams) > 0
}

func (c *Args) Stream(s ...string) *Args {
	c.streams = append(c.streams, s...)
	return c
}

func (c *Args) ShowChapters() *Args {
	c.showChapters = true
	return c
}

func (c Args) HasEntries() bool {
	return len(c.Entries) > 0
}

func (c *Args) Entry(key string, vals ...string) *Args {
	c.entries.Add(key, vals...)
	return c
}

func (c *Args) FormatEntry(vals ...string) *Args {
	c.entries.AddFormat(vals...)
	return c
}

func (c *Args) StreamEntry(vals ...string) *Args {
	c.entries.AddStream(vals...)
	return c
}

func (c Args) HasFormat() bool {
	return len(c.format) > 1
}

func (c *Args) Json() *Args {
	c.format = append(c.format, ofJson)
}

func (c *Args) Plain() *Args {
	c.format = append(c.format, ofPlain)
}

func (c Args) HasInput() bool {
	return c.input != ""
}

func (c *Args) Input(i string) *Args {
	c.input = i
	return c
}
