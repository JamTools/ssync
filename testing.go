package main

import (
  "io"
  "os"
  "time"
  "strings"
  "testing"
  "io/ioutil"
  "path/filepath"
)

type TestFile struct {
  Name, Contents, Date string
}

var testFiles = []*TestFile{
  {".ssync-test", ".ssync-test\ndir1\ndir1/dir2\ndir1/dir2/file3\n", ""},
  {"file1", "file1Contents", "2018-01-01"},
  {"dir1/file2", "file2Contents", "2018-01-01"},
  {"dir1/dir2/file3", "file3Contents", "2018-01-01"},
}

var testFiles2 = []*TestFile{
  {".ssync-test", ".ssync-test\ndir1\ndir1/dir2\ndir1/dir2/file3\n", ""},
  {"file1", "file1Contents2", "2017-01-01"},
  {"file4", "file4Contents", "2018-01-01"},
  {"dir3/file5", "file2Contents", "2018-01-01"},
  {"dir1/dir2/file3", "file3Contents2", "2018-02-01"},
}

func createTestFiles(t *testing.T, files []*TestFile) (string, []string) {
  td, err := ioutil.TempDir("", "")
  if err != nil {
    t.Fatal(err)
  }

  paths := []string{}
  for i := range files {
    if len(files[i].Name) == 0 {
      continue
    }

    // set default date
    if len(files[i].Date) == 0 {
      files[i].Date = "2006-01-02"
    }

    writeFile(t, td, files[i])

    // append to paths
    paths = append(paths, filepath.Join(td, files[i].Name))
  }

  return td, paths
}

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

  _, _ = in.Seek(0, os.SEEK_SET)

  f(in)
}

// writes files also creating parent directories
func writeFile(t *testing.T, path string, testFile *TestFile) {
  pa := strings.Split(testFile.Name, "/")

  // return if filename blank
  if len(pa) == 0 || len(pa[len(pa)-1]) == 0 {
    return
  }

  // create parent dirs
  p := filepath.Join(path, filepath.Join(pa[:len(pa)-1]...))
  if len(pa) > 1 {
    err := os.MkdirAll(p, 0777)
    if err != nil {
      t.Fatal(err)
    }
  }

  // create file
  fullpath := filepath.Join(p, pa[len(pa)-1])
  err := ioutil.WriteFile(fullpath, []byte(testFile.Contents), 0644)
  if err != nil {
    t.Fatal(err)
  }

  // parse modified timestamp
  modTime, err := time.Parse("2006-01-02", testFile.Date)
  if err != nil {
    t.Fatal(err)
  }

  // set modified timestamp
  err = os.Chtimes(fullpath, modTime, modTime)
  if err != nil {
    t.Fatal(err)
  }
}
