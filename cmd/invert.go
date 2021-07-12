package cmd

import (
	"github.com/spf13/cobra"

	md "github.com/austintraver/hype/hallmark"
)

var invertCmd = cobra.Command{
	Use:   "invert FILE",
	Short: "convert an HTML file into Markdown",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var converter = md.Converter{}
		converter.Use(md.Table())
		converter.Use(md.Strikethrough(""))
		converter.Use(md.TaskListItems())
		converter.ConvertString("")
	},
}

func init() {
	rootCmd.AddCommand()
}
