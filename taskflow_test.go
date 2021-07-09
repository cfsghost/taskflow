package taskflow

import "testing"

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
	doneCh := make(chan bool)

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

func TestFanOutData(t *testing.T) {

	str := "Customized Content"
	doneCh := make(chan bool)

	// Inherit previous tasks
	source := testTaskFlow.GetTask(2)

	task := NewTask(1, 0)
	task.SetHandler(func(message *Message) {
		if message.Data.(string) != str {
			doneCh <- false
			return
		}

		doneCh <- true
	})
	testTaskFlow.AddTask(task)

	// Link two tasks
	testTaskFlow.Link(source, 0, task, 0)

	// Push data to task flow
	testTaskFlow.Push(source.GetID(), 0, "empty")

	success := <-doneCh
	if !success {
		t.Fail()
	}
}

func TestUnlink(t *testing.T) {

	// Getting task which has connections
	targetTask := testTaskFlow.GetTask(2)
	if targetTask == nil {
		t.Fail()
	}

	// Remove connection from task
	slots := targetTask.GetOutputSlots()
	connections := slots[0].GetConnections()
	for _, conn := range connections {
		err := testTaskFlow.Unlink(conn.id)
		if err != nil {
			t.Error(err)
		}
	}

	// Check
	conns := testTaskFlow.GetConnections()
	if len(conns) > 0 {
		t.Fail()
	}
}
