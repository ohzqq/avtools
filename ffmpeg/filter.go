package ffmpeg

import "strings"

type Filters []Filter

func (f Filters) String() string {
	var filters []string
	for _, f := range f {
		filters = append(filters, f.String())
	}
	return strings.Join(filters, ",")
}

type Filter struct {
	Name   string
	Params [][]string
}

func NewFilter(name string) Filter {
	return Filter{
		Name: name,
	}
}

func (f *Filter) Set(val ...string) {
	f.Params = append(f.Params, val)
}

func (f Filter) String() string {
	var params []string
	for _, val := range f.Params {
		params = append(params, strings.Join(val, "="))
	}

	p := strings.Join(params, ":")

	if f.Name == "" {
		return p
	}

	return f.Name + "=" + p
}
