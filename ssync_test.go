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

var testFiles = map[string][]string{
  "file1": { "file1Contents", "2018-01-01" },
  "dir1/file2": { "file2Contents", "2018-01-01" },
  "dir1/dir2/file3": { "file3Contents", "2018-01-01" },
}

var testFiles2 = map[string][]string{
  "file1": { "file1Contents2", "2017-01-01" },
  "file4": { "file4Contents", "2018-01-01" },
  "dir3/file5": { "file2Contents", "2018-01-01" },
  "dir1/dir2/file3": { "file3Contents2", "2018-02-01" },
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

func writeFile(t *testing.T, p, c, ts string) {
  err := ioutil.WriteFile(p, []byte(c), 0644)
  if err != nil {
    t.Fatal(err)
  }

  // parse modified timestamp
  modTime, err := time.Parse("2006-01-02", ts)
  if err != nil {
    t.Fatal(err)
  }

  // set modified timestamp
  err = os.Chtimes(p, modTime, modTime)
  if err != nil {
    t.Fatal(err)
  }
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
      writeFile(t, fullpath, c[0], c[1])
    }
  }

  return td
}

func testExec(t *testing.T, i, path1, path2 string) {
  tmpFile(t, i, func(in *os.File){
    err := exec([]string{"test", path1, path2}, in)
    if err != nil {
      t.Fatal(err)
    }
  })
}

// test initial new state
func testExecInitial(t *testing.T, path1, path2 string) {
  testExec(t, "Y", path1, path2)

  result := []string{
    ".ssync-test", "dir1", "dir1/dir2", "dir1/dir2/file3", "dir1/file2",
    "dir3", "dir3/file5", "file1", "file4",
  }

  // check combined paths equal
  for _, p := range []string{path1, path2} {
    ps, err := stringSliceFromPathWalk(p)
    if err != nil {
      t.Fatal(err)
    }

    if strings.Join(ps, "\n") != strings.Join(result, "\n") {
      t.Errorf("Expected %v, got %v", result, ps)
    }
  }

  // check file1 overwritten on path2
  c, _ := ioutil.ReadFile(filepath.Join(path2, "file1"))
  if string(c) != "file1Contents" {
    t.Errorf("Expected %v, got %v", "file1Contents", string(c))
  }

  // check dir1/dir2/file3 overwritten on path1
  c, _ = ioutil.ReadFile(filepath.Join(path1, "dir1", "dir2", "file3"))
  if string(c) != "file3Contents2" {
    t.Errorf("Expected %v, got %v", "file3Contents2", string(c))
  }

  // check saved state equals result on both paths
  for _, p := range []string{path1, path2} {
    r, _ := stringSliceFromFile(filepath.Join(p, ".ssync-test"))
    if strings.Join(r, "\n") != strings.Join(result[1:], "\n") {
      t.Errorf("Expected %v, got %v", result, r)
    }
  }
}

// test delete skip sync
func testExecDeleteSkip(t *testing.T, path1, path2 string) {
  os.RemoveAll(filepath.Join(path1, "dir1", "dir2"))
  os.Remove(filepath.Join(path2, "file4"))

  testExec(t, "N", path1, path2)

  delTests := map[string][]string{
    path1: []string{ ".ssync-test", "dir1", "dir1/file2", "dir3",
      "dir3/file5", "file1", "file4" },
    path2: []string{ ".ssync-test", "dir1", "dir1/dir2", "dir1/dir2/file3",
      "dir1/file2", "dir3", "dir3/file5", "file1" },
  }

  // check delete was skipped
  for k, v := range delTests {
    ps, err := stringSliceFromPathWalk(k)
    if err != nil {
      t.Fatal(err)
    }

    if strings.Join(ps, "\n") != strings.Join(v, "\n") {
      t.Errorf("Expected %v, got %v", v, ps)
    }
  }

  // check saved state still includes all files
  result := []string{ ".ssync-test", "dir1", "dir1/dir2", "dir1/dir2/file3",
    "dir1/file2", "dir3", "dir3/file5", "file1", "file4" }

  for _, p := range []string{path1, path2} {
    r, _ := stringSliceFromFile(filepath.Join(p, ".ssync-test"))
    if strings.Join(r, "\n") != strings.Join(result, "\n") {
      t.Errorf("Expected %v, got %v", result, r)
    }
  }
}

// test delete sync
func testExecDelete(t *testing.T, path1, path2 string) {
  os.RemoveAll(filepath.Join(path1, "dir1", "dir2"))
  os.Remove(filepath.Join(path2, "file4"))

  testExec(t, "Y", path1, path2)

  result := []string{
    ".ssync-test", "dir1", "dir1/file2", "dir3", "dir3/file5", "file1",
  }

  // check combined paths equal
  for _, p := range []string{path1, path2} {
    ps, err := stringSliceFromPathWalk(p)
    if err != nil {
      t.Fatal(err)
    }

    if strings.Join(ps, "\n") != strings.Join(result, "\n") {
      t.Errorf("Expected %v, got %v", result, ps)
    }
  }
}

// test update sync
func testExecUpdate(t *testing.T, path1, path2 string) {
  writeFile(t, filepath.Join(path1, "file1"), "newFile1Contents", "2018-03-01")

  testExec(t, "Y", path1, path2)

  c, _ := ioutil.ReadFile(filepath.Join(path2, "file1"))
  if string(c) != "newFile1Contents" {
    t.Errorf("Expected %v, got %v", "newFile1Contents", c)
  }
}

func TestExec(t *testing.T) {
  flagConfirm = true
  defer func() {
    flagConfirm = false
  }()

  path1 := createTestFiles(testFiles, t)
  defer os.RemoveAll(path1)

  path2 := createTestFiles(testFiles2, t)
  defer os.RemoveAll(path2)

  testExecInitial(t, path1, path2)
  testExecDeleteSkip(t, path1, path2)
  testExecDelete(t, path1, path2)
  testExecUpdate(t, path1, path2)
}
