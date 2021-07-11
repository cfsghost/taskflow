package taskflow

import (
	"github.com/lithammer/go-jump-consistent-hash"
)

type Scheduler struct {
	taskflow *TaskFlow
	queues   []chan *Message
}

func NewScheduler(tf *TaskFlow) *Scheduler {
	return &Scheduler{
		taskflow: tf,
		queues:   make([]chan *Message, tf.options.WorkerCount),
	}
}

func (scheduler *Scheduler) initWorker(queue chan *Message) {

	// Wating for messages
	for message := range queue {
		if message.CurrentLog.OutputConnection == nil {
			continue
		}

		message.CurrentLog.OutputConnection.Execute(message)
	}
}

func (scheduler *Scheduler) Init() error {

	for i := int32(0); i < int32(scheduler.taskflow.options.WorkerCount); i++ {

		// Initializing queue channel
		queue := make(chan *Message, scheduler.taskflow.options.SchedulerQueueSize)
		scheduler.queues[i] = queue

		go scheduler.initWorker(queue)
	}

	return nil
}

func (scheduler *Scheduler) Close() {
	for _, queue := range scheduler.queues {
		close(queue)
	}
}

func (scheduler *Scheduler) Push(key uint64, message *Message) {

	// Calculate pipeline by task
	pipelineID := jump.Hash(key, int32(scheduler.taskflow.options.WorkerCount))

	scheduler.queues[pipelineID] <- message
}
