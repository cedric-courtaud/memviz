package internal

import (
	"github.com/realistschuckle/testify/assert"
	"memrec/internal/flatbuffers"
	"testing"
)

func Test_phaseStats_handleAccess(t *testing.T) {
	conf, _ := NewStatsConfig("2:2")
	p := newPhaseStats("p1", conf)

	p.handleAccess(&Access{
		AccessType: flatbuffers.AccessTypeI,
		InstAddr:   0x0,
		DestAddr:   0x0,
		InstBefore: 0,
	})

	assert.Equal(t, p.AccessCount, uint64(1))
	assert.Equal(t, p.InversionCount, uint64(0))

	assert.Equal(t, p.addrDiffCount[3][0],  uint64(0))
	assert.Equal(t, p.addrDiffCount[3][1],  uint64(0))
	assert.Equal(t, p.addrDiffCount[3][2],  uint64(0))
	assert.Equal(t, p.addrDiffCount[3][3],  uint64(0))

	assert.Equal(t, p.addrDiffCount[12][0], uint64(0))
	assert.Equal(t, p.addrDiffCount[12][1], uint64(0))
	assert.Equal(t, p.addrDiffCount[12][2], uint64(0))
	assert.Equal(t, p.addrDiffCount[12][3], uint64(0))

	p.handleAccess(&Access{
		AccessType: flatbuffers.AccessTypeR,
		InstAddr:   0x0,
		DestAddr:   0x5,
		InstBefore: 2,
	})

	assert.Equal(t, p.AccessCount, uint64(2))
	assert.Equal(t, p.InversionCount, uint64(0))

	assert.Equal(t, p.addrDiffCount[3][0],  uint64(0))
	assert.Equal(t, p.addrDiffCount[3][1],  uint64(1))
	assert.Equal(t, p.addrDiffCount[3][2],  uint64(0))
	assert.Equal(t, p.addrDiffCount[3][3],  uint64(0))

	assert.Equal(t, p.addrDiffCount[12][0], uint64(0))
	assert.Equal(t, p.addrDiffCount[12][1], uint64(1))
	assert.Equal(t, p.addrDiffCount[12][2], uint64(0))
	assert.Equal(t, p.addrDiffCount[12][3], uint64(0))

	p.handleAccess(&Access{
		AccessType: flatbuffers.AccessTypeW,
		InstAddr:   0x0,
		DestAddr:   0xf,
		InstBefore: 1,
	})

	assert.Equal(t, p.AccessCount, uint64(3))
	assert.Equal(t, p.InversionCount, uint64(1))

	assert.Equal(t, p.addrDiffCount[3][0],  uint64(0))
	assert.Equal(t, p.addrDiffCount[3][1],  uint64(1))
	assert.Equal(t, p.addrDiffCount[3][2],  uint64(1))
	assert.Equal(t, p.addrDiffCount[3][3],  uint64(0))

	assert.Equal(t, p.addrDiffCount[12][0], uint64(0))
	assert.Equal(t, p.addrDiffCount[12][1], uint64(1))
	assert.Equal(t, p.addrDiffCount[12][2], uint64(1))
	assert.Equal(t, p.addrDiffCount[12][3], uint64(0))

	p.handleAccess(&Access{
		AccessType: flatbuffers.AccessTypeR,
		InstAddr:   0x0,
		DestAddr:   0x0,
		InstBefore: 1,
	})

	assert.Equal(t, p.AccessCount, uint64(4))
	assert.Equal(t, p.InversionCount, uint64(2))

	assert.Equal(t, p.addrDiffCount[3][0],  uint64(0))
	assert.Equal(t, p.addrDiffCount[3][1],  uint64(2))
	assert.Equal(t, p.addrDiffCount[3][2],  uint64(1))
	assert.Equal(t, p.addrDiffCount[3][3],  uint64(0))

	assert.Equal(t, p.addrDiffCount[12][0], uint64(0))
	assert.Equal(t, p.addrDiffCount[12][1], uint64(2))
	assert.Equal(t, p.addrDiffCount[12][2], uint64(1))
	assert.Equal(t, p.addrDiffCount[12][3], uint64(0))
}

func TestShannonEntropy(t *testing.T) {
	m := make(map[uint64]uint64)

	h := ShannonEntropy(m, 0)
	assert.Equal(t, h, 0.0)

	m[0] += 1

	h = ShannonEntropy(m, 1)
	assert.Equal(t, h, 0.0)

	m[1] += 1
	m[2] += 1
	m[3] += 1

	h = ShannonEntropy(m, 4)
	assert.Equal(t, h, 2.0)
}

func TestStreamStats(t *testing.T) {
	conf, _ := NewStatsConfig("4")
	stats := NewStreamStats(conf)

	err := stats.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeI,
		InstAddr:   0x0,
		DestAddr:   0x0,
		InstBefore: 0,
	})

	assert.NotNil(t, err)

	stats.HandleCheckpoint(&Checkpoint{"p1", 0, 0})
	p, err := stats.getCurrentPhase()
	assert.Nil(t, err)
	assert.Equal(t, p.AccessCount, uint64(0))

	err = stats.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeR,
		InstAddr:   0x0,
		DestAddr:   0x4,
		InstBefore: 1,
	})
	assert.Nil(t, err)

	err = stats.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeW,
		InstAddr:   0x0,
		DestAddr:   0x8,
		InstBefore: 2,
	})
	assert.Nil(t, err)

	assert.Equal(t, p.AccessCount, uint64(2))
	assert.Equal(t, p.InversionCount, uint64(1))

	assert.Equal(t, p.addrDiffCount[0xf][0], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][1], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][2], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][3], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][4], uint64(1))
	assert.Equal(t, p.addrDiffCount[0xf][5], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][6], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][7], uint64(0))

	stats.HandleCheckpoint(&Checkpoint{"p1", 0, 0})
	p, err = stats.getCurrentPhase()
	assert.Nil(t, err)
	assert.Equal(t, p.AccessCount, uint64(0))

	err = stats.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeR,
		InstAddr:   0x0,
		DestAddr:   0xc,
		InstBefore: 2,
	})
	assert.Nil(t, err)

	assert.Equal(t, p.AccessCount, uint64(1))
	assert.Equal(t, p.InversionCount, uint64(1))

	assert.Equal(t, p.addrDiffCount[0xf][0], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][1], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][2], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][3], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][4], uint64(1))
	assert.Equal(t, p.addrDiffCount[0xf][5], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][6], uint64(0))
	assert.Equal(t, p.addrDiffCount[0xf][7], uint64(0))
}