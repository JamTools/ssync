package main

import (
  "os"
  "io"
  "sort"
  "strings"
  "log"
  "bufio"
  "path/filepath"
)

// read file into string slice
func loadList(path string) (result []string) {
  fi, err := os.Stat(path)
  if err != nil || fi.IsDir() {
    return
  }

  f, err := os.Open(path)
  if err != nil {
    return
  }
  defer f.Close()

  rd := bufio.NewReader(f)
  for {
    line, err := rd.ReadString('\n')
    if err != nil {
      if err == io.EOF {
        break
      }

      log.Fatalf("read file line error: %v", err)
    }

    line = strings.TrimSpace(line)
    if line != "" {
      result = append(result, line)
    }
  }

  return
}

// recursive step through directory creating string slice
func pathList(p string) []string {
  paths := make([]string, 0)

  e := filepath.Walk(p, func(path string, f os.FileInfo, err error) error {
    // remove base path
    path = strings.Replace(path, p, "", 1)

    // skip blank / remove separator
    if len(path) > 1 {
      paths = append(paths, path[1:])
    }

    return err
  })

  if e != nil {
    log.Fatalf("%v", e)
  }

  sort.Strings(paths)
  return paths
}
