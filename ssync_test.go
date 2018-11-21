package main

import (
  "os"
  "fmt"
  "strings"
  "testing"
  "io/ioutil"
  "path/filepath"
)

// check saved state equals result on both paths
func stateEqual(t *testing.T, path1, path2 string, result []string) string {
  for _, p := range []string{path1, path2} {
    r, _ := stringSliceFromFile(filepath.Join(p, ".ssync-test"))
    if strings.Join(r, "\n") != strings.Join(result, "\n") {
      return fmt.Sprintf("Expected %q, got %q", r, result)
    }
  }
  return ""
}

// check combined paths equal
func pathsEqual(t *testing.T, path1, path2 string, result []string) string {
  for _, p := range []string{path1, path2} {
    ps, err := stringSliceFromPathWalk(p)
    if err != nil {
      return err.Error()
    }

    if strings.Join(ps, "\n") != strings.Join(result, "\n") {
      return fmt.Sprintf("Expected %v, got %v", result, ps)
    }
  }
  return ""
}

func testExec(t *testing.T, i, path1, path2 string) {
  tmpFile(t, i, func(in *os.File){
    err := exec([]string{"test", path1, path2}, in)
    if err != nil {
      t.Fatal(err)
    }
  })
}

// test initial
func testExecInitial(t *testing.T, path1, path2 string) {
  testExec(t, "Y", path1, path2)

  result := []string{
    ".ssync-test", "dir1", "dir1/dir2", "dir1/dir2/file3", "dir1/file2",
    "dir3", "dir3/file5", "file1", "file4",
  }

  if e := pathsEqual(t, path1, path2, result); e != "" {
    t.Errorf(e)
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

  if e := stateEqual(t, path1, path2, result); e != "" {
    t.Errorf(e)
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

  if e := stateEqual(t, path1, path2, result); e != "" {
    t.Errorf(e)
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

  if e := pathsEqual(t, path1, path2, result); e != "" {
    t.Errorf(e)
  }
}

// test update sync
func testExecUpdate(t *testing.T, path1, path2 string) {
  writeFile(t, path1, &TestFile{"file1", "newFile1Contents", "2018-03-01"})

  testExec(t, "Y", path1, path2)

  c, _ := ioutil.ReadFile(filepath.Join(path2, "file1"))
  if string(c) != "newFile1Contents" {
    t.Errorf("Expected %v, got %v", "newFile1Contents", c)
  }
}

// test when both paths contain a few folder
// avoid merging, instead add (X) to 2nd path's folder name
func testExecCommonNewFolder(t *testing.T, path1, path2 string) {
  writeFile(t, path1, &TestFile{"dir99/dir22/file0", "FileContents", "2018-03-01"})
  writeFile(t, path2, &TestFile{"dir99/dir22/file1", "FileDifConts", "2018-02-01"})

  testExec(t, "Y", path1, path2)

  result := []string{
    ".ssync-test", "dir1", "dir1/file2", "dir3", "dir3/file5",
    "dir99", "dir99 (1)", "dir99 (1)/dir22", "dir99 (1)/dir22/file1",
    "dir99/dir22", "dir99/dir22/file0", "file1",
  }

  if e := pathsEqual(t, path1, path2, result); e != "" {
    t.Errorf(e)
  }

  if e := stateEqual(t, path1, path2, result); e != "" {
    t.Errorf(e)
  }
}

func TestExec(t *testing.T) {
  flagConfirm = true
  defer func() {
    flagConfirm = false
  }()

  path1, _ := createTestFiles(t, testFiles)
  defer os.RemoveAll(path1)

  path2, _ := createTestFiles(t, testFiles2)
  defer os.RemoveAll(path2)

  testExecInitial(t, path1, path2)
  testExecDeleteSkip(t, path1, path2)
  testExecDelete(t, path1, path2)
  testExecUpdate(t, path1, path2)
  testExecCommonNewFolder(t, path1, path2)
}
