package main

import (
  "os"
  "io/ioutil"
  "fmt"
  "strings"
  "sort"
  "log"
  "flag"
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

func (a *Args) process() {
  fileName := ".ssync-" + a.Label

  // load previous paths from file: $path/.ssync-$label
  // check both paths in case one accidentally deleted
  prevList := make([]string, 0)
  for i := range a.Paths {
    prevList = loadList(filepath.Join(a.Paths[i], fileName))
    if len(prevList) > 0 {
      break
    }
  }
  sort.Strings(prevList)

  // slice for updated paths file to be saved at end
  outList := make([]string, 0)

  // DEBUG
  fmt.Printf("prev: %v\n\n", prevList)

  for i := range a.Paths {
    paths := pathList(a.Paths[i])
    otherPath := a.Paths[bool2int(!int2bool(i))]

    fmt.Printf("%v\n\n", a.Paths[i])

    deleteList := notIn(prevList, paths)
    if len(deleteList) > 0 {
      fmt.Printf("delete from %s : %v\n\n", otherPath, deleteList)

      del := true
      if flagConfirm {
        // ask to confirm deleted
        del = removeDeleted(deleteList, otherPath, true)
        fmt.Println()
      }

      if del {
        // actually delete
        _ = removeDeleted(deleteList, otherPath, false)
      } else {
        // skip delete (add to outList to ask to confirm delete next time)
        outList = append(outList, notIn(outList, deleteList)...)
      }
    }

    newList := notIn(paths, prevList)
    if len(newList) > 0 {
      fmt.Printf("new in %s : %v\n\n", otherPath, newList)

      copyNew(newList, a.Paths[i], otherPath)
    }
  }

  // update modified
  update(prevList, a.Paths[0], a.Paths[1])

  // append to outList
  for i := range a.Paths {
    outList = append(outList, notIn(pathList(a.Paths[i]), outList)...)
  }
  sort.Strings(outList)

  // write outList to file
  d1 := []byte(strings.Join(outList, "\n")+"\n")
  for i := range a.Paths {
    err := ioutil.WriteFile(filepath.Join(a.Paths[i], fileName), d1, 0644)
    if err != nil {
      log.Fatalf("%v", err)
    }
  }

  return
}
