# `hype`

[![Go Reference](https://pkg.go.dev/badge/github.com/austintraver/hype.svg)](https://pkg.go.dev/github.com/austintraver/hype)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Introduction

`hype` is a CLI utility to help convert Markdown into HTML.

## Installation

As a Go module, `hype` can be installed like any traditional Go package.

```shell script
go install github.com/austintraver/hype
```

## Usage

Using `hype` to convert a Markdown file into HTML can be performed using the
`convert` subcommand, whose usage is as follows:

```shell script
hype convert < input.md
```

Alternatively, input and output file locations can be specified using the `-i`
(`--input`) and `-o` (`--output`) flags respectively. For example:

```shell script
hype convert -i input.md -o output.html
```

## Configuration

`hype` searches for a configuration file located at
`${XDG_CONFIG_HOME}/hyperc.yaml`. If the file is found, `hype` will treat each
configuration specified as if the user had provided the equivalent `--flag` on
the command line.

For this reason, any `--flag` that is *actually* provided on the command line
will take priority, overriding the value set within the configuration file.

Support for configuration files in the home directory is not supported at this
time.

An example configuration file is provided below:

```yaml
# Use basic Markdown syntax, removing support for common extensions
basic: false

# When running the HTTP server to preview Markdown files 
# via the `server`, subcommand, listen for connections on port 1411
port: 1411
```
