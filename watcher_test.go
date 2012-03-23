package fsevents

import "testing"

import (
  "time"
  "os"
)

func TestFileChanges(t *testing.T) {
  ch := WatchPaths([]string{"."})

  dummyfile := "dummyfile.txt"

  os.Create(dummyfile)

  select {
  case <-ch:
  case <-time.After(time.Second * 2):
    t.Errorf("timed out")
  }

  os.Remove(dummyfile)
}
