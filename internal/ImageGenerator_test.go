package internal

import (
	"github.com/cedric-courtaud/memviz/internal/flatbuffers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImageGenerator_getXPos(t *testing.T) {
	slicing, _ := ParseAddrSlicing("8:8:11:5")
	generator := NewImageGenerator(10, 10, slicing)

	assert.Equal(t, generator.getXPos(1), 0)
	assert.Equal(t, generator.getXPos(10), 1)
	assert.Equal(t, generator.getXPos(11), 1)
	assert.Equal(t, generator.getXPos(34), 3)
}

func TestYPositioner_GetYPos(t *testing.T) {
	slicing, _ := ParseAddrSlicing("2:2")
	p1 := NewYPositioner(slicing.Slices[0], 4, slicing.Total)
	p2 := NewYPositioner(slicing.Slices[1], 4, slicing.Total)

	assert.Equal(t, p1.GetYPos(0), 0)
	assert.Equal(t, p2.GetYPos(0), 2)

	assert.Equal(t, p1.GetYPos(1), 0)
	assert.Equal(t, p2.GetYPos(1), 2)

	assert.Equal(t, p1.GetYPos(2), 1)
	assert.Equal(t, p2.GetYPos(2), 3)

	assert.Equal(t, p1.GetYPos(3), 1)
	assert.Equal(t, p2.GetYPos(3), 3)
}

func TestImageGenerator_HandleAccess(t *testing.T) {
	slicing, _ := ParseAddrSlicing("2:2")
	generator := NewImageGenerator(4, 100, slicing)

	// C0
	generator.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeI,
		DestAddr:   0x0,
		InstAddr:   0x0,
		InstBefore: 0,
	})

	// + 0x8

	generator.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeR,
		DestAddr:   0x8,
		InstAddr:   0x0,
		InstBefore: 50,
	})
	// C1
	// + 0x4
	generator.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeW,
		DestAddr:   0xc,
		InstAddr:   0x0,
		InstBefore: 50,
	})

	// + 0x4
	generator.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeR,
		DestAddr:   0x0,
		InstAddr:   0x0,
		InstBefore: 50,
	})

	b := generator.buff
	assert.Equal(t, generator.Width, 2)

	assert.Equal(t, b[0][0].NI, 0)
	assert.Equal(t, b[0][0].NR, 1)
	assert.Equal(t, b[0][0].NW, 0)

	assert.Equal(t, b[0][1].NI, 0)
	assert.Equal(t, b[0][1].NR, 0)
	assert.Equal(t, b[0][1].NW, 0)

	assert.Equal(t, b[0][2].NI, 0)
	assert.Equal(t, b[0][2].NR, 0)
	assert.Equal(t, b[0][2].NW, 0)

	assert.Equal(t, b[0][3].NI, 0)
	assert.Equal(t, b[0][3].NR, 1)
	assert.Equal(t, b[0][3].NW, 0)

	assert.Equal(t, b[1][0].NI, 0)
	assert.Equal(t, b[1][0].NR, 1)
	assert.Equal(t, b[1][0].NW, 1)

	assert.Equal(t, b[1][1].NI, 0)
	assert.Equal(t, b[1][1].NR, 0)
	assert.Equal(t, b[1][1].NW, 0)

	assert.Equal(t, b[1][2].NI, 0)
	assert.Equal(t, b[1][2].NR, 1)
	assert.Equal(t, b[1][2].NW, 1)

	assert.Equal(t, b[1][3].NI, 0)
	assert.Equal(t, b[1][3].NR, 0)
	assert.Equal(t, b[1][3].NW, 0)
}
