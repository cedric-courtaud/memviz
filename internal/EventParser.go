package internal

import (
	"bufio"
	"bytes"
	"errors"
	"memrec/internal/flatbuffers"
	"strconv"
)

type Access struct {
	AccessType flatbuffers.AccessType
	InstAddr   uint64
	DestAddr   uint64
	InstBefore uint64
}

// Different from the one flatbuffer
type Checkpoint struct {
	Id         string
	Pos        int
	InstBefore uint64
}

type Event interface {
	VisitHandler(handler EventHandler)
}

type EventHandler interface {
	HandleAccess(access *Access) error
	HandleCheckpoint(checkpoint *Checkpoint) error
	Start()
	Stop()
}

type EventParser struct {
	Handler EventHandler;
	pos int
}

func (p *EventParser) parseLine(line []byte) error {
	if len(line) == 0 {
		return nil
	}

	fields := bytes.Fields(line)

	if bytes.Equal(fields[0], []byte("C")) {
		id := string(fields[1])
		accessBefore, err := strconv.ParseUint(string(fields[2]), 0, 64)
		if err != nil {
			return err
		}
		p.Handler.HandleCheckpoint(&Checkpoint{id, p.pos, accessBefore})

	} else if accessType, ok := flatbuffers.EnumValuesAccessType[string(fields[0])]; ok {
		instAddr, err := strconv.ParseUint(string(fields[1]), 0, 64)
		if err != nil {
			return err
		}

		destAddr, err := strconv.ParseUint(string(fields[2]), 0, 64)
		if err != nil {
			return err
		}

		accessBefore, err := strconv.ParseUint(string(fields[3]), 0, 64)
		if err != nil {
			return err
		}

		p.pos = p.pos + 1
		p.Handler.HandleAccess(&Access{accessType, instAddr, destAddr, accessBefore})
	} else {
		return errors.New("Invalid event identifier")
	}

	return nil
}

func (p *EventParser) Parse(reader * bufio.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		err := p.parseLine(scanner.Bytes())

		if err != nil {
			return err
		}
	}

	return nil
}



func (a * Access) VisitHandler (handler EventHandler) {
	handler.HandleAccess(a)
}

func (c * Checkpoint) VisitHandler (handler EventHandler) {
	handler.HandleCheckpoint(c)
}
