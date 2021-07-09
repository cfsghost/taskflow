package taskflow

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

func NewMessage() *Message {
	return &Message{
		Logs: make([]*TaskLog, 0),
	}
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

	p := NewMessage()
	p.Context = message.Context
	p.Data = message.Data
	p.Task = message.Task
	p.CurrentLog = message.CurrentLog

	for _, taskLog := range message.Logs {
		p.Logs = append(p.Logs, taskLog)
	}

	return p
}
