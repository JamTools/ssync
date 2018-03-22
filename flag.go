package main

import (
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

var flagConfirm, flagVerbose, flagVersion bool

func init() {
  // setup options
  flag.BoolVar(&flagConfirm, "confirm", false, "confirm before deleting files")
  flag.BoolVar(&flagVerbose, "v", false, "verbose: print additional output")
  flag.BoolVar(&flagVersion, "version", false, "prints program version")

  // --help
  flag.Usage = func() {
    fmt.Printf(printUsage, version, description, args)
    // print options from built-in flag helper
    flag.PrintDefaults()
    fmt.Println()
  }
}

// handle flag --help && --version
func processFlags() ([]string, bool) {
  flag.Parse()
  a := flag.Args()

  // --version
  if flagVersion {
    fmt.Printf("%s\n", version)
    return a, false
  }

  // show --help unless args
  if len(a) != 3 {
    flag.Usage()
    return a, false
  }

  return a, true
}
