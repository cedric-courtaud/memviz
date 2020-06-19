package internal

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEventDispatcher(t *testing.T) {
	l1 := EventLogger{}
	l2 := EventLogger{}
	d := EventDispatcher{eventQueue: make(chan Event)}

	d.AddHandler(&l1)
	d.AddHandler(&l2)
	d.Start()

	p := EventParser{&d, 0}

	example := []string {"C p1 12",
		"I 0x11 0x12 13",
		"R 0x11 0x12 13",
		"W 0x11 0x12 13",
		"C p3 14"}

	r := bufio.NewReader(strings.NewReader(strings.Join(example, "\n")))

	p.Parse(r)

	d.Stop()

	for i, _ := range example {
		assert.Equal(t, example[i], l1.Events[i])
		assert.Equal(t, example[i], l2.Events[i])
	}
}