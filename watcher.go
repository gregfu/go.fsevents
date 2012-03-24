package fsevents

/*
#cgo LDFLAGS: -framework CoreServices
#include <CoreServices/CoreServices.h>
FSEventStreamRef fswatch_stream_for_paths(char** paths, int paths_n);
*/
import "C"
import "unsafe"

const (
  FlagUnmount = uint32(C.kFSEventStreamEventFlagUnmount)
  FlagMount   = uint32(C.kFSEventStreamEventFlagMount)

  FlagItemCreated       = uint32(C.kFSEventStreamEventFlagItemCreated)
  FlagItemRemoved       = uint32(C.kFSEventStreamEventFlagItemRemoved)
  FlagItemInodeMetaMod  = uint32(C.kFSEventStreamEventFlagItemInodeMetaMod)
  FlagItemRenamed       = uint32(C.kFSEventStreamEventFlagItemRenamed)
  FlagItemModified      = uint32(C.kFSEventStreamEventFlagItemModified)
  FlagItemFinderInfoMod = uint32(C.kFSEventStreamEventFlagItemFinderInfoMod)
  FlagItemChangeOwner   = uint32(C.kFSEventStreamEventFlagItemChangeOwner)
  FlagItemXattrMod      = uint32(C.kFSEventStreamEventFlagItemXattrMod)

  FlagItemIsFile    = uint32(C.kFSEventStreamEventFlagItemIsFile)
  FlagItemIsDir     = uint32(C.kFSEventStreamEventFlagItemIsDir)
  FlagItemIsSymlink = uint32(C.kFSEventStreamEventFlagItemIsSymlink)
)

type watchingInfo struct {
  channel chan []PathEvent
  runloop C.CFRunLoopRef
}

var watchers = make(map[C.FSEventStreamRef]watchingInfo)

type PathEvent struct {
  Path  string
  Flags uint32
}

func Unwatch(ch chan []PathEvent) {
  for stream, info := range watchers {
    if ch == info.channel {
      C.FSEventStreamStop(stream)
      C.FSEventStreamInvalidate(stream)
      C.FSEventStreamRelease(stream)
      C.CFRunLoopStop(info.runloop)
    }
  }
}

func WatchPaths(paths []string) chan []PathEvent {
  type watchSuccessData struct {
    runloop C.CFRunLoopRef
    stream  C.FSEventStreamRef
  }

  successChan := make(chan *watchSuccessData)

  go func() {
    var cpaths []*C.char
    for _, dir := range paths {
      path := C.CString(dir)
      defer C.free(unsafe.Pointer(path))
      cpaths = append(cpaths, path)
    }

    stream := C.fswatch_stream_for_paths(&cpaths[0], C.int(len(cpaths)))

    ok := C.FSEventStreamStart(stream) != 0
    if ok {
      successChan <- &watchSuccessData{
        runloop: C.CFRunLoopGetCurrent(),
        stream:  stream,
      }
      C.CFRunLoopRun()
    } else {
      successChan <- nil
    }
  }()

  watchingData := <-successChan

  if watchingData == nil {
    return nil
  }

  newChan := make(chan []PathEvent)
  watchers[watchingData.stream] = watchingInfo{
    channel: newChan,
    runloop: watchingData.runloop,
  }
  return newChan
}

//export watchDirsCallback
func watchDirsCallback(stream C.FSEventStreamRef, count C.size_t, paths **C.char, flags *C.FSEventStreamEventFlags) {
  var events []PathEvent

  for i := 0; i < int(count); i++ {
    cpaths := uintptr(unsafe.Pointer(paths)) + (uintptr(i) * unsafe.Sizeof(*paths))
    cpath := *(**C.char)(unsafe.Pointer(cpaths))

    cflags := uintptr(unsafe.Pointer(flags)) + (uintptr(i) * unsafe.Sizeof(*flags))
    cflag := *(*C.FSEventStreamEventFlags)(unsafe.Pointer(cflags))

    events = append(events, PathEvent{
      Path:  C.GoString(cpath),
      Flags: uint32(cflag),
    })
  }

  watchers[stream].channel <- events
}
