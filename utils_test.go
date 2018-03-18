package main

import (
  "os"
  "io"
  "io/ioutil"
  "testing"
)

func testAskConfirm(input string, result bool, t *testing.T) {
  in, err := ioutil.TempFile("", "")
  if err != nil {
    t.Fatal(err)
  }
  defer in.Close()

  _, err = io.WriteString(in, input)
  if err != nil {
    t.Fatal(err)
  }

  _, err = in.Seek(0, os.SEEK_SET)
  if err != nil {
    t.Fatal(err)
  }

  r := askConfirm(in)
  if r != result {
    t.Errorf("Expected %v, got %v", result, r)
  }
}

func TestAskConfirm(t *testing.T) {
  askTests := map[bool][]string{
    true: { "y", "Y", "yes", "Yes", "YES" },
    false: { "n", "N", "no", "No", "NO", "1", "0", "A", "*", "z" },
  }

  for r := range askTests {
    for _, v := range askTests[r] {
      testAskConfirm(v, r, t)
    }
  }
}
