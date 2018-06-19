package main

import (
  "fmt"
  "flag"
)

const (
  version = "0.0.1"
  description = "Synchronize audio collection among group of individuals"
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

var flagForcePath int
var flagConfirm, flagVersion bool

func init() {
  // setup options
  flag.IntVar(&flagForcePath, "force", 0, "update modified using this path regardless" +
    " of modified timestamp (0=timestamp, 1=PATH1, 2=PATH2)")
  flag.BoolVar(&flagConfirm, "confirm", false, "confirm before deleting files")
  flag.BoolVar(&flagVersion, "version", false, "print program version, then exit")

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
