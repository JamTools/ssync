package main

import (
  "os"
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

var flagMain *flag.FlagSet
var flagForcePath int
var flagConfirm, flagVersion bool

func init() {
  flagMain = flag.NewFlagSet("main", flag.ContinueOnError)

  // setup options
  flagMain.IntVar(&flagForcePath, "force", 0, "update modified using this path regardless" +
    " of modified timestamp (0=timestamp, 1=PATH1, 2=PATH2)")
  flagMain.BoolVar(&flagConfirm, "confirm", false, "confirm before deleting files")
  flagMain.BoolVar(&flagVersion, "version", false, "print program version, then exit")

  // --help
  flagMain.Usage = func() {
    fmt.Printf(printUsage, version, description, args)
    // print options from built-in flag helper
    flagMain.PrintDefaults()
    fmt.Println()
  }
}

// handle flag --help && --version
func processFlags() ([]string, bool) {
  // show --help unless args
  if len(os.Args) < 4 {
    flagMain.Usage()
    return os.Args, false
  }

  flagMain.Parse(os.Args[1:])
  a := flagMain.Args()

  // --version
  if flagVersion {
    fmt.Printf("%s\n", version)
    return a, false
  }

  return a, true
}
