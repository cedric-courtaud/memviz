package internal

import "github.com/cedric-courtaud/memviz/internal/flatbuffers"

type Access struct {
	AccessType flatbuffers.AccessType
	InstAddr   uint64
	DestAddr   uint64
	InstBefore uint64
}

func (a *Access) VisitHandler(handler EventHandler) {
	handler.HandleAccess(a)
}
