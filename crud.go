package main

import (
  "io"
  "os"
  "fmt"
  "log"
  "path/filepath"
)

// remove what was deleted from one path from the other path
func removeDeleted(list []string, path string, confirm bool) bool {
  if confirm {
    fmt.Printf("Simulate delete from '%s'...\n", path)
  } else {
    fmt.Printf("Delete from '%s'...\n", path)
  }

  count := 0
  for i := range list {
    fullpath := filepath.Join(path, list[i])

    fi, err := os.Stat(fullpath)
    if err != nil {
      continue
    }

    count += 1
    fmt.Printf("%s\n", list[i])

    if !confirm {
      if fi.IsDir() {
        os.RemoveAll(fullpath)
      } else {
        os.Remove(fullpath)
      }
    }
  }
  fmt.Println()

  if confirm && count > 0 {
    fmt.Printf("Confirm delete files? (yes/no) ")
    return askConfirm()
  }

  return true
}

func copyNew(list []string, srcPath, destPath string) {
  for i := range list {
    fi, err := os.Stat(filepath.Join(srcPath, list[i]))
    if err != nil {
      continue
    }

    if fi.IsDir() {
      err = os.MkdirAll(filepath.Join(destPath, list[i]), 0777)
      if err != nil {
        log.Fatalf("%v", err)
      }
    } else {
      copyFile(
        filepath.Join(srcPath, list[i]),
        filepath.Join(destPath, list[i]),
      )
    }
  }

}

func update(list []string, path1, path2 string) {
  for i := range list {
    src, dest := mostRecentlyModified(list[i], path1, path2)
    if len(src) > 0 && len(dest) > 0 {
      copyFile(src, dest)
    }
  }
}

func mostRecentlyModified(file, path1, path2 string) (string, string) {
  src1 := filepath.Join(path1, file)
  src2 := filepath.Join(path2, file)

  fi1, err := os.Stat(src1)
  if err != nil {
    return "", ""
  }
  fi2, err := os.Stat(src2)
  if err != nil {
    return "", ""
  }

  // if not directory use most recently modified
  if !fi1.IsDir() && !fi2.IsDir() {
    if fi1.ModTime().Unix() > fi2.ModTime().Unix() {
      // update on path2
      return src1, src2
    } else if fi2.ModTime().Unix() > fi1.ModTime().Unix() {
      // update on path1
      return src2, src1
    }
  }

  return "", ""
}

func copyFile(srcPath, destPath string) {
  srcFile, err := os.Open(srcPath)
  if err != nil {
    log.Fatalf("%v", err)
  }
  defer srcFile.Close()

  destFile, err := os.Create(destPath)
  if err != nil {
    log.Fatalf("%v", err)
  }
  defer destFile.Close()

  _, err = io.Copy(destFile, srcFile)
  if err != nil {
    log.Fatalf("%v", err)
  }

  err = destFile.Sync()
  if err != nil {
    log.Fatalf("%v", err)
  }
}
