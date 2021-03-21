# `hype`

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

Configuration files can be placed either at `${XDG_CONFIG_HOME}/hyperc.yaml` or
`${XDG_CONFIG_HOME}/hype/hyperc.yaml`

Support for configuration files in the home directory is not supported at this
time.

An example configuration file is provided below:

```yaml
basic: false
```
