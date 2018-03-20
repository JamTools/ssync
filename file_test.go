package main

import (
  "os"
  "io"
  "io/ioutil"
  "strings"
  "path/filepath"
  "testing"
)

func TestStringSliceFromFileNotExist(t *testing.T) {
  _, err := stringSliceFromFile(".ssync-this-file-does-not-exist")
  if err == nil {
    t.Errorf("Expected file not found error, got nil")
  }
}

func TestStringSliceFromFile(t *testing.T) {
  sliceTests := []map[string][]string{
    { "hello\nworld\n!\n": { "hello", "world", "!" } },
    { "\n\n  \n0 \n": { "0" } },
  }

  for i := range sliceTests {
    for k, v := range sliceTests[i] {
      in, err := ioutil.TempFile("", "")
      if err != nil {
        t.Fatal(err)
      }
      defer os.Remove(in.Name())
      defer in.Close()

      _, err = io.WriteString(in, k)
      if err != nil {
        t.Fatal(err)
      }

      r, _ := stringSliceFromFile(in.Name())
      if strings.Join(r, "\n") != strings.Join(v, "\n") {
        t.Errorf("Expected %v, got %v", v, r)
      }
    }
  }
}

func createTestFiles(paths map[string]string, t *testing.T) string {
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

    // created parent dirs
    if len(pa) > 1 {
      err := os.MkdirAll(path, 0777)
      if err != nil {
        t.Fatal(err)
      }
    }

    // create file
    if len(pa[len(pa)-1]) > 0 {
      err = ioutil.WriteFile(filepath.Join(path, pa[len(pa)-1]), []byte(c), 0644)
      if err != nil {
        t.Fatal(err)
      }
    }
  }

  return td
}

func TestStringSliceFromPathWalk(t *testing.T) {
  walkPaths := map[string]string{
    "file1": "file1Contents",
    "dir1/file2": "file2Contents",
    "dir1/dir2/file3": "file3Contents",
  }

  result := []string{
    "dir1",
    "dir1/dir2",
    "dir1/dir2/file3",
    "dir1/file2",
    "file1",
  }

  dir := createTestFiles(walkPaths, t)
  defer os.RemoveAll(dir)

  paths, err := stringSliceFromPathWalk(dir)
  if err != nil {
    t.Fatal(err)
  }

  if strings.Join(paths, "\n") != strings.Join(result, "\n") {
    t.Errorf("Expected %v, got %v", result, paths)
  }
}

// TODO: test copyFile modified timestamp equals