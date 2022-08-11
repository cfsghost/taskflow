package taskflow

import "sync"

type TaskLog struct {
	Task             *Task
	Input            *InputSlot
	Output           *OutputSlot
	OutputConnection *Connection
}

type Message struct {
	Context    *Context
	Logs       []*TaskLog
	CurrentLog *TaskLog
	Task       *Task
	Data       interface{}
}

var msgPool = sync.Pool{
	New: func() interface{} {
		return &Message{}
	},
}

func NewMessage(logSize int) *Message {
	m := msgPool.Get().(*Message)
	m.Logs = make([]*TaskLog, 0, logSize)
	return m
}

func (message *Message) ApplyTask(task *Task) {

	taskLog := &TaskLog{
		Task: task,
	}

	message.Logs = append(message.Logs, taskLog)
	message.CurrentLog = taskLog
	message.Task = task
}

func (message *Message) Send(outputID int, result interface{}) error {
	message.Data = result
	return message.Task.Emit(outputID, message)
}

func (message *Message) Clone() *Message {

	curLogSize := len(message.Logs)

	p := NewMessage(curLogSize + 1)
	p.Context = message.Context
	p.Data = message.Data
	p.Task = message.Task
	p.CurrentLog = message.CurrentLog
	p.Logs = append(p.Logs, message.Logs...)

	return p
}

func (message *Message) Release() {
	message.Context = nil
	message.Data = nil
	message.CurrentLog = nil
	message.Logs = nil
	message.Task = nil
	msgPool.Put(message)
}
