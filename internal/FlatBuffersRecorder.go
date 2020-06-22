package internal

import (
	"github.com/cedric-courtaud/memviz/internal/flatbuffers"
	fb "github.com/google/flatbuffers/go"
	"io"
)

const DEFAULT_CAPACITY = (1 << 20)

type layers struct {
	instAddr   []uint64
	destAddr   []uint64
	instBefore []uint64
	accessType []flatbuffers.AccessType
}

type FlatBuffersRecorder struct {
	checkpoints []*Checkpoint
	layers      layers

	writer io.Writer
}

func NewFlatBuffersRecorder(writer io.Writer) *FlatBuffersRecorder {
	return &FlatBuffersRecorder{
		layers: layers{
			instAddr:   make([]uint64, 0, DEFAULT_CAPACITY),
			destAddr:   make([]uint64, 0, DEFAULT_CAPACITY),
			instBefore: make([]uint64, 0, DEFAULT_CAPACITY),
			accessType: make([]flatbuffers.AccessType, 0, DEFAULT_CAPACITY),
		},

		writer: writer,
	}
}

func (f *FlatBuffersRecorder) HandleAccess(access *Access) error {
	f.layers.instAddr = append(f.layers.instAddr, access.InstAddr)
	f.layers.destAddr = append(f.layers.destAddr, access.DestAddr)
	f.layers.instBefore = append(f.layers.instBefore, access.InstBefore)
	f.layers.accessType = append(f.layers.accessType, access.AccessType)

	return nil
}

func (f *FlatBuffersRecorder) HandleCheckpoint(checkpoint *Checkpoint) error {
	f.checkpoints = append(f.checkpoints, checkpoint)

	return nil
}

func (f *FlatBuffersRecorder) Start() {
}

func (f *FlatBuffersRecorder) Stop() {
}

func (f FlatBuffersRecorder) Finalize() {
	builder := fb.NewBuilder(DEFAULT_CAPACITY)

	// serialize checkpoints
	var cs []fb.UOffsetT
	for _, checkpoint := range f.checkpoints {
		id := builder.CreateString(checkpoint.Id)
		flatbuffers.CheckpointStart(builder)
		flatbuffers.CheckpointAddId(builder, id)
		flatbuffers.CheckpointAddPos(builder, uint64(checkpoint.Pos))
		flatbuffers.CheckpointAddInstBefore(builder, checkpoint.InstBefore)
		cs = append(cs, flatbuffers.CheckpointEnd(builder))
	}

	// Create checkpoints vector
	flatbuffers.BaseProfileStartCheckpointsVector(builder, len(cs))
	for i := len(cs) - 1; i >= 0; i -= 1 {
		builder.PrependUOffsetT(cs[i])
	}
	checkpoints := builder.EndVector(len(cs))

	// Create layer vectors
	size := len(f.layers.destAddr)

	// InstAddr layer
	flatbuffers.BaseProfileStartInstAddrVector(builder, size)
	for i := size - 1; i >= 0; i -= 1 {
		builder.PrependUint64(f.layers.instAddr[i])
	}
	instAddr := builder.EndVector(size)

	// DestAddr layer
	flatbuffers.BaseProfileStartDestAddrVector(builder, size)
	for i := size - 1; i >= 0; i -= 1 {
		builder.PrependUint64(f.layers.destAddr[i])
	}
	destAddr := builder.EndVector(size)

	// InstBefore layer
	flatbuffers.BaseProfileStartInstBeforeVector(builder, size)
	for i := size - 1; i >= 0; i -= 1 {
		builder.PrependUint64(f.layers.instBefore[i])
	}
	instBefore := builder.EndVector(size)

	// accessType layer
	flatbuffers.BaseProfileStartAccessTypeVector(builder, size)
	for i := size - 1; i >= 0; i -= 1 {
		builder.PrependInt8(int8(f.layers.accessType[i]))
	}
	accessTypes := builder.EndVector(size)

	flatbuffers.BaseProfileStart(builder)
	flatbuffers.BaseProfileAddCheckpoints(builder, checkpoints)
	flatbuffers.BaseProfileAddDestAddr(builder, destAddr)
	flatbuffers.BaseProfileAddInstAddr(builder, instAddr)
	flatbuffers.BaseProfileAddInstBefore(builder, instBefore)
	flatbuffers.BaseProfileAddAccessType(builder, accessTypes)
	profile := flatbuffers.BaseProfileEnd(builder)

	builder.Finish(profile)
	buf := builder.FinishedBytes()

	_, err := f.writer.Write(buf)

	if err != nil {
		panic(err)
	}
}
