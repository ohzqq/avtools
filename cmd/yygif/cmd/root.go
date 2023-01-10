package cmd

import (
	"os"
	"path/filepath"

	"github.com/ohzqq/avtools/ff"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	scale      []string
	crop       []string
	eq         []string
	colortemp  []string
	filterFlag []string
	smartblur  string
	fps        string
	setpts     string
	bayerScale string
	startTime  string
	endTime    string
	proFile    string
	dither     string
	outName    string
	metadata   string
	verbose    bool
	yadif      bool
	overwrite  bool
)

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "A brief description of your application",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringSliceVarP(&scale, "scale", "a", []string{}, "")
	rootCmd.PersistentFlags().StringSliceVarP(&crop, "crop", "c", []string{}, "")
	rootCmd.PersistentFlags().StringSliceVarP(&eq, "eq", "e", []string{}, "")
	rootCmd.PersistentFlags().StringSliceVarP(&filterFlag, "filter", "f", []string{}, "")
	rootCmd.PersistentFlags().StringVar(&bayerScale, "bs", "3", "")
	rootCmd.PersistentFlags().StringVarP(&dither, "dither", "d", "bayer", "")
	rootCmd.PersistentFlags().StringVarP(&outName, "output", "o", "tmp", "")
	rootCmd.PersistentFlags().StringVar(&metadata, "meta", "metadata-default.yml", "")
	rootCmd.PersistentFlags().StringVarP(&startTime, "ss", "s", "", "")
	rootCmd.PersistentFlags().StringVarP(&endTime, "to", "t", "", "")
	rootCmd.PersistentFlags().StringVarP(&smartblur, "smartblur", "b", "0.5", "")
	rootCmd.PersistentFlags().StringVarP(&proFile, "profile", "p", "default", "")
	rootCmd.PersistentFlags().StringSliceVarP(&colortemp, "colortemp", "m", []string{}, "")
	rootCmd.PersistentFlags().StringVarP(&fps, "fps", "r", "", "")
	rootCmd.PersistentFlags().StringVar(&setpts, "setpts", "", "")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "")
	rootCmd.PersistentFlags().BoolVarP(&overwrite, "overwrite", "y", true, "")
	rootCmd.PersistentFlags().BoolVar(&yadif, "yadif", false, "")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	path := filepath.Join(home, ".config/avtools/profiles.yml")
	ff.ReadConfig(path)
	//fmt.Printf("%+V\n", yygif.Profiles)
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
