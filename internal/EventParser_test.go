package internal

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"memrec/internal/flatbuffers"
	"strings"
	"testing"
)


func TestEventParser_Parse(t *testing.T) {
	example := []string {"C p1 12",
		"I 0x11 0x12 13",
		"R 0x11 0x12 13",
		"W 0x11 0x12 13",
		"C p3 14"}

	r := bufio.NewReader(strings.NewReader(strings.Join(example, "\n")))

	h := EventLogger{}
	p := EventParser{&h, 0}

	assert.Nil(t, p.Parse(r))

	for i := range example {
		assert.Equal(t, example[i], h.Events[i])
	}

	assert.Equal(t, p.pos, 3)
}

func TestEventParser_parseLine(t *testing.T) {
	h := EventLogger{}
	p := EventParser{&h, 0}

	example := []byte("C p1 0")
	assert.Nil(t, p.parseLine(example))
	assert.Equal(t, h.Events[0], string(example))

	example = []byte("I 0x11 0x12 13")
	assert.Nil(t, p.parseLine(example))
	assert.Equal(t, h.Events[1], string(example))

	example = []byte("R 0x11 0x12 13")
	assert.Nil(t, p.parseLine(example))
	assert.Equal(t, h.Events[2], string(example))

	example = []byte("W 0x11 0x12 13")
	assert.Nil(t, p.parseLine(example))
	assert.Equal(t, h.Events[3], string(example))

	example = []byte(" 0x11 0x12 13")
	assert.Error(t, p.parseLine(example))
}

func TestEventStringer_HandleAccess(t *testing.T) {
	s := EventLogger{}

	_ = s.HandleAccess(&Access{flatbuffers.AccessTypeI, 0x11, 0x12, 13})
	assert.Equal(t, s.Events[0], "I 0x11 0x12 13")

	_ = s.HandleAccess(&Access{flatbuffers.AccessTypeW, 0x11, 0x12, 13})
	assert.Equal(t, s.Events[1], "W 0x11 0x12 13")

	_ = s.HandleAccess(&Access{flatbuffers.AccessTypeR, 0x11, 0x12, 13})
	assert.Equal(t, s.Events[2], "R 0x11 0x12 13")
}

func TestEventStringer_HandleCheckpoint(t *testing.T) {
	s := EventLogger{}
	_ = s.HandleCheckpoint(&Checkpoint{"__p1", 0, 0})
	assert.Equal(t, s.Events[0], "C __p1 0")
}