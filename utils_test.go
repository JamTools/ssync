package main

import (
  "os"
  "io"
  "sort"
  "strings"
  "io/ioutil"
  "testing"
)

func TestNotIn(t *testing.T) {
  notInTests := [][][]string{
    { { "abc", "xyz", "123" }, { "abc", "xyz", "123" }, {} },
    { { "abc", "123" }, { "xyz", "123" }, { "abc" } },
    { { "xyz", "123" }, { "abc", "123" }, { "xyz" } },
    { { "abc", "xyz" }, {}, { "abc", "xyz" } },
    { {}, { "abc", "xyz" }, {} },
  }

  for _, v := range notInTests {
    for i := range v {
      // v[1], v[2] MUST be sorted for binary search
      sort.Strings(v[i])
    }

    r := notIn(v[0], v[1])
    if strings.Join(r, "\n") != strings.Join(v[2], "\n") {
      t.Errorf("Expected %v, got %v", v[2], r)
    }
  }
}

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
