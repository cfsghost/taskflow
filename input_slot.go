package taskflow

import "sync"

type InputSlot struct {
	id          int
	task        *Task
	connections map[uint64]*Connection
	mutex       sync.RWMutex
}

func NewInputSlot(id int, task *Task) *InputSlot {
	return &InputSlot{
		id:          id,
		task:        task,
		connections: make(map[uint64]*Connection),
	}
}

func (slot *InputSlot) GetConnections() []*Connection {

	slot.mutex.RLock()
	connections := make([]*Connection, len(slot.connections))
	i := 0
	for _, conn := range slot.connections {
		connections[i] = conn
	}
	slot.mutex.RUnlock()

	return connections
}

func (slot *InputSlot) AddConnection(conn *Connection) error {
	slot.mutex.Lock()
	slot.connections[conn.id] = conn
	slot.mutex.Unlock()
	return nil
}

func (slot *InputSlot) RemoveConnection(conn *Connection) error {
	slot.mutex.Lock()
	delete(slot.connections, conn.id)
	slot.mutex.Unlock()

	return nil
}

func (slot *InputSlot) Push(message *Message) {

	message.ApplyTask(slot.task)

	// Update input slot information
	message.CurrentLog.Input = slot

	// Push to task handler
	slot.task.Push(slot.id, message)
}
