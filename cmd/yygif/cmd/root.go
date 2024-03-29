package cmd

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/ohzqq/avtools"
	"github.com/ohzqq/avtools/ff"
	"github.com/ohzqq/avtools/media"
	"github.com/ohzqq/avtools/meta"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	outName   string
	inName    string
	startTime string
	endTime   string
	metadata  string
	proFile   string
	verbose   bool
	overwrite bool
	chapterNo int
	// filters
	scale      []string
	crop       []string
	eq         []string
	colortemp  []string
	filterFlag []string
	smartblur  string
	fps        string
	setpts     string
	bayerScale string
	dither     string
	yadif      bool
)

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "make gifs!",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		meta := getMedia(cmd)
		//c := media.CutChapter(meta, meta.Chapters()[0])
		//c.Run()

		for _, ch := range getChapters(cmd, meta) {
			g := media.CutChapter(meta, ch)
			if c, ok := ch.Tags["crop"]; ok {
				crop := strings.Split(c, ":")
				g.Filters.Set("crop", crop...)
			}
			//  g := gif.MkGif(meta.Input.Abs, ch)
			ff := ParseFlags(cmd, &g)
			ff.Compile()

			if cmd.Flags().Changed("verbose") {
				fmt.Println(ff.String())
			}

			err := ff.Run()
			if err != nil {
				log.Fatal(err)
			}
		}

	},
}

func LoadGifMeta(input string) *media.Media {
	meta := meta.LoadIni(input)
	src := avtools.NewMedia().SetMeta(meta)
	vid := meta.Tags()["title"]
	return &media.Media{
		Media:   src,
		Input:   media.NewFile(vid),
		Profile: "gif",
	}
}

func getMedia(cmd *cobra.Command) *media.Media {
	var meta *media.Media
	if cmd.Flags().Changed("input") {
		meta = media.New(inName)
		meta.Profile = "gif"
	}
	if cmd.Flags().Changed("meta") {
		if cmd.Flags().Changed("input") {
			meta.LoadIni(metadata)
		} else {
			meta = LoadGifMeta(metadata)
		}
	}
	if meta == nil {
		log.Fatal("either a video or meta file is required")
	}
	return meta
}

func getChapters(cmd *cobra.Command, meta *media.Media) []*avtools.Chapter {
	var chapters []*avtools.Chapter
	switch {
	case cmd.Flags().Changed("num"):
		ch := meta.GetChapter(chapterNo)
		chapters = append(chapters, ch)
	case cmd.Flags().Changed("ss"), cmd.Flags().Changed("to"):
		if !cmd.Flags().Changed("input") {
			log.Fatal("this command requires an input video src")
		}

		start := "00:00"
		if cmd.Flags().Changed("ss") {
			start = startTime
		}

		end := meta.GetTag("duration")
		if cmd.Flags().Changed("to") {
			end = endTime
		}

		ch := &avtools.Chapter{
			Title: fmt.Sprintf("%s-%s-%s", meta.Input.Name, start, end),
		}
		ch.SS(start).To(end)
		chapters = append(chapters, ch)
	default:
		if meta.HasChapters() {
			chapters = meta.Chapters()
		} else {
			log.Fatal("no gifs to make")
		}
	}

	return chapters
}

func ParseFlags(cmd *cobra.Command, ffCmd *ff.Cmd) *ff.Cmd {
	orig := ffCmd
	if cmd.Flags().Changed("profile") {
		//ffCmd = ff.New(proFile)
		c := ff.New(proFile)
		orig = &c
		ffCmd.In(orig.File, orig.Input.Args)
		//ffCmd.Output = orig.Output.Merge(ffCmd.Output.Args)
		ffCmd.Output = ff.NewOutput(orig.Output.Args, ffCmd.Output.Args)
	}

	if cmd.Flags().Changed("output") {
		ffCmd.Output.Name(outName).Pad("")
	}

	if cmd.Flags().Changed("verbose") {
		ffCmd.Verbose()
	}

	if !cmd.Flags().Changed("overwrite") {
		ffCmd.Overwrite()
	}

	if cmd.Flags().Changed("ss") {
		ffCmd.Start(startTime)
	}

	if cmd.Flags().Changed("to") {
		ffCmd.End(endTime)
	}

	if cmd.Flags().Changed("fps") {
		filter := ff.Fps(fps)
		ffCmd.Filters.Add("fps", filter)
	}

	if cmd.Flags().Changed("setpts") {
		filter := ff.Setpts(setpts)
		ffCmd.Filters.Add("setpts", filter)
	}

	if cmd.Flags().Changed("crop") {
		ffCmd.Filters.Set("crop", crop...)
	}

	if cmd.Flags().Changed("scale") {
		ffCmd.Filters.Set("scale", scale...)
	}

	if cmd.Flags().Changed("yadif") {
		filter := ff.Yadif()
		ffCmd.Filters.Add("yadif", filter)
	}

	if cmd.Flags().Changed("colortemp") {
		filter := ff.Colortemp(colortemp...)
		ffCmd.Filters.Add("colortemperature", filter)
	}

	if cmd.Flags().Changed("eq") {
		filter := ff.Eq(eq...)
		ffCmd.Filters.Add("eq", filter)
	}

	if cmd.Flags().Changed("smartblur") {
		filter := ff.Smartblur(smartblur)
		ffCmd.Filters.Add("smartblur", filter)
	}

	if cmd.Flags().Changed("filter") {
		for n, f := range FilterFlag() {
			ffCmd.Filters.Add(n, f)
		}
	}

	var gif []string
	if cmd.Flags().Changed("bs") {
		gif = append(gif, "bs="+bayerScale)
	}
	if cmd.Flags().Changed("dither") {
		gif = append(gif, "dither="+dither)
	}
	if len(gif) > 0 {
		ffCmd.Filters.Set("palette", gif...)
	}

	return ffCmd
}

func FilterFlag() ff.Filters {
	filters := make(ff.Filters)
	for _, filter := range filterFlag {
		split := strings.Split(filter, ":")
		var name, args string
		switch l := len(split); l {
		case 2:
			args = split[1]
			fallthrough
		case 1:
			name = split[0]
		}
		f := ff.NewFilter(args)
		filters[name] = f
	}
	return filters
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	mime.AddExtensionType(".ini", "text/plain")
	mime.AddExtensionType(".cue", "text/plain")
	mime.AddExtensionType(".m4b", "audio/mp4")

	rootCmd.PersistentFlags().StringVarP(&metadata, "meta", "m", "", "")
	rootCmd.PersistentFlags().StringVarP(&outName, "output", "o", "tmp", "")
	rootCmd.PersistentFlags().StringVarP(&inName, "input", "i", "", "input video")
	rootCmd.PersistentFlags().StringVarP(&proFile, "profile", "p", "default", "")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "")
	rootCmd.PersistentFlags().BoolVarP(&overwrite, "overwrite", "y", true, "")
	rootCmd.PersistentFlags().StringVarP(&startTime, "ss", "s", "", "")
	rootCmd.PersistentFlags().StringVarP(&endTime, "to", "t", "", "")
	rootCmd.Flags().IntVarP(&chapterNo, "num", "n", 1, "chapter number")

	// filters
	rootCmd.PersistentFlags().StringSliceVarP(&filterFlag, "filter", "f", []string{}, "")
	rootCmd.PersistentFlags().StringSliceVarP(&scale, "scale", "a", []string{}, "")
	rootCmd.PersistentFlags().StringSliceVarP(&crop, "crop", "c", []string{}, "")
	rootCmd.PersistentFlags().StringSliceVarP(&eq, "eq", "e", []string{}, "")
	rootCmd.PersistentFlags().StringVar(&bayerScale, "bs", "3", "")
	rootCmd.PersistentFlags().StringVarP(&dither, "dither", "d", "bayer", "")
	rootCmd.PersistentFlags().StringVarP(&smartblur, "smartblur", "b", "0.5", "")
	rootCmd.PersistentFlags().StringSliceVarP(&colortemp, "colortemp", "k", []string{}, "")
	rootCmd.PersistentFlags().StringVarP(&fps, "fps", "r", "", "")
	rootCmd.PersistentFlags().StringVar(&setpts, "setpts", "", "")
	rootCmd.PersistentFlags().BoolVar(&yadif, "yadif", false, "")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	path := filepath.Join(home, ".config/avtools/profiles.yml")
	ff.ReadConfig(path)
	viper.SetConfigName("profiles")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(home, ".config/avtools"))
	//fmt.Printf("%+V\n", yygif.Profiles)
}
