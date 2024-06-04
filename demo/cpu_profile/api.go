package cpuprofile

import "time"

// StartCPUProfiler uses to start to run the global parallelCPUProfiler.
func StartCPUProfiler(window time.Duration, interval time.Duration) error {
	profileWindow = window
	profileInterval = interval
	return globalCPUProfiler.start()
}

// StopCPUProfiler uses to stop the global parallelCPUProfiler.
func StopCPUProfiler() {
	globalCPUProfiler.stop()
}

// Register register a ProfileConsumer into the global CPU profiler.
// Normally, the registered ProfileConsumer will receive the cpu profile data per second.
// If the ProfileConsumer (channel) is full, the latest cpu profile data will not be sent to it.
// This function is thread-safe.
// WARN: ProfileConsumer should not be closed before unregister.
func Register(ch ProfileConsumer) {
	globalCPUProfiler.register(ch)
}

// Unregister unregister a ProfileConsumer from the global CPU profiler.
// The unregistered ProfileConsumer won't receive the cpu profile data any more.
// This function is thread-safe.
func Unregister(ch ProfileConsumer) {
	globalCPUProfiler.unregister(ch)
}
