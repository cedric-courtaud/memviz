package internal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"memrec/internal/flatbuffers"
	"testing"
)

func TestFlatBuffersRecorder(t *testing.T) {
	buf := new(bytes.Buffer)
	rec := NewFlatBuffersRecorder(buf)

	rec.HandleCheckpoint(&Checkpoint{
		Id:         "p1",
		Pos:        0,
		InstBefore: 1,
	})

	rec.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeI,
		InstAddr:   0x12,
		DestAddr:   0x13,
		InstBefore: 14,
	})

	rec.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeR,
		InstAddr:   0x22,
		DestAddr:   0x23,
		InstBefore: 24,
	})

	rec.HandleCheckpoint(&Checkpoint{
		Id:         "p2",
		Pos:        3,
		InstBefore: 2,
	})

	rec.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeW,
		InstAddr:   0x32,
		DestAddr:   0x33,
		InstBefore: 34,
	})

	rec.HandleAccess(&Access{
		AccessType: flatbuffers.AccessTypeR,
		InstAddr:   0x42,
		DestAddr:   0x43,
		InstBefore: 44,
	})

	rec.Finalize()

	profile := flatbuffers.GetRootAsBaseProfile(buf.Bytes(),0)

	assert.Equal(t, profile.CheckpointsLength(), 2)
	c := flatbuffers.Checkpoint{}
	profile.Checkpoints(&c, 0)
	assert.Equal(t, c.Pos(), uint64(0))
	assert.Equal(t, c.Id(), []byte("p1"))
	assert.Equal(t, c.InstBefore(), uint64(1))

	profile.Checkpoints(&c, 1)
	assert.Equal(t, c.Pos(), uint64(3))
	assert.Equal(t, c.Id(), []byte("p2"))
	assert.Equal(t, c.InstBefore(), uint64(2))

	assert.Equal(t, profile.DestAddrLength(), 4)

	assert.Equal(t, profile.AccessType(0), flatbuffers.AccessTypeI)
	assert.Equal(t, profile.InstAddr(0), uint64(0x12))
	assert.Equal(t, profile.DestAddr(0), uint64(0x13))
	assert.Equal(t, profile.InstBefore(0), uint64(14))

	assert.Equal(t, profile.AccessType(1), flatbuffers.AccessTypeR)
	assert.Equal(t, profile.InstAddr(1), uint64(0x22))
	assert.Equal(t, profile.DestAddr(1), uint64(0x23))
	assert.Equal(t, profile.InstBefore(1), uint64(24))

	assert.Equal(t, profile.AccessType(2), flatbuffers.AccessTypeW)
	assert.Equal(t, profile.InstAddr(2), uint64(0x32))
	assert.Equal(t, profile.DestAddr(2), uint64(0x33))
	assert.Equal(t, profile.InstBefore(2), uint64(34))

	assert.Equal(t, profile.AccessType(3), flatbuffers.AccessTypeR)
	assert.Equal(t, profile.InstAddr(3), uint64(0x42))
	assert.Equal(t, profile.DestAddr(3), uint64(0x43))
	assert.Equal(t, profile.InstBefore(3), uint64(44))
}

