package taskflow

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type TaskFlow struct {
	taskCounter uint64
	connCounter uint64
	tasks       sync.Map
	connections sync.Map
	options     *Options
	scheduler   *Scheduler
}

func NewTaskFlow(options *Options) *TaskFlow {

	tf := &TaskFlow{
		options: options,
	}

	tf.scheduler = NewScheduler(tf)

	return tf
}

func (tf *TaskFlow) Start() error {

	err := tf.scheduler.Init()

	return err
}

func (tf *TaskFlow) Stop() {
	tf.scheduler.Close()
}

func (tf *TaskFlow) GetScheduler() *Scheduler {
	return tf.scheduler
}

func (tf *TaskFlow) AddTask(task *Task) uint64 {
	task.taskflow = tf
	task.id = atomic.AddUint64(&tf.taskCounter, 1)
	tf.tasks.Store(task.id, task)
	return task.id
}

func (tf *TaskFlow) RemoveTask(id uint64) {

	task := tf.GetTask(id)
	if task == nil {
		return
	}

	// Unlink all for the task
	for _, slot := range task.GetInputSlots() {

		connections := slot.GetConnections()
		for _, conn := range connections {
			conn.Close()
			tf.connections.Delete(conn.id)
		}
	}

	for _, slot := range task.GetOutputSlots() {

		connections := slot.GetConnections()
		for _, conn := range connections {
			conn.Close()
			tf.connections.Delete(conn.id)
		}
	}

	tf.tasks.Delete(id)
}

func (tf *TaskFlow) GetTask(id uint64) *Task {
	v, ok := tf.tasks.Load(id)
	if !ok {
		return nil
	}

	return v.(*Task)
}

func (tf *TaskFlow) Push(id uint64, inputSlotID int, data interface{}) error {
	return tf.PushWithContext(id, inputSlotID, NewContext(), data)
}

func (tf *TaskFlow) PushWithContext(id uint64, inputSlotID int, ctx *Context, data interface{}) error {

	task := tf.GetTask(id)
	if task == nil {
		return fmt.Errorf("No such task: %d", id)
	}

	slot := task.GetInputSlot(inputSlotID)
	if slot == nil {
		return fmt.Errorf("Not found input slot \"%d\" of task \"%s\"(id=%d)", inputSlotID, task.name, id)
	}

	// Prepare payload and push to input slot
	payload := NewMessage()
	payload.Context = ctx
	payload.Data = data
	slot.Push(payload)

	return nil
}

func (tf *TaskFlow) Link(source *Task, sourceOutputSlotID int, dest *Task, destInputSlotID int) (uint64, error) {

	// Getting source slot
	outputSlot := source.GetOutputSlot(sourceOutputSlotID)
	if outputSlot == nil {
		return 0, fmt.Errorf("Not found output slot \"%d\" of task \"%s\"(id=%d)", sourceOutputSlotID, source.name, source.id)
	}

	// Getting destination slot
	inputSlot := dest.GetInputSlot(destInputSlotID)
	if inputSlot == nil {
		return 0, fmt.Errorf("Not found input slot \"%d\" of task \"%s\"(id=%d)", destInputSlotID, dest.name, dest.id)
	}

	// Create a new connection
	connID := atomic.AddUint64(&tf.connCounter, 1)
	connection := NewConnection(connID, outputSlot, inputSlot)
	tf.connections.Store(connID, connection)

	return connID, connection.Apply()
}

func (tf *TaskFlow) Unlink(id uint64) error {

	conn := tf.GetConnection(id)
	if conn == nil {
		return fmt.Errorf("No such connection: %d", id)
	}

	err := conn.Close()
	if err != nil {
		return err
	}

	tf.connections.Delete(id)

	return nil
}

func (tf *TaskFlow) GetConnections() []*Connection {

	connections := make([]*Connection, 0)
	tf.connections.Range(func(k interface{}, v interface{}) bool {
		connections = append(connections, v.(*Connection))
		return true
	})

	return connections
}

func (tf *TaskFlow) GetConnection(id uint64) *Connection {

	v, ok := tf.connections.Load(id)
	if !ok {
		return nil
	}

	return v.(*Connection)
}
