package taskflow

import (
	"sync"
)

type OutputSlot struct {
	id          int
	task        *Task
	connections map[uint64]*Connection
	mutex       sync.RWMutex
}

func NewOutputSlot(id int, task *Task) *OutputSlot {
	return &OutputSlot{
		id:          id,
		task:        task,
		connections: make(map[uint64]*Connection),
	}
}

func (slot *OutputSlot) GetConnections() []*Connection {

	slot.mutex.RLock()
	connections := make([]*Connection, len(slot.connections))
	i := 0
	for _, conn := range slot.connections {
		connections[i] = conn
		i++
	}
	slot.mutex.RUnlock()

	return connections
}

func (slot *OutputSlot) AddConnection(conn *Connection) error {
	slot.mutex.Lock()
	slot.connections[conn.id] = conn
	slot.mutex.Unlock()
	return nil
}

func (slot *OutputSlot) RemoveConnection(conn *Connection) error {
	slot.mutex.Lock()
	delete(slot.connections, conn.id)
	slot.mutex.Unlock()

	return nil
}

func (slot *OutputSlot) Push(message *Message) {

	// Update output slot information
	message.CurrentLog.Output = slot

	slot.mutex.RLock()
	//	isFirstConnection := true
	for _, targetConn := range slot.connections {

		//		if !isFirstConnection {
		// Clone a new message which inherited logs
		m := message.Clone()
		m.CurrentLog.OutputConnection = targetConn

		// Push to slot of next task
		targetConn.Push(m)
		/*
				continue
			}

			isFirstConnection = false

			// First slot has original message to push
			message.CurrentLog.OutputConnection = targetConn

			targetConn.Push(message)
		*/
	}
	slot.mutex.RUnlock()

}
