package internal

type Event interface {
	VisitHandler(handler EventHandler)
}
