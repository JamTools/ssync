package main

import (
  "os"
  "io/ioutil"
  "fmt"
  "strings"
  "sort"
  "log"
  "path/filepath"
)

type Args struct {
  Label string
  Paths []string
  In []string
  Out []string
}

// main package entry
func main() {
  args, cont := processFlags(nil)
  if !cont {
    os.Exit(1)
  }

  a := &Args{ Label: args[0], Paths: make([]string, 2),
    In: make([]string, 0), Out: make([]string, 0) }

  // check paths exist and are directories
  for x, i := range []int{1, 2} {
    path := filepath.Clean(args[i])

    fi, err := os.Stat(path)
    if err != nil || !fi.IsDir() {
      log.Fatalf("'%s' is not a directory", path)
    }

    a.Paths[x] = path
  }

  fmt.Printf("\nLabel: %s\nPath1: %s\nPath2: %s\n\n", 
    a.Label, a.Paths[0], a.Paths[1])

  // TODO: ensure state exists on both paths and is equal
  // this ensures it wasn't overwritten accidentally with another sync

  a.load()
  a.process()
  a.update()
  a.save()
}

// load previous paths from file: $path/.ssync-$label
// check both paths in case one accidentally deleted
func (a *Args) load() {
  for i := range a.Paths {
    a.In, _ = stringSliceFromFile(filepath.Join(a.Paths[i], ".ssync-" + a.Label))
    if len(a.In) > 0 {
      break
    }
  }
  sort.Strings(a.In)

  if flagVerbose {
    fmt.Printf("State: %v\n\n", a.In)
  }
}

// delete, copy new
func (a *Args) process() {
  for i := range a.Paths {
    paths, err := stringSliceFromPathWalk(a.Paths[i])
    if err != nil {
      log.Fatalf("%v", err)
    }
    sort.Strings(paths)

    if flagVerbose {
      fmt.Printf("%v\n\n", a.Paths[i])
    }

    // handle deleted files
    deleteList := notIn(a.In, paths)
    if len(deleteList) > 0 {
      del := true
      if flagConfirm {
        // ask to confirm deleted
        del = deleteConfirm(deleteList, a.Paths[1^i], nil)
      }

      if del {
        delete(deleteList, a.Paths[1^i])
      } else {
        // skip delete (add to a.Out to ask to confirm delete next time)
        a.Out = append(a.Out, notIn(a.Out, deleteList)...)
      }
    }

    // handle new files
    newList := notIn(paths, a.In)
    if len(newList) > 0 {
      err = copyAll(newList, a.Paths[i], a.Paths[1^i])
      if err != nil {
        log.Fatalf("%v", err)
      }
    }
  }
}

// update common files if one is more recently updated
func (a *Args) update() {
  for i := range a.In {
    src, dest := mostRecentlyModified(a.In[i], a.Paths[0], a.Paths[1])
    if len(src) > 0 && len(dest) > 0 {
      err := copyFile(src, dest)
      if err != nil {
        log.Fatalf("%v", err)
      }
    }
  }
}

// append and save updated paths to: $path/.ssync-$label
func (a *Args) save() {
  // append to a.Out
  for i := range a.Paths {
    paths, err := stringSliceFromPathWalk(a.Paths[i])
    if err != nil {
      log.Fatalf("%v", err)
    }

    a.Out = append(a.Out, notIn(paths, a.Out)...)
    sort.Strings(a.Out)
  }

  // write a.Out to file
  d1 := []byte(strings.Join(a.Out, "\n")+"\n")
  for i := range a.Paths {
    err := ioutil.WriteFile(filepath.Join(a.Paths[i], ".ssync-" + a.Label), d1, 0644)
    if err != nil {
      log.Fatalf("%v", err)
    }
  }
}
