package internal

type EventQueue struct {
	Queue chan Event
	Done chan bool
	handler *EventHandler
}

func NewEventQueue(handler *EventHandler) *EventQueue {
	return &EventQueue{handler: handler, Queue: make(chan Event), Done: make(chan bool)}
}

func (e *EventQueue) Start() {
	go WaitForEvent(*e.handler, e.Queue, e.Done)
}

func (e *EventQueue) Stop() {
	close(e.Queue)
}

func (e *EventQueue) HandleAccess(access *Access) error {
	e.Queue <- access
	return nil
}

func (e *EventQueue) HandleCheckpoint(checkpoint *Checkpoint) error {
	e.Queue <- checkpoint
	return nil
}
