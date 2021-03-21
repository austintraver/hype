# hype

## Introduction

`hype` is a Markdown utility tool for the command-line interface. Hype can be
used to convert Markdown to HTML using the `convert` subcommand, as follows:

## Usage

```shell script
hype convert < input.md > output.md
```

Alternatively, input and output file locations can be specified using the `-i`
(`--input`) and `-o` (`--output`) flags respectively. For example:

```shell script
hype convert -i input.md -o output.md
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
