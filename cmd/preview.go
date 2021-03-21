package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// port determines the port for the HTTP server to listen for requests on
var port = 1411

var rootDir string

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview [DIR]",
	Short: "Preview the Markdown files within directory DIR in a web browser",
	Long: `Preview Markdown files located in directory DIR in a web browser.
Launches a server listening on port 1411.

Usage: access file FILE within directory DIR by visiting the following URL:
http://localhost:1411/FILE`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) != 0 {
			err := cmd.Usage()
			if err != nil {
				panic(err)
			}
		}
		port = viper.GetInt("port")

		// Register a handler for the HTTP server
		http.HandleFunc("/", handleIt)

		address := fmt.Sprintf(":%v", port)

		rootDir, err = filepath.Abs(viper.GetString("root"))
		if err != nil {
			panic(err)
		}

		if verbose {
			fmt.Printf(
				"Listening for connections on port " +
					color.BlueString(fmt.Sprintf(":%v", viper.GetInt("port"))) +
					"\n",
			)
			fmt.Printf(
				"Serving files under root directory: %v\n",
				color.YellowString(rootDir),
			)
		}

		// Begin listening for client connections to the HTTP server
		log.Fatal(http.ListenAndServe(address, http.HandlerFunc(handleIt)))
	},
}

// handleIt handles the current HTTP request
func handleIt(writer http.ResponseWriter, request *http.Request) {

	if verbose {
		// When we receive a new request, notify the client
		fmt.Printf(
			"Received request for file %v from host %v",
			color.YellowString(request.RequestURI),
			color.BlueString(request.RemoteAddr),
		)
	}

	ifile := path.Join(
		rootDir,
		filepath.Clean(request.URL.Path),
	)

	// If we cannot find the file they are asking for, notify the client
	_, err := os.Stat(ifile)
	if errors.Is(err, os.ErrNotExist) {
		writer.WriteHeader(404)
		message := fmt.Sprintf("<h1>Error 404</h1>\n"+
			"<p>File %v not found.</p>",
			request.URL.Path,
		)
		_, err = writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
		return
	} else if err != nil {
		panic(err)
	}

	// Read the Markdown contained in the input file
	buffer, err := os.ReadFile(ifile)
	if err != nil {
		writer.WriteHeader(500)
		message := fmt.Sprintf("<h1>Error 500</h1>\n"+
			"<p>%v</p>",
			err.Error(),
		)
		_, err = writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
		return
	}

	// Convert the Markdown into HTML
	html, err := Convert(buffer)
	if err != nil {
		writer.WriteHeader(500)
		message := fmt.Sprintf("<h1>Error 500</h1>\n"+
			"<p>%v</p>",
			err.Error(),
		)
		_, err = writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
		return
	}

	// Output the resultant HTML
	_, err = writer.Write(html)
	if err != nil {
		writer.WriteHeader(500)
		message := fmt.Sprintf("<h1>Error 500</h1>\n"+
			"<p>%v</p>",
			err.Error(),
		)
		_, err = writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
		return
	}
}

func init() {
	rootCmd.AddCommand(previewCmd)
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	previewCmd.Flags().SortFlags = false
	previewCmd.Flags().StringVarP(
		&rootDir,
		"root",
		"r",
		cwd,
		"root directory for the server to server files from\n",
	)
	previewCmd.Flags().IntVarP(
		&port,
		"port",
		"p",
		1411,
		"port for the server to listen on\n",
	)
	err = viper.BindPFlags(previewCmd.Flags())
	if err != nil {
		panic(err)
	}

}
