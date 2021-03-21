/*
Copyright © 2021 Austin Traver <atraver@usc.edu>
It's open source, so make a *copy*, right?
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"os"

	"github.com/spf13/cobra"
)

// inputFile represents the filepath to the input file
var inputFile string

// outputFile represents the filepath to the input file
var outputFile string

// force determines whether output should overwrite existing files
var force bool

// basic determines whether extensions to Markdown syntax are ignored
var basic bool

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [FILE]",
	Short: "Transform Markdown into HTML",
	Run:   handleConvert,
}

// handleConvert converts Markdown provided from input
// into HTML, and writes it to output

// the underlying conversion from Markdown to HTML
func handleConvert(cmd *cobra.Command, args []string) {

	// Confirm that the path to the file actually exists
	_, err := os.Stat(inputFile)
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
	_, err = os.Stat(outputFile)
	if outputFile != os.Stdout.Name() {
		if os.IsExist(err) && !(force) {
			message := fmt.Sprintf(
				"error: file %s already exists, use --force to overwrite",
				outputFile,
			)
			_, err = os.Stderr.WriteString(message)
			if err != nil {
				panic(err)
			}
		}
	}

	// Read Markdown from the input file
	input, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	// Initialize a list of extensions and options for the Markdown parser
	var extensions = []goldmark.Extender{
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
	}
	var parserOptions = []parser.Option{
		parser.WithAutoHeadingID(),
	}
	var rendererOptions []renderer.Option

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
	converter := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(parserOptions...),
		goldmark.WithRendererOptions(rendererOptions...),
	)
	var html bytes.Buffer
	err = converter.Convert(input, &html)
	if err != nil {
		panic(err)
	}

	// Write the resultant HTML to the output file
	err = os.WriteFile(outputFile, html.Bytes(), 644)
	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(convertCmd)
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
	convertCmd.Flags().BoolVarP(
		&basic,
		"basic",
		"b",
		basic,
		"ignore extensions to Markdown syntax (default: false)",
	)

	err := viper.BindPFlags(convertCmd.Flags())
	if err != nil {
		panic(err)
	}
}
