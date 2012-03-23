package fsevents

import "testing"
import "github.com/sdegutis/assert"

import (
  "time"
  "os"
  "path/filepath"
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

func TestCanGetPath(t *testing.T) {
  ch := WatchPaths([]string{"."})

  dummyfile := "dummyfile.txt"

  os.Create(dummyfile)

  select {
  case events := <-ch:
    assert.Equals(t, len(events), 1)

    fullpath, _ := filepath.Abs(dummyfile)
    assert.Equals(t, events[0].Path, fullpath)
  case <-time.After(time.Second * 2):
    t.Errorf("timed out")
  }

  os.Remove(dummyfile)
}

func TestOnlyWatchesSpecifiedPaths(t *testing.T) {
  ch := WatchPaths([]string{"imaginaryfile"})

  dummyfile := "dummyfile.txt"

  os.Create(dummyfile)

  select {
  case <-ch:
    t.Errorf("should have timed out, but got some file event")
  case <-time.After(time.Second * 1):
  }

  os.Remove(dummyfile)
}
