package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/cedric-courtaud/memviz/internal/flatbuffers"
	"strconv"
)

type EventParser struct {
	Handler EventHandler
	pos     int
}

func NewEventParser(h EventHandler) *EventParser {
	return &EventParser{h, 0}
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
		s := fmt.Sprintf("Line %d: %s\nInvalid event identifier", p.pos, line)
		return errors.New(s)
	}

	return nil
}

func (p *EventParser) start(done chan bool, errc chan error, queue chan []byte) {
	for line := range queue {
		err := p.parseLine(line)
		if err != nil {
			errc <- err
		}
	}

	done <- true
}

func (p *EventParser) Parse(reader *bufio.Reader) error {
	scanner := bufio.NewScanner(reader)
	done := make(chan bool)
	queue := make(chan []byte)
	errc := make(chan error, 1)

	go p.start(done, errc, queue)

	for scanner.Scan() {
		line := scanner.Bytes()
		buff := make([]byte, len(line))
		copy(buff, line)
		queue <- buff

		select {
		case err := <-errc:
			return err
		default:
			continue
		}
	}

	close(queue)
	<-done

	return nil
}
