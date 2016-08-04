package runtime

import (
	"runtime"
	"time"

	"github.com/segmentio/kit/log"
	"github.com/segmentio/kit/schema"
	"github.com/segmentio/kit/stats"
)

// Package runtime reports metrics, in a different goroutine,
// to the stats package every two minutes:
// 		- Number of goroutines (as Positive Number)
//      - Allocated memory (as Megabytes)
//		- Heap memory (as Megabytes)
//		- Stack memory in use (as Megabytes)
//		- Garbage Collector pauses (as Milliseconds)
func Init(serviceSchema schema.Service) error {
	go reportRuntime()
	return nil
}

func reportRuntime() {
	for {
		// TODO(vince): Make this a variable
		time.Sleep(2 * time.Minute)
		log.Debug("Sending runtime metrics to provider")

		// Goroutines
		stats.Gauge("numGoroutines", runtime.NumGoroutine())

		// Memory profiling
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		// Memory reports in MB
		stats.Gauge("memAllocated", int(m.Alloc/1000000))
		stats.Gauge("memHeap", int(m.HeapAlloc/1000000))
		stats.Gauge("memStackInUse", int(m.StackInuse/1000000))
		// Pause reports in milliseconds
		stats.Gauge("gcPause", int((m.PauseNs[(m.NumGC+255)%256])/1000000))
	}
}
