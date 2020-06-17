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

type EventDispatcher struct {
	eventQueue chan Event
	Handlers []*EventQueue
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{eventQueue: make(chan Event)}
}

func (e * EventDispatcher) AddHandler(handler EventHandler) {
	e.Handlers = append(e.Handlers, NewEventQueue(&handler))
}

func (e * EventDispatcher) HandleAccess(access *Access) error {
	for _, q := range e.Handlers {
		q.Queue <- access
	}

	return nil
}

func (e * EventDispatcher) HandleCheckpoint(checkpoint *Checkpoint) error {
	for _, q := range e.Handlers {
		q.Queue <- checkpoint
	}
	return nil
}

func (e *EventDispatcher) Start() {
	for _, q := range e.Handlers {
		q.Start()
	}

	go WaitForEvent(e, e.eventQueue, make(chan bool))
}

func (e *EventDispatcher) Stop() {
	for _, q := range e.Handlers {
		q.Stop()
		<- q.Done
	}

	close(e.eventQueue)
}

func WaitForEvent(handler EventHandler, queue chan Event, done chan bool) {
	for event := range queue {
		event.VisitHandler(handler)
	}

	done <- true
}


