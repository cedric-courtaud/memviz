package internal

import (
	"fmt"
	"memrec/internal/flatbuffers"
)

type EventLogger struct {
	Events []string
}

func (e *EventLogger) Finalize() {
}

func (e *EventLogger) Start() {
}

func (e *EventLogger) Stop() {
}

func (e *EventLogger) HandleAccess(access *Access) error {
	e.Events = append(e.Events, fmt.Sprintf("%s 0x%x 0x%x %d", flatbuffers.EnumNamesAccessType[access.AccessType],
		access.InstAddr, access.DestAddr, access.InstBefore))

	return nil
}

func (e *EventLogger) HandleCheckpoint(c * Checkpoint) error {
	e.Events = append(e.Events, fmt.Sprintf("C %s %d", c.Id, c.InstBefore))
	return nil
}

