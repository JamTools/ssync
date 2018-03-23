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
  args, cont := processFlags()
  if !cont {
    os.Exit(1)
  }

  if err := exec(args, os.Stdin); err != nil {
    log.Fatalf("%v", err)
  }
}

// main process
func exec(args []string, in *os.File) error {
  a := &Args{ Label: args[0], Paths: []string{args[1], args[2]},
    In: make([]string, 0), Out: make([]string, 0) }

  // check paths exist and are directories
  for x := range a.Paths {
    p, err := checkDir(a.Paths[x])
    if err != nil {
      return err
    }
    a.Paths[x] = p
  }

  fmt.Printf("\nLabel: %s\nPath1: %s\nPath2: %s\n", 
    a.Label, a.Paths[0], a.Paths[1])

  err := a.loadState()
  if err != nil {
    return err
  }

  err = a.process(in)
  if err != nil {
    return err
  }

  err = a.saveState()
  if err != nil {
    return err
  }

  fmt.Printf("\nssync completed\n")
  return nil
}

// load previous state
func (a *Args) loadState() error {
  l := ".ssync-" + a.Label
  s1, _ := stringSliceFromFile(filepath.Join(a.Paths[0], l))
  s2, _ := stringSliceFromFile(filepath.Join(a.Paths[1], l))

  // new state
  if len(s1) == 0 && len(s2) == 0 {
    return nil
  }

  // ensure states are equal (avoid accidental overwrite with another sync)
  if strings.Join(s1, "\n") == strings.Join(s2, "\n") {
    a.In = s1
    return nil
  }

  // handle unequal states
  return fmt.Errorf("shared state unequal")
}

// delete, copy new, update
func (a *Args) process(in *os.File) error {
  for i := range a.Paths {
    paths, err := stringSliceFromPathWalk(a.Paths[i])
    if err != nil {
      return err
    }

    fmt.Printf("\nProcessing: %v\n", a.Paths[i])

    // handle deleted files
    if len(a.In) > 0 {
      deleteList := notIn(a.In, paths)
      if len(deleteList) > 0 {
        del := true
        if flagConfirm {
          // ask to confirm deleted
          del = deleteConfirm(deleteList, a.Paths[1^i], in)
        }

        if del {
          delete(deleteList, a.Paths[1^i])
        } else {
          // skip delete (add to a.Out to ask to confirm delete next time)
          a.Out = append(a.Out, notIn(a.Out, deleteList)...)
        }
      }
    }

    // handle new files
    newList := notIn(paths, a.In)
    if len(newList) > 0 {
      err = copyAll(newList, a.Paths[i], a.Paths[1^i])
      if err != nil {
        return err
      }
    }
  }

  // update common files if one is more recently updated
  for i := range a.In {
    src, dest := mostRecentlyModified(a.In[i], a.Paths[0], a.Paths[1])
    if len(src) > 0 && len(dest) > 0 {
      err := copyFile(src, dest)
      if err != nil {
        return err
      }
    }
  }

  return nil
}

// append and save updated paths to: $path/.ssync-$label
func (a *Args) saveState() error {
  // append to a.Out
  for i := range a.Paths {
    paths, err := stringSliceFromPathWalk(a.Paths[i])
    if err != nil {
      return err
    }

    a.Out = append(a.Out, notIn(paths, a.Out)...)
    sort.Strings(a.Out)
  }

  // write a.Out to file
  l := ".ssync-" + a.Label
  d1 := []byte(strings.Join(a.Out, "\n")+"\n")
  for i := range a.Paths {
    err := ioutil.WriteFile(filepath.Join(a.Paths[i], l), d1, 0644)
    if err != nil {
      return err
    }
  }

  return nil
}
