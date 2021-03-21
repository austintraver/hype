package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// inputFile represents the filepath to the input file
var inputFile string

// outputFile represents the filepath to the input file
var outputFile string

// force determines whether output should overwrite existing files
var force bool

// convertCmd represents the command to convert Markdown into HTML
var convertCmd = &cobra.Command{
	Use:   "convert [FILE]",
	Short: "Transform Markdown into HTML",
	Run:   handleConvert,
}

// handleConvert receives Markdown provided as input and outputs it as HTML
func handleConvert(cmd *cobra.Command, args []string) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeCharDevice) != 0 && len(args) > 0 {
		inputFile = args[0]
	} else {
		inputFile = viper.GetString("input")
	}

	// If the user supplied one argument, assume it is the name
	// of the file to be converted from Markdown to HTML

	// Confirm that the path to the file actually exists
	_, err := os.Stat(viper.GetString("input"))
	if os.IsNotExist(err) {
		message := fmt.Sprintf(
			"hype: file %s not found",
			inputFile,
		)
		_, err := os.Stderr.WriteString(message)
		if err != nil {
			panic(err)
		}
	}

	// Confirm that the path to the output file does not already exist
	_, err = os.Stat(viper.GetString("output"))
	if outputFile != os.Stdout.Name() {
		if os.IsExist(err) && !(viper.GetBool("force")) {
			message := fmt.Sprintf(
				"hype: file %s already exists, use --force to overwrite",
				outputFile,
			)
			_, err = os.Stderr.WriteString(message)
			if err != nil {
				panic(err)
			}
			return
		}
	}
	// Read Markdown from the input file
	markdown, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	html, err := Convert(markdown)
	if err != nil {
		panic(err)
	}

	// Write the resultant HTML to the output file
	err = os.WriteFile(outputFile, html, 644)
	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().SortFlags = false
	convertCmd.Flags().StringVarP(
		&inputFile,
		"input",
		"i",
		os.Stdin.Name(),
		"filepath to read Markdown input",
	)
	convertCmd.Flags().StringVarP(
		&outputFile,
		"output",
		"o",
		os.Stdout.Name(),
		"filepath to write HTML output",
	)
	convertCmd.Flags().BoolVarP(
		&force,
		"force",
		"f",
		force,
		"allow output to overwrite a file if it already exists "+
			"(default: false)",
	)

	err := viper.BindPFlags(convertCmd.Flags())
	if err != nil {
		panic(err)
	}
}
