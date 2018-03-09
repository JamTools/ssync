package main

import (
  "sort"
)

func bool2int(b bool) int {
  if b {
    return 1
  }
  return 0
}

func int2bool(i int) bool {
  if i == 0 {
    return false
  }
  return true
}

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
