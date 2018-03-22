package main

import (
  "testing"
)

func TestProcessFlags(t *testing.T) {
  flagVersion = true
  _, cont := processFlags([]string{})
  if cont == true {
    t.Errorf("Expected %v, got %v", false, cont)
  }
  flagVersion = false

  _, cont = processFlags([]string{})
  if cont == true {
    t.Errorf("Expected %v, got %v", false, cont)
  }

  _, cont = processFlags([]string{"label", "path1", "path2"})
  if cont == false {
    t.Errorf("Expected %v, got %v", true, cont)
  }
}
