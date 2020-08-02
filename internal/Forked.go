package internal

type Forked struct {
	PId               uint32
	PPId              uint32
	ParentAccessCount uint64
}

func (f *Forked) VisitHandler(handler EventHandler) {
	handler.HandleForked(f)
}
