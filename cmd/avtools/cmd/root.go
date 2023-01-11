package cmd

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/ohzqq/avtools/ff"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	outName    string
	inName     string
	proFile    string
	verbose    bool
	overwrite  bool
	filterFlag []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "tools for working with a/v files",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outName, "output", "o", "tmp", "")
	rootCmd.PersistentFlags().StringVarP(&proFile, "profile", "p", "default", "")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "")
	rootCmd.PersistentFlags().BoolVarP(&overwrite, "overwrite", "y", true, "")
	rootCmd.PersistentFlags().StringVarP(&inName, "input", "i", "", "input video")
	rootCmd.PersistentFlags().StringSliceVarP(&filterFlag, "filter", "f", []string{}, "")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cmd")
		path := filepath.Join(home, ".config/avtools/profiles.yml")
		ff.ReadConfig(path)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
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

	if cmd.Flags().Changed("filter") {
		for n, f := range FilterFlag() {
			ffCmd.Filters.Add(n, f)
		}
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
