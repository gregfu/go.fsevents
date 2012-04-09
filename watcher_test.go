package fsevents

import "testing"
import "github.com/sdegutis/go.assert"

import (
  "os"
  "path/filepath"
  "time"
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
    case <-time.After(time.Second * 1):
      t.Errorf("should have got some file event, but timed out")
    }
  })
}

func TestEventFlags(t *testing.T) {
  ch := WatchPaths([]string{"."})

  withCreate(func(dummyfile string) {
    select {
    case events := <-ch:
      assert.Equals(t, len(events), 1)
      assert.True(t, events[0].Flags&FlagItemCreated != 0)
    case <-time.After(time.Second * 1):
      t.Errorf("should have got some file event, but timed out")
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

func TestCanUnwatch(t *testing.T) {
  ch := WatchPaths([]string{"."})

  Unwatch(ch)

  withCreate(func(dummyfile string) {
    select {
    case <-ch:
      t.Errorf("should have timed out, but got some file event")
    case <-time.After(time.Second * 1):
    }
  })
}
