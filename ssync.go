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
  newList := make([][]string, 2)

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

    // save new files
    newList[i] = notIn(paths, a.In)
  }

  // update common files if one is more recently updated
  if len(a.In) > 0 {
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
  }

  // rename folder where new folder exists on both sides
  newList, err := a.commonFolders(newList)
  if err != nil {
    return err
  }

  // copy new files
  for i := range newList {
    if len(newList[i]) > 0 {
      fmt.Printf("\nCopy new files to: %v\n", a.Paths[1^i])
      err := copyAll(newList[i], a.Paths[i], a.Paths[1^i])
      if err != nil {
        return err
      }
    }
  }

  return nil
}

// renames folder to (1) in the instance where new folder exists on both sides
func (a *Args) commonFolders(files [][]string) ([][]string, error) {
  m := make(map[string]bool)
  for i := range files[0] {
    _, err := checkDir(filepath.Join(a.Paths[0], files[0][i]))
    if err == nil {
      m[files[0][i]] = true
    }
  }

  var prev, prevRep string
  for i := range files[1] {
    if prev != "" && prevRep != "" {
      // replace new folder name if exists in path
      files[1][i] = strings.Replace(files[1][i], prev, prevRep, 1)
    }

    f := files[1][i]
    _, err := checkDir(filepath.Join(a.Paths[1], f))
    if err != nil {
      continue
    }

    _, ok := m[f]
    if ok && (prev == "" || strings.Index(f, prev) != 0) {
      // if a folder and not part of previous path
      // rename folder appending (X), incrementing X until does not exist
      fp := filepath.Join(a.Paths[1], f)
      newFolder, err := RenameFolder(fp, fp)
      if err != nil {
        return files, err
      }

      // update file name in slice
      files[1][i] = strings.Replace(newFolder, a.Paths[1], "", 1)[1:]

      // save updated name to replace within subpaths
      prev = f
      prevRep = files[1][i]
    }
  }

  return files, nil
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
