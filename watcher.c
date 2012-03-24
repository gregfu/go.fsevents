#include <CoreServices/CoreServices.h>
#include "_cgo_export.h"

void
fswatch_callback(ConstFSEventStreamRef streamRef,
                 void *clientCallBackInfo,
                 size_t numEvents,
                 void *eventPaths,
                 const FSEventStreamEventFlags eventFlags[],
                 const FSEventStreamEventId eventIds[])
{
  watchDirsCallback(
      (FSEventStreamRef)streamRef,
      numEvents,
      eventPaths,
      (FSEventStreamEventFlags*)eventFlags);
}

FSEventStreamRef fswatch_stream_for_paths(CFMutableArrayRef pathsToWatch) {
  return FSEventStreamCreate(
      NULL,
      fswatch_callback,
      NULL,
      pathsToWatch,
      kFSEventStreamEventIdSinceNow,
      0.1,
      kFSEventStreamCreateFlagNoDefer | kFSEventStreamCreateFlagFileEvents);
}
