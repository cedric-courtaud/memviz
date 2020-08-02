package internal

type EventHandler interface {
	HandleAccess(access *Access) error
	HandleCheckpoint(checkpoint *Checkpoint) error
	HandleForked(forked *Forked) error
	Start()
	Stop()
	Finalize()
}
