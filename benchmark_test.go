package taskflow

import (
	"sync"
	"testing"
)

var benchTaskFlow *TaskFlow

func BenchmarkSingleTask(b *testing.B) {

	var wg sync.WaitGroup

	// Create a taskflow
	options := NewOptions()
	tf := NewTaskFlow(options)

	err := tf.Start()
	if err != nil {
		b.Error(err)
	}

	defer tf.Stop()
	benchTaskFlow = tf

	// Create a task
	task1 := NewTask(1, 0)
	task1.SetHandler(func(message *Message) {
		wg.Done()
	})
	benchTaskFlow.AddTask(task1)

	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		// Push data to task flow
		err := benchTaskFlow.Push(task1.GetID(), 0, "empty")
		if err != nil {
			b.Error(err)
		}
	}

	wg.Wait()
}

func BenchmarkTwoTasks(b *testing.B) {

	var wg sync.WaitGroup

	// Create a taskflow
	options := NewOptions()
	tf := NewTaskFlow(options)

	err := tf.Start()
	if err != nil {
		b.Error(err)
	}

	defer tf.Stop()
	benchTaskFlow = tf

	// Create a task
	task1 := NewTask(1, 1)
	task1.SetHandler(func(message *Message) {
		err := message.Send(0, "TEST")
		if err != nil {
			b.Error(err)
		}
	})
	benchTaskFlow.AddTask(task1)

	// Create final task
	task2 := NewTask(1, 0)
	task2.SetHandler(func(message *Message) {
		wg.Done()
	})
	benchTaskFlow.AddTask(task2)

	// Link two tasks
	benchTaskFlow.Link(task1, 0, task2, 0)

	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		// Push data to task flow
		err := benchTaskFlow.Push(task1.GetID(), 0, "empty")
		if err != nil {
			b.Error(err)
		}
	}

	wg.Wait()
}
