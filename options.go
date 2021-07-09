package taskflow

type Options struct {
	WorkerCount        int
	SchedulerQueueSize int
}

func NewOptions() *Options {
	return &Options{
		WorkerCount:        4,
		SchedulerQueueSize: 1024000,
	}
}
