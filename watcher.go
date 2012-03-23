package fsevents

/*
#cgo LDFLAGS: -framework CoreServices
#include <CoreServices/CoreServices.h>
FSEventStreamRef fswatch_stream_for_paths(char** paths, int paths_n);
void fswatch_unwatch_stream(FSEventStreamRef stream);
*/
import "C"
import "unsafe"

const (
  FlagItemCreated = uint32(C.kFSEventStreamEventFlagItemCreated)
)

type watchingInfo struct {
  channel chan []PathEvent
  runloop C.CFRunLoopRef
}

var watchers = make(map[C.FSEventStreamRef]watchingInfo)

type PathEvent struct {
  Path string
  Flags uint32
}

func Unwatch(ch chan []PathEvent) {
  for stream, info := range watchers {
    if ch == info.channel {
      C.fswatch_unwatch_stream(stream)
      C.CFRunLoopStop(info.runloop)
    }
  }
}

func WatchPaths(paths []string) chan []PathEvent {
  type watchSuccessData struct{
    runloop C.CFRunLoopRef
    stream C.FSEventStreamRef
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
        stream: stream,
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
    path := C.GoString(cpath)

    cflags := uintptr(unsafe.Pointer(flags)) + (uintptr(i) * unsafe.Sizeof(*flags))
    cflag := *(*C.FSEventStreamEventFlags)(unsafe.Pointer(cflags))
    flag := uint32(cflag)

    events = append(events, PathEvent{
      Path: path,
      Flags: flag,
    })
  }

  ch := watchers[stream].channel
  ch <- events
}
