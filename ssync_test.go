package main

import (
  "os"
  "io"
  "io/ioutil"
  "testing"
)

func tmpFile(input string, t *testing.T) *os.File {
  in, err := ioutil.TempFile("", "")
  if err != nil {
    t.Fatal(err)
  }

  _, err = io.WriteString(in, input)
  if err != nil {
    t.Fatal(err)
  }

  _, err = in.Seek(0, os.SEEK_SET)
  if err != nil {
    t.Fatal(err)
  }

  return in
}
