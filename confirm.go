package main

import (
  "fmt"
  "log"
)

func askConfirm() bool {
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
  _, err := fmt.Scanln(&response)
  if err != nil {
    log.Fatal(err)
  }

  if posString(yesResponses, response) >= 0 {
    return true
  } else if posString(noResponses, response) >= 0 {
    return false
  }

  return false
}
