package ff

import (
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Filters map[string]Filter
type Filter ffmpeg.KwArgs

type filter func(*ffmpeg.Stream) *ffmpeg.Stream

func newFilter(name string, args ffmpeg.KwArgs) filter {
	return func(stream *ffmpeg.Stream) *ffmpeg.Stream {
		return stream.Filter(name, ffmpeg.Args{}, args)
	}
}

func NewFilter(args ...string) Filter {
	return Filter(ArgsToKwArgs(args))
}

func ArgsToKwArgs(args []string) ffmpeg.KwArgs {
	parsed := make(ffmpeg.KwArgs)
	for _, arg := range args {
		var key, val string
		a := strings.Split(arg, "=")
		switch len(a) {
		case 2:
			val = a[1]
			fallthrough
		case 1:
			key = a[0]
		}
		parsed[key] = val
	}
	return parsed
}

func MergeFilters(filters []Filters) Filters {
	a := Filters{}
	for _, b := range filters {
		for c := range b {
			a[c] = b[c]
		}
	}
	return a
}

func (f Filters) Set(name string, args ...string) {
	filter := f.Get(name)
	merged := filter.MergeArgs(args...)
	f[name] = Filter(merged)
}

func (f Filters) Add(name string, filt Filter) {
	filter := f.Get(name)
	merged := filter.Merge(filt)
	f[name] = Filter(merged)
}

func (f Filters) Get(name string) Filter {
	if filter, ok := f[name]; ok {
		return filter
	}
	return make(Filter)
}

func (f Filter) Args() ffmpeg.KwArgs {
	return ffmpeg.KwArgs(f)
}

func (f Filter) MergeArgs(args ...string) Filter {
	parsed := Filter(ArgsToKwArgs(args))
	return f.Merge(parsed)
}

func (f Filter) Merge(filter Filter) Filter {
	args := []ffmpeg.KwArgs{f.Args(), ffmpeg.KwArgs(filter)}
	merged := ffmpeg.MergeKwArgs(args)
	return Filter(merged)
}

var filterOrder = []string{
	"ffmetadata",
	"yadif",
	"thumbnail",
	"fps",
	"setpts",
	"colortemperature",
	"eq",
	"crop",
	"scale",
	"smartblur",
	"palette",
}

func (f Filters) Compile() []filter {
	var filters []filter

	for _, name := range filterOrder {
		if args, ok := f[name]; ok {
			var filter filter
			switch name {
			case "crop":
				filter = Crop(ffmpeg.KwArgs(args))
			case "scale":
				filter = Scale(ffmpeg.KwArgs(args))
			case "palette":
				filter = Palette(ffmpeg.KwArgs(args))
			default:
				filter = newFilter(name, args.Args())
			}
			filters = append(filters, filter)
		}
	}

	return filters
}

// colortemperature filter arguments
// see https://ffmpeg.org/ffmpeg-filters.html#colortemperature
// for details about the options
func Colortemp(args ...string) Filter {
	var filter []string
	for _, arg := range args {
		split := strings.Split(arg, "=")
		val := split[1]
		switch key := split[0]; key {
		case "t", "temp":
			filter = append(filter, "temperature="+val)
		case "m", "mix":
			filter = append(filter, "mix="+val)
		case "pl":
			filter = append(filter, "pl="+val)
		default:
			filter = append(filter, key+"="+val)
		}
	}
	filter = []string(filter)
	return NewFilter(filter...)
}

// paletteuse filter arguments
// see https://ffmpeg.org/ffmpeg-filters.html#paletteuse
// for details about the options
// palettegen filter arguments
// see https://ffmpeg.org/ffmpeg-filters.html#palettegen-1
// for details about the options
func Palette(args ffmpeg.KwArgs) filter {
	return func(stream *ffmpeg.Stream) *ffmpeg.Stream {
		genArgs := make(ffmpeg.KwArgs)
		useArgs := make(ffmpeg.KwArgs)

		for key, val := range args {
			switch key {
			case "mx", "max_colors":
				genArgs["max_colors"] = val
			case "rt", "reserve_transparent":
				genArgs["reserve_transparent"] = val
			case "tc", "transparency_color":
				genArgs["transparency_color"] = val
			case "s", "stats_mode", "sm":
				if val != "full" || val != "diff" || val != "single" {
					val = "full"
				}
				genArgs["stats_mode"] = val
			case "n", "new":
				useArgs["new"] = val
			case "bs", "bayer_scale":
				useArgs["bayer_scale"] = val
			case "d", "dither":
				if val != "bayer" && val != "heckbert" && val != "floyd_steinberg" && val != "sierra2" && val != "sierra2_4a" {
					val = "floyd_steinberg"
				}
				useArgs["dither"] = val
			case "dm", "diff", "diff_mode":
				if val != "rectangle" {
					val = "none"
				}
				useArgs["diff_mode"] = val
			case "at", "alpha", "alpha_threshold":
				useArgs["alpha_threshold"] = val
			}
		}

		var hasGen, hasUse bool
		if len(genArgs) > 0 {
			hasGen = true
		}
		if len(useArgs) > 0 {
			hasUse = true
		}

		var ns *ffmpeg.Stream
		if hasGen || hasUse {
			split := stream.Split()
			use := ffmpeg.NewStream(split, "FilterableStream", "0", "")
			use = use.Filter("palettegen", ffmpeg.Args{}, genArgs)

			ns = ffmpeg.Filter(
				[]*ffmpeg.Stream{split.Get("1"), use},
				"paletteuse",
				ffmpeg.Args{},
				useArgs,
			)
		}
		return ns
	}
}

// eq filter arguments
// see https://ffmpeg.org/ffmpeg-filters.html#eq
// for details about the options
func Eq(args ...string) Filter {
	var filter []string
	if len(args) > 0 {
		for _, arg := range args {
			split := strings.Split(arg, "=")
			val := split[1]
			switch key := split[0]; key {
			case "b":
				filter = append(filter, "brightness="+val)
			case "c":
				filter = append(filter, "contrast="+val)
			case "g":
				filter = append(filter, "gamma="+val)
			case "s":
				filter = append(filter, "saturation="+val)
			case "gr":
				filter = append(filter, "gamma_r="+val)
			case "gg":
				filter = append(filter, "gamma_g="+val)
			case "gb":
				filter = append(filter, "gamma_b="+val)
			case "gw":
				filter = append(filter, "gamma_weight="+val)
			default:
				filter = append(filter, key+"="+val)
			}
		}
	}
	return NewFilter(filter...)
}

// scale filter args
// see https://ffmpeg.org/ffmpeg-filters.html#scale-1
// for details about the arguments.
func Scale(args ffmpeg.KwArgs) filter {
	if _, ok := args["w"]; !ok {
		args["w"] = "-2"
	}
	if _, ok := args["h"]; !ok {
		args["h"] = "-2"
	}

	return newFilter("scale", args)
}

func Crop(args ffmpeg.KwArgs) filter {
	if _, ok := args["w"]; !ok {
		args["w"] = "iw"
	}
	if _, ok := args["h"]; !ok {
		args["h"] = "ih"
	}

	return newFilter("crop", args)
}

// smartblur filter args
// see https://ffmpeg.org/ffmpeg-filters.html#smartblur-1
// for details about the options
func Smartblur(args ...string) Filter {
	var filter []string

	if len(args) == 1 {
		args = []string{"ls=" + args[0]}
	}

	if len(args) > 0 {
		for _, arg := range args {
			split := strings.Split(arg, "=")
			val := split[1]
			switch key := split[0]; key {
			case "ls", "luma_strength":
				filter = append(filter, "ls="+val)
			default:
				filter = append(filter, key+"="+val)
			}
		}
	}

	return NewFilter(filter...)
}

func Fps(args ...string) Filter {
	var filter []string

	if len(args) == 1 {
		args = []string{"fps=" + args[0]}
	}

	if len(args) > 0 {
		filter = args
	}

	return NewFilter(filter...)
}

func Setpts(args ...string) Filter {
	var filter []string
	if len(args) == 1 {
		args = []string{args[0] + "*PTS"}
	}
	if len(args) > 0 {
		filter = args
	}
	return NewFilter(filter...)
}

func Yadif(args ...string) Filter {
	var filter []string
	if len(args) > 0 {
		filter = args
	}
	return NewFilter(filter...)
}
