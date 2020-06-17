package internal

type Checkpoint struct {
	Id         string
	Pos        int
	InstBefore uint64
}

func (c * Checkpoint) VisitHandler (handler EventHandler) {
	handler.HandleCheckpoint(c)
}
