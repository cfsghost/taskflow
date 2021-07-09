package taskflow

import (
	"fmt"
	"sync"
)

type TaskHandler func(*Message)

type Task struct {
	taskflow *TaskFlow
	id       uint64
	name     string
	inputs   []*InputSlot
	outputs  []*OutputSlot
	handler  TaskHandler

	inputMutex  sync.RWMutex
	outputMutex sync.RWMutex
}

func NewTask(inputCount int, outputCount int) *Task {

	task := &Task{
		inputs:  make([]*InputSlot, inputCount),
		outputs: make([]*OutputSlot, outputCount),
		handler: func(message *Message) {
			message.Send(0, message)
		},
	}

	// Initializing input slots
	for i := 0; i < inputCount; i++ {
		slot := NewInputSlot(i, task)
		task.inputs[i] = slot
	}

	// Initializing output slots
	for i := 0; i < outputCount; i++ {
		slot := NewOutputSlot(i, task)
		task.outputs[i] = slot
	}

	return task
}

func (task *Task) GetID() uint64 {
	return task.id
}

func (task *Task) SetName(name string) {
	task.name = name
}

func (task *Task) GetName() string {
	return task.name
}

func (task *Task) SetHandler(fn func(*Message)) {
	task.handler = fn
}

func (task *Task) GetInputSlots() []*InputSlot {

	task.inputMutex.RLock()
	slots := make([]*InputSlot, len(task.inputs))

	for i, slot := range task.inputs {
		slots[i] = slot
	}

	task.inputMutex.RUnlock()

	return slots
}

func (task *Task) GetInputSlot(id int) *InputSlot {

	task.inputMutex.RLock()
	defer task.inputMutex.RUnlock()

	if id < 0 || id >= len(task.inputs) {
		return nil
	}

	return task.inputs[id]
}

func (task *Task) GetOutputSlots() []*OutputSlot {

	task.outputMutex.RLock()
	slots := make([]*OutputSlot, len(task.outputs))

	for i, slot := range task.outputs {
		slots[i] = slot
	}

	task.outputMutex.RUnlock()

	return slots
}

func (task *Task) GetOutputSlot(id int) *OutputSlot {

	task.outputMutex.RLock()
	defer task.outputMutex.RUnlock()

	if id < 0 || id >= len(task.outputs) {
		return nil
	}

	return task.outputs[id]
}

func (task *Task) Push(inputSlotID int, message *Message) error {

	// Getting
	slot := task.GetInputSlot(inputSlotID)
	if slot == nil {
		return fmt.Errorf("No such input slot: %d", inputSlotID)
	}

	task.handler(message)

	return nil
}

func (task *Task) Emit(outputSlotID int, message *Message) error {

	// Getting
	slot := task.GetOutputSlot(outputSlotID)
	if slot == nil {
		return fmt.Errorf("No such output slot: %d", outputSlotID)
	}

	slot.Push(message)

	return nil
}
