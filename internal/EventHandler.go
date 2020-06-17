package internal

type EventHandler interface {
	HandleAccess(access *Access) error
	HandleCheckpoint(checkpoint *Checkpoint) error
	Start()
	Stop()
	Finalize()
}
