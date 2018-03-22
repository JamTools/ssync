package main

import (
  "os"
  "testing"
)

func TestProcessFlagsVersion(t *testing.T) {
  os.Args = []string{"ssync", "-version"}
  defer func() { flagVersion = false }()

  if _, cont := processFlags(); cont == true {
    t.Errorf("Expected %v, got %v", false, cont)
  }
}

func TestProcessFlagsMissing(t *testing.T) {
  os.Args = []string{"ssync"}

  if _, cont := processFlags(); cont == true {
    t.Errorf("Expected %v, got %v", false, cont)
  }
}

func TestProcessFlags(t *testing.T) {
  os.Args = []string{"ssync", "label", "path1", "path2"}

  _, cont := processFlags()
  if cont == false {
    t.Errorf("Expected %v, got %v", true, cont)
  }
}
