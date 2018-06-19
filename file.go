package main

import (
  "os"
  "io"
  "fmt"
  "sort"
  "bufio"
  "strings"
  "path/filepath"
)

// check if path is directory
func checkDir(path string) (string, error) {
  path = filepath.Clean(path)

  fi, err := os.Stat(path)
  if err != nil || !fi.IsDir() {
    return path, fmt.Errorf("'%s' is not a directory", path)
  }

  return path, nil
}

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

  sort.Strings(lines)
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
  sort.Strings(paths)
  return
}

// check if path exists & exec pathFunction for each iteration
type pathFunction func(fi os.FileInfo, path string)

func pathsThatExist(list []string, path string, f pathFunction) int {
  count := 0

  for i := range list {
    fullpath := filepath.Join(path, list[i])

    fi, err := os.Stat(fullpath)
    if err != nil {
      continue
    }

    count += 1
    fmt.Printf("%s\n", list[i])

    if f != nil {
      f(fi, fullpath)
    }
  }

  return count
}

// prompt confirmation before deleting files
func deleteConfirm(list []string, path string, in *os.File) bool {
  fmt.Printf("\nSimulate delete from '%s'...\n", path)

  result := false
  count := pathsThatExist(list, path, nil)

  if count > 0 {
    fmt.Printf("\nConfirm delete files? (yes/no) ")
    result = askConfirm(in)
  }

  return result
}

// remove all paths (dir & file)
func delete(list []string, path string) {
  fmt.Printf("\nDelete from '%s'...\n", path)

  _ = pathsThatExist(list, path, func(fi os.FileInfo, fullpath string) {
    if fi.IsDir() {
      os.RemoveAll(fullpath)
    } else {
      os.Remove(fullpath)
    }
  })
}

// copy new files & folders from srcPath to destPath
func copyAll(paths []string, srcPath, destPath string) (err error) {
  for i := range paths {
    // skip path if error while reading
    fi, err := os.Stat(filepath.Join(srcPath, paths[i]))
    if err != nil {
      continue
    }

    if fi.IsDir() {
      err = os.MkdirAll(filepath.Join(destPath, paths[i]), 0777)
    } else {
      src, dest := mostRecentlyModified(paths[i], srcPath, destPath)

      // only copy if dne or more recently modified
      if src == srcPath || (len(src) == 0 && len(dest) == 0) {
        err = copyFile(paths[i], srcPath, destPath)
      }
    }

    if err != nil {
      // return any error while writing
      return err
    }
  }
  return
}

// return file path if one is more recently modified
func mostRecentlyModified(file, path1, path2 string) (string, string) {
  fi1, err := os.Stat(filepath.Join(path1, file))
  if err != nil || fi1.IsDir() {
    return "", ""
  }

  fi2, err := os.Stat(filepath.Join(path2, file))
  if err != nil || fi2.IsDir() {
    return "", ""
  }

  // override flag option
  switch flagForcePath {
  case 1:
    return path1, path2
  case 2:
    return path2, path1
  }

  // compared modified times
  if fi1.ModTime().Unix() > fi2.ModTime().Unix() {
    // update on path2
    return path1, path2
  } else if fi2.ModTime().Unix() > fi1.ModTime().Unix() {
    // update on path1
    return path2, path1
  }

  // equal, do nothing
  return "", ""
}

// copy file srcPath to destPath
func copyFile(file, srcPath, destPath string) (err error) {
  srcFile, err := os.Open(filepath.Join(srcPath, file))
  if err != nil {
    return
  }
  defer srcFile.Close()

  srcInfo, err := srcFile.Stat()
  if err != nil {
    return
  }

  fmt.Printf("%s\n", file)

  destFullpath := filepath.Join(destPath, file)

  destFile, err := os.Create(destFullpath)
  if err != nil {
    return
  }
  defer destFile.Close()

  _, err = io.Copy(destFile, srcFile)
  if err != nil {
    return
  }

  err = destFile.Sync()
  if err != nil {
    return
  }

  err = os.Chtimes(destFile.Name(), srcInfo.ModTime(), srcInfo.ModTime())
  return
}
