package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
)

// configFile represents the filepath to the configuration file for hype
var configFile string

// basic determines whether extensions to Markdown syntax are ignored
var basic bool

// verbose determines if additional output will be logged to stdout
var verbose = false

var converter goldmark.Markdown

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "hype",
	Short: "A pretty *hype* CLI to help convert Markdown to hypertext markup" +
		" language (HTML)",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().StringVarP(
		&configFile,
		"config",
		"c",
		configFile,
		"user configuration file \n "+
			"(default \"${XDG_CONFIG_HOME}/hyperc.yaml\")",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&basic,
		"basic",
		"b",
		basic,
		fmt.Sprintf(
			"ignore extensions to Markdown syntax \n (default %v)",
			basic,
		),
	)
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		verbose,
		fmt.Sprintf(
			"ignore extensions to Markdown syntax \n (default %v)",
			verbose,
		),
	)
	err := viper.BindPFlags(rootCmd.Flags())
	if err != nil {
		panic(err)
	}

}

// Convert accepts Markdown as input and returns the equivalent HTML as output
func Convert(input []byte) (output []byte, err error) {

	// Initialize a list of extensions and options for the Markdown parser
	var extensions = []goldmark.Extender{
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
		emoji.Emoji,
	}
	var parserOptions = []parser.Option{
		parser.WithAutoHeadingID(),
	}
	var rendererOptions = []renderer.Option{
		html.WithUnsafe(),
	}

	// If the user specified to use basic Markdown syntax,
	// remove support for the popular extensions that are by default
	// included
	if viper.Get("basic") == true {
		// if basic == true {
		extensions = nil
		parserOptions = nil
		rendererOptions = nil
	}

	// Construct the markdown interface
	converter = goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(parserOptions...),
		goldmark.WithRendererOptions(rendererOptions...),
	)

	var buffer bytes.Buffer
	err = converter.Convert(input, &buffer)
	output = buffer.Bytes()
	return
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
	}

	// read in environment variables set specifically for hype
	// e.g.: HYPE_BASIC
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil && viper.GetBool("verbose") {
		_, err = fmt.Fprintln(
			os.Stderr,
			"Using config file:",
			color.YellowString(viper.ConfigFileUsed()),
		)
		if err != nil {
			panic(err)
		}
	}

}
