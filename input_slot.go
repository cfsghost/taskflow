package taskflow

type InputSlot struct {
	id   int
	task *Task
}

func NewInputSlot(id int, task *Task) *InputSlot {
	return &InputSlot{
		id:   id,
		task: task,
	}
}

func (slot *InputSlot) Push(message *Message) {

	message.ApplyTask(slot.task)

	// Update input slot information
	message.CurrentLog.Input = slot

	// Push to task handler
	slot.task.Push(slot.id, message)
}
