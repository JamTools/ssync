package main

import (
  "os"
  "io"
  "time"
  "strings"
  "testing"
  "io/ioutil"
  "path/filepath"
)

func tmpFile(t *testing.T, input string, f func(in *os.File)) {
  in, err := ioutil.TempFile("", "")
  if err != nil {
    t.Fatal(err)
  }
  defer os.Remove(in.Name())
  defer in.Close()

  _, err = io.WriteString(in, input)
  if err != nil {
    t.Fatal(err)
  }

  _, err = in.Seek(0, os.SEEK_SET)
  if err != nil {
    t.Fatal(err)
  }

  f(in)
}

var testFiles = map[string][]string{
  "file1": { "file1Contents", "2018-01-01" },
  "dir1/file2": { "file2Contents", "2018-01-01" },
  "dir1/dir2/file3": { "file3Contents", "2018-01-01" },
}

func createTestFiles(paths map[string][]string, t *testing.T) string {
  td, err := ioutil.TempDir("", "")
  if err != nil {
    t.Fatal(err)
  }

  for p, c := range paths {
    pa := strings.Split(p, "/")

    if len(pa) == 0 {
      continue
    }

    path := filepath.Join(td, filepath.Join(pa[:len(pa)-1]...))

    // create parent dirs
    if len(pa) > 1 {
      err := os.MkdirAll(path, 0777)
      if err != nil {
        t.Fatal(err)
      }
    }

    // create file
    if len(pa[len(pa)-1]) > 0 {
      fullpath := filepath.Join(path, pa[len(pa)-1])

      err = ioutil.WriteFile(fullpath, []byte(c[0]), 0644)
      if err != nil {
        t.Fatal(err)
      }

      // parse modified timestamp
      modTime, err := time.Parse("2006-01-02", c[1])
      if err != nil {
        t.Fatal(err)
      }
      // set modified timestamp
      err = os.Chtimes(fullpath, modTime, modTime)
      if err != nil {
        t.Fatal(err)
      }
    }
  }

  return td
}
