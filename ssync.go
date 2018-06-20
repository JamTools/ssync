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

  fmt.Printf("\nssync finished.\n")
  return nil
}

// load previous state
func (a *Args) loadState() error {
  l := ".ssync-" + a.Label
  s1, _ := stringSliceFromFile(filepath.Join(a.Paths[0], l))
  s2, _ := stringSliceFromFile(filepath.Join(a.Paths[1], l))

  // ensure states are equal
  // error indicates accidental overwrite with another sync
  if strings.Join(s1, "\n") != strings.Join(s2, "\n") {
    return fmt.Errorf("shared state unequal")
  }

  a.In = s1
  return nil
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
      fmt.Printf("\nCopy new files to: %v\n", a.Paths[1^i])
      err = copyAll(newList, a.Paths[i], a.Paths[1^i])
      if err != nil {
        return err
      }
    }
  }

  // update common files if one is more recently updated
  fmt.Printf("\nUpdate modified:\n")
  for i := range a.In {
    src, dest, found := mostRecentlyModified(a.In[i], a.Paths[0], a.Paths[1])
    if found && len(src) > 0 && len(dest) > 0 {
      err := copyFile(a.In[i], src, dest)
      if err != nil {
        return err
      }
    }
  }

  return nil
}

// save shared state to file
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

  // write a.Out state to file
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
