package taskflow

import (
	"sync"
	"testing"
)

var testTaskFlow *TaskFlow

func TestCreateTaskFlow(t *testing.T) {

	options := NewOptions()
	tf := NewTaskFlow(options)

	err := tf.Start()
	if err != nil {
		t.Fail()
	}

	testTaskFlow = tf
}

func TestCreateEmptyTask(t *testing.T) {

	task := NewTask(1, 1)
	testTaskFlow.AddTask(task)

	// Check whether task exists or not
	targetTask := testTaskFlow.GetTask(task.id)
	if targetTask == nil {
		t.Fail()
	}
}

func TestRemoveTask(t *testing.T) {

	testTaskFlow.RemoveTask(1)

	// Check whether task exists or not
	targetTask := testTaskFlow.GetTask(2)
	if targetTask != nil {
		t.Fail()
	}
}

func TestCreateCustomizedTask(t *testing.T) {

	str := "Customized Content"
	doneCh := make(chan bool, 1)

	// First task
	task1 := NewTask(1, 1)
	task1.SetHandler(func(message *Message) {
		err := message.Send(0, str)
		if err != nil {
			t.Error(err)
		}
	})
	testTaskFlow.AddTask(task1)

	// Second task
	task2 := NewTask(1, 0)
	task2.SetHandler(func(message *Message) {
		if message.Data.(string) != str {
			doneCh <- false
			return
		}

		doneCh <- true
	})
	testTaskFlow.AddTask(task2)

	// Link two tasks
	testTaskFlow.Link(task1, 0, task2, 0)

	// Push data to task flow
	testTaskFlow.Push(task1.GetID(), 0, "empty")

	success := <-doneCh
	if !success {
		t.Fail()
	}
}

func TestMultipleSend(t *testing.T) {

	str := "Customized Content"
	doneCh := make(chan bool, 10)

	// First task
	task1 := NewTask(1, 1)
	task1.SetHandler(func(message *Message) {

		for i := 0; i < 10; i++ {
			err := message.Send(0, str)
			if err != nil {
				t.Error(err)
			}
		}
	})
	testTaskFlow.AddTask(task1)

	// Second task
	task2 := NewTask(1, 0)
	task2.SetHandler(func(message *Message) {
		if message.Data.(string) != str {
			doneCh <- false
			return
		}

		doneCh <- true
	})
	testTaskFlow.AddTask(task2)

	// Link two tasks
	testTaskFlow.Link(task1, 0, task2, 0)

	// Push data to task flow
	testTaskFlow.Push(task1.GetID(), 0, "empty")

	counter := 0
	for success := range doneCh {
		if !success {
			t.Fail()
		}

		counter++
		if counter == 10 {
			break
		}
	}
}

func TestFanOutData(t *testing.T) {

	str := "Customized Content"
	var wg sync.WaitGroup

	// Source
	source := NewTask(1, 1)
	source.SetHandler(func(message *Message) {
		err := message.Send(0, str)
		if err != nil {
			t.Error(err)
		}
	})
	testTaskFlow.AddTask(source)

	wg.Add(10)
	for i := 0; i < 10; i++ {
		// Linking tasks
		task := NewTask(1, 0)
		task.SetHandler(func(message *Message) {
			if message.Data.(string) != str {
				t.Fail()
				return
			}

			wg.Done()
		})
		testTaskFlow.AddTask(task)

		// Link tasks
		testTaskFlow.Link(source, 0, task, 0)
	}

	// Push data to task flow
	testTaskFlow.Push(source.GetID(), 0, "empty")

	wg.Wait()
}

func TestUnlink(t *testing.T) {

	// Getting task which has connections
	targetTask := testTaskFlow.GetTask(2)
	if targetTask == nil {
		t.Fail()
	}

	totalConn := len(testTaskFlow.GetConnections())

	// Remove connection from task
	slots := targetTask.GetOutputSlots()
	connections := slots[0].GetConnections()
	for _, conn := range connections {
		err := testTaskFlow.Unlink(conn.id)
		if err != nil {
			t.Error(err)
		}

		totalConn--
	}

	// Check
	if totalConn != len(testTaskFlow.GetConnections()) {
		t.Fail()
	}
}
func TestPrivateData(t *testing.T) {

	privData := make(chan interface{}, 1)

	task := NewTask(1, 1)
	task.SetHandler(func(message *Message) {
		privData <- message.Context.GetPrivData()
	})
	testTaskFlow.AddTask(task)

	// Push data to task flow
	ctx := NewContext()
	ctx.SetPrivData("private data")
	testTaskFlow.PushWithContext(task.GetID(), 0, ctx, "test")

	data := <-privData

	if data.(string) != "private data" {
		t.Fail()
	}
}
