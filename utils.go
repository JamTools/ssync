package main

import (
  "io"
  "os"
  "fmt"
  "sort"
  "log"
)

func notIn(a, b []string) []string {
  list := make([]string, 0)

  for i := range a {
    x := sort.Search(len(b), func(x int) bool { return b[x] >= a[i] })

    if len(b) == 0 || x >= len(b) || (x < len(b) && b[x] != a[i]) {
      list = append(list, a[i])
    }
  }

  return list
}

func askConfirm(in *os.File) bool {
  yesResponses := []string{"y", "Y", "yes", "Yes", "YES"}
  noResponses := []string{"n", "N", "no", "No", "NO"}

  // closure to check response
  posString := func (slice []string, element string) int {
    for index, elem := range slice {
      if elem == element {
        return index
      }
    }
    return -1
  }

  var response string
  _, err := fmt.Fscan(in, &response)
  if err != nil && err != io.EOF {
    log.Fatal(err)
  }

  if posString(yesResponses, response) >= 0 {
    return true
  } else if posString(noResponses, response) >= 0 {
    return false
  }

  return false
}
