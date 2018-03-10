package main

import (
  "os"
  "fmt"
  "flag"
)

const (
  version = "0.0.1"
  description = "smart syncs files between two paths"
)

// positional args template
const args = `
Positional Args:
  LABEL           give label for subsequent use
  PATH1           1st directory path
  PATH2           2nd directory path
`

// usage template
const printUsage = `
ssync v%s
%s

Usage: ssync [OPTIONS] LABEL PATH1 PATH2
%s
Options:
`

var flagConfirm bool

// init: handle flag --help && --version
func init() {
  // setup options
  var showVersion bool
  flag.BoolVar(&showVersion, "version", false, "prints program version")
  flag.BoolVar(&flagConfirm, "confirm", false, "confirm before deleting files")

  // --help
  flag.Usage = func() {
    // build & print usage
    fmt.Printf(printUsage, version, description, args)
    // print options from built-in flag helper
    flag.PrintDefaults()
    fmt.Println()
  }

  flag.Parse()

  // --version
  if showVersion {
    fmt.Printf("%s\n", version)
    os.Exit(1)
  }

  // show --help unless args
  if len(flag.Args()) != 3 {
    flag.Usage()
    os.Exit(1)
  }
}
