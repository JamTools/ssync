package main

import (
  "os"
  "sort"
  "strings"
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

func TestAskConfirm(t *testing.T) {
  askTests := map[bool][]string{
    true: { "y", "Y", "yes", "Yes", "YES" },
    false: { "n", "N", "no", "No", "NO", "1", "0", "A", "*", "z" },
  }

  for result := range askTests {
    for _, v := range askTests[result] {
      tmpFile(t, v, func(in *os.File){
        r := askConfirm(in)
        if r != result {
          t.Errorf("Expected %v, got %v", result, r)
        }
      })
    }
  }
}
