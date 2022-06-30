package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ohzqq/avtools/avtools"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	flags   avtools.Options
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "avtools",
	Short: "",
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/avtools/config.yml)")
	rootCmd.PersistentFlags().StringVarP(&flags.Output, "output", "o", "", "set output name")
	rootCmd.PersistentFlags().StringVarP(&flags.Profile, "profile", "p", "default", "set profile")
	rootCmd.PersistentFlags().StringVarP(&flags.CoverFile, "artFile", "A", "", "set album art file")
	rootCmd.PersistentFlags().StringVarP(&flags.CueFile, "cuesheet", "C", "", "set cue sheet")
	rootCmd.PersistentFlags().StringVarP(&flags.MetaFile, "metaFile", "M", "", "set ffmetadata file")
	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "print ffmpeg/ffprobe command string")
	rootCmd.PersistentFlags().BoolVarP(&flags.Overwrite, "overwrite", "y", false, "overwrite existing files")
	rootCmd.PersistentFlags().BoolVarP(&flags.MetaSwitch, "meta", "m", false, "toggle ffmetadata")
	rootCmd.PersistentFlags().BoolVarP(&flags.CueSwitch, "cue", "c", false, "toggle cue sheet")
	//rootCmd.PersistentFlags().BoolVarP(&flags.ChapSwitch, "chaps", "s", false, "toggle chapters")
	rootCmd.PersistentFlags().BoolVarP(&flags.CoverSwitch, "albumArt", "a", false, "toggle album art")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".avtools" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".config/avtools"))
		viper.AddConfigPath(filepath.Join(home, "Sync/code/avtools/tmp/"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		avtools.InitProfiles(viper.Sub("defaults"), viper.Sub("profiles"))
	}
}
