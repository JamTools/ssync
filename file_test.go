package main

import (
  "os"
  "time"
  "strings"
  "testing"
  "io/ioutil"
  "path/filepath"
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
      tmpFile(t, k, func(in *os.File){
        r, _ := stringSliceFromFile(in.Name())
        if strings.Join(r, "\n") != strings.Join(v, "\n") {
          t.Errorf("Expected %v, got %v", v, r)
        }
      })
    }
  }
}

func TestStringSliceFromPathWalk(t *testing.T) {
  result := []string{
    "dir1",
    "dir1/dir2",
    "dir1/dir2/file3",
    "dir1/file2",
    "file1",
  }

  dir := createTestFiles(testFiles, t)
  defer os.RemoveAll(dir)

  paths, err := stringSliceFromPathWalk(dir)
  if err != nil {
    t.Fatal(err)
  }

  if strings.Join(paths, "\n") != strings.Join(result, "\n") {
    t.Errorf("Expected %v, got %v", result, paths)
  }
}

// TestDeleteConfirm & TestDelete also fulfill testing of pathsThatExist

func TestDeleteConfirm(t *testing.T) {
  dir := createTestFiles(testFiles, t)
  defer os.RemoveAll(dir)

  removes := []string{
    "file1",
    "dir1/dir2",
    "extra-does-not-exist-path",
  }

  delTests := map[string]bool{
    "Y": true,
    "N": false,
    "": false,
  }

  for input, result := range delTests {
    tmpFile(t, input, func(in *os.File){
      r := deleteConfirm(removes, dir, in)
      if r != result {
        t.Errorf("Expected %v, got %v", result, r)
      }
    })
  }
}

func TestDelete(t *testing.T){
  dir := createTestFiles(testFiles, t)
  defer os.RemoveAll(dir)

  removes := []string{
    "file1",
    "dir1/dir2",
  }

  delete(removes, dir)

  for _, v := range removes {
    fullpath := filepath.Join(dir, filepath.Join(strings.Split(v, "/")...))
    if _, err := os.Stat(fullpath); err == nil {
      t.Errorf("Expected '%v' to be deleted", v)
    }
  }
}

type testCopyAllFunc func(in, out string, ip, op []string)

func testCopyAll(t *testing.T, f testCopyAllFunc){
  srcPath := createTestFiles(testFiles, t)
  defer os.RemoveAll(srcPath)

  srcPaths, err := stringSliceFromPathWalk(srcPath)
  if err != nil {
    t.Fatal(err)
  }

  destPath, err := ioutil.TempDir("", "")
  if err != nil {
    t.Fatal(err)
  }
  defer os.RemoveAll(destPath)

  err = copyAll(srcPaths, srcPath, destPath)
  if err != nil {
    t.Fatal(err)
  }

  destPaths, err := stringSliceFromPathWalk(destPath)
  if err != nil {
    t.Fatal(err)
  }

  f(srcPath, destPath, srcPaths, destPaths)
}

func TestCopyAll(t *testing.T){
  testCopyAll(t, func(srcPath, destPath string, srcPaths, destPaths []string){
    // ensure all srcPaths equal destPaths
    if strings.Join(srcPaths, "\n") != strings.Join(destPaths, "\n") {
      t.Errorf("Expected %v, got %v", srcPaths, destPaths)
    }

    // ensure specified modified timestamp was set
    modTime, _ := time.Parse("2006-01-02", testFiles["file1"][1])
    destFullpath := filepath.Join(destPath, "file1")
    fi, _ := os.Stat(destFullpath)
    if fi.ModTime() != modTime {
      t.Errorf("Expected %v, got %v", modTime, fi.ModTime())
    }
  })
}

func TestMostRecentlyModified(t *testing.T){
  testCopyAll(t, func(srcPath, destPath string, srcPaths, destPaths []string){

    // ensure modified timestamp is preserved in copyFile
    a, b := mostRecentlyModified("file1", srcPath, destPath)
    if a != "" || b != "" {
      t.Errorf("Expected equal timestamps")
    }

    // ensure blank when checking directory
    a, b = mostRecentlyModified("dir1", srcPath, destPath)
    if a != "" || b != "" {
      t.Errorf("Expected blank timestamps for directory")
    }

    ct := time.Now().Local()

    // ensure when srcPath most recently modified
    srcFullpath := filepath.Join(srcPath, "dir1/file2")
    if err := os.Chtimes(srcFullpath, ct, ct); err != nil {
      t.Fatal(err)
    }

    a, b = mostRecentlyModified("dir1/file2", srcPath, destPath)
    if a != srcFullpath {
      t.Errorf("Expected %v, got %v", srcFullpath, a)
    }

    // ensure when destPath most recently modified
    destFullpath := filepath.Join(destPath, "dir1/dir2/file3")
    if err := os.Chtimes(destFullpath, ct, ct); err != nil {
      t.Fatal(err)
    }

    a, b = mostRecentlyModified("dir1/dir2/file3", srcPath, destPath)
    if a != destFullpath {
      t.Errorf("Expected %v, got %v", destFullpath, a)
    }
  })
}
