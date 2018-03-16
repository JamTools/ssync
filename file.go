package main

import (
  "os"
  "io"
  "fmt"
  "strings"
  "bufio"
  "path/filepath"
)

// read text file by newline into string slice
func stringSliceFromFile(path string) (lines []string, err error) {
  f, err := os.Open(path)
  if err != nil {
    return
  }
  defer f.Close()

  rd := bufio.NewReader(f)
  for {
    line, e := rd.ReadString('\n')
    if e != nil {
      if e == io.EOF {
        break
      }
      return lines, e
    }

    line = strings.TrimSpace(line)
    if len(line) > 0 {
      lines = append(lines, line)
    }
  }

  return
}

// recursive walk directory path creating string slice of child paths
func stringSliceFromPathWalk(p string) (paths []string, err error) {
  // closure to pass to filepath.Walk
  walkFunc := func(path string, f os.FileInfo, err error) error {
    // remove base path
    path = strings.Replace(path, p, "", 1)

    // if not blank (including separator)
    if len(path) > 1 {
      paths = append(paths, path[1:])
    }

    return err
  }

  err = filepath.Walk(p, walkFunc)
  return
}

// prompt confirmation before deleting files
func deleteConfirm(list []string, path string) bool {
  fmt.Printf("Simulate delete from '%s'...\n", path)

  count := 0
  for i := range list {
    fullpath := filepath.Join(path, list[i])

    _, err := os.Stat(fullpath)
    if err != nil {
      continue
    }

    count += 1
    fmt.Printf("%s\n", list[i])
  }
  fmt.Println()

  if count > 0 {
    fmt.Printf("Confirm delete files? (yes/no) ")
    return askConfirm()
  }

  return false
}

// remove all paths (dir & file)
func delete(list []string, path string) {
  fmt.Printf("Delete from '%s'...\n", path)

  for i := range list {
    fullpath := filepath.Join(path, list[i])

    fi, err := os.Stat(fullpath)
    if err != nil {
      continue
    }

    fmt.Printf("%s\n", list[i])

    if fi.IsDir() {
      os.RemoveAll(fullpath)
    } else {
      os.Remove(fullpath)
    }
  }
  fmt.Println()
}

// copy new files from srcPath to destPath
func create(paths []string, srcPath, destPath string) (err error) {
  for i := range paths {
    fi, err := os.Stat(filepath.Join(srcPath, paths[i]))
    if err != nil {
      // skip path if error while reading
      continue
    }

    if fi.IsDir() {
      err = os.MkdirAll(filepath.Join(destPath, paths[i]), 0777)
    } else {
      err = copyFile(
        filepath.Join(srcPath, paths[i]),
        filepath.Join(destPath, paths[i]),
      )
    }

    if err != nil {
      // return any error while writing
      return err
    }
  }
  return
}

// update common files if one is more recently updated
func update(paths []string, path1, path2 string) (err error) {
  for i := range paths {
    src, dest := mostRecentlyModified(paths[i], path1, path2)
    if len(src) > 0 && len(dest) > 0 {
      err := copyFile(src, dest)
      if err != nil {
        return err
      }
    }
  }
  return
}

// return file path if one is more recently modified
func mostRecentlyModified(file, path1, path2 string) (string, string) {
  src1 := filepath.Join(path1, file)
  src2 := filepath.Join(path2, file)

  fi1, err := os.Stat(src1)
  if err != nil || fi1.IsDir() {
    return "", ""
  }

  fi2, err := os.Stat(src2)
  if err != nil || fi2.IsDir() {
    return "", ""
  }

  // compared modified times
  if fi1.ModTime().Unix() > fi2.ModTime().Unix() {
    // update on path2
    return src1, src2
  } else if fi2.ModTime().Unix() > fi1.ModTime().Unix() {
    // update on path1
    return src2, src1
  }

  return "", ""
}

// copy file srcPath to destPath
func copyFile(srcPath, destPath string) (err error) {
  srcFile, err := os.Open(srcPath)
  if err != nil {
    return
  }
  defer srcFile.Close()

  destFile, err := os.Create(destPath)
  if err != nil {
    return
  }
  defer destFile.Close()

  _, err = io.Copy(destFile, srcFile)
  if err != nil {
    return
  }

  err = destFile.Sync()
  return
}
