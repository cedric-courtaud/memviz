package internal

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


