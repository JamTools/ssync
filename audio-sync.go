package main

import (
  "os"
  "io"
  "io/ioutil"
  "fmt"
  "strings"
  "sort"
  "log"
  "flag"
  "bufio"
  "path/filepath"
)

type Args struct {
  Label string
  Paths []string
}

// main package entry
func main() {
  a := &Args{ Label: flag.Arg(0), Paths: make([]string, 2) }

  // check paths exist and are directories
  for x, i := range []int{1, 2} {
    path := filepath.Clean(flag.Arg(i))

    fi, err := os.Stat(path)
    if err != nil || !fi.IsDir() {
      log.Fatalf("'%s' is not a directory", path)
    }

    a.Paths[x] = path
  }

  fmt.Printf("\nLabel: %s\nPath1: %s\nPath2: %s\n\n", 
    a.Label, a.Paths[0], a.Paths[1])

  a.process()
}

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

    // skip blank & hidden files
    if len(path) > 1 && path[1:2] != "." {
      paths = append(paths, path[1:])
    }

    return err
  })

  if e != nil {
    log.Fatalf("%v", e)
  }

  return paths
}

// remove what was deleted from one path from the other path
func removeDeleted(list []string, path string) {
  fmt.Printf("delete: %v\n\n", list)
}

func copyNew(list []string, path string) {
  fmt.Printf("new: %v\n\n", list)
}

func (a *Args) process() {
  fileName := ".audio-sync-" + a.Label

  // load previous paths from file: $path/.audio-sync-$label
  // check both paths in case one accidentally deleted
  prev := make([]string, 0)
  for i := range a.Paths {
    prev = loadList(filepath.Join(a.Paths[i], fileName))
    if len(prev) > 0 {
      break
    }
  }
  sort.Strings(prev)
  
  fmt.Printf("prev: %v\n\n", prev)

  for i := range a.Paths {
    fileList := pathList(a.Paths[i])
    otherPath := a.Paths[bool2int(!int2bool(i))]

    fmt.Printf("%s\n\n", a.Paths[i])

    deleteList := notIn(prev, fileList)
    removeDeleted(deleteList, otherPath)

    newList := notIn(fileList, prev)
    copyNew(newList, otherPath)
  }

  // write new combined path list file
  fileList := pathList(a.Paths[0])
  d1 := []byte(strings.Join(fileList, "\n")+"\n")
  for i := range a.Paths {
    err := ioutil.WriteFile(filepath.Join(a.Paths[i], fileName), d1, 0644)
    if err != nil {
      log.Fatalf("%v", err)
    }
  }

  return
}
