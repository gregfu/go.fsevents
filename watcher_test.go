package fsevents

import "testing"

import (
  "time"
  "fmt"
  "os"
)

func TestFileChanges(t *testing.T) {
  ch := WatchPaths([]string{"."})

  dummyfile := "dummyfile.txt"

  os.Create(dummyfile)

  select {
  case <-ch:
    fmt.Println("woot!")
  case <-time.After(time.Second * 2):
    fmt.Println("aww")
    t.Errorf("timed out")
  }

  os.Remove(dummyfile)
}
