package fsevents

import "testing"
import "github.com/sdegutis/assert"

import (
  "time"
  "os"
  "path/filepath"
)

func withCreate(action func(string)) {
  dummyfile := "dummyfile.txt"
  os.Create(dummyfile)

  action(dummyfile)

  os.Remove(dummyfile)
}

func TestFileChanges(t *testing.T) {
  ch := WatchPaths([]string{"."})

  withCreate(func(dummyfile string) {
    select {
    case <-ch:
    case <-time.After(time.Second * 2):
      t.Errorf("timed out")
    }
  })
}

func TestCanGetPath(t *testing.T) {
  ch := WatchPaths([]string{"."})

  withCreate(func(dummyfile string) {
    select {
    case events := <-ch:
      assert.Equals(t, len(events), 1)

      fullpath, _ := filepath.Abs(dummyfile)
      assert.Equals(t, events[0].Path, fullpath)
    case <-time.After(time.Second * 2):
      t.Errorf("timed out")
    }
  })
}

func TestOnlyWatchesSpecifiedPaths(t *testing.T) {
  ch := WatchPaths([]string{"imaginaryfile"})

  withCreate(func(dummyfile string) {
    select {
    case <-ch:
      t.Errorf("should have timed out, but got some file event")
    case <-time.After(time.Second * 1):
    }
  })
}
