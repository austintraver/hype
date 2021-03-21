package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

// configFile represents the filepath to the configuration file for hype
var configFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "hype [FILE]",
	Short: "A pretty *hype* CLI to help convert Markdown to hypertext markup" +
		" language (HTML)",
	Run: handleConvert,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(
		&configFile,
		"config",
		"c",
		configFile,
		"user configuration file \n(default: "+
			"\"${XDG_CONFIG_HOME}/hyperc.yaml\" or \n"+
			"\"${XDG_CONFIG_HOME}/hype/.hyperc.yaml\")",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in with the name "hyperc" (without extension).
		viper.SetConfigName("hyperc")

		xdgConfigDir, set := os.LookupEnv("XDG_CONFIG_HOME")
		if !(set) {
			xdgConfigDir = path.Join(home, ".config")
		}
		viper.AddConfigPath(xdgConfigDir)
		hypeConfigDir := path.Join(xdgConfigDir, "hype")
		viper.AddConfigPath(hypeConfigDir)
		_, err = os.Stat(hypeConfigDir)
		if os.IsExist(err) {
			fmt.Printf("Using config directory: %v\n", hypeConfigDir)
			viper.AddConfigPath(hypeConfigDir)
			viper.SetConfigName("hyperc")
		}
	}

	// read in environment variables set specifically for hype
	// e.g.: HYPE_BASIC
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		_, err = fmt.Fprintln(os.Stderr, "Using config file:",
			viper.ConfigFileUsed())
		if err != nil {
			panic(err)
		}
	}

}
