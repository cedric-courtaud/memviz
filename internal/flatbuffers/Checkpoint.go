// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flatbuffers

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Checkpoint struct {
	_tab flatbuffers.Table
}

func GetRootAsCheckpoint(buf []byte, offset flatbuffers.UOffsetT) *Checkpoint {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Checkpoint{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Checkpoint) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Checkpoint) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Checkpoint) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Checkpoint) Pos() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Checkpoint) MutatePos(n uint64) bool {
	return rcv._tab.MutateUint64Slot(6, n)
}

func (rcv *Checkpoint) InstBefore() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Checkpoint) MutateInstBefore(n uint64) bool {
	return rcv._tab.MutateUint64Slot(8, n)
}

func CheckpointStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func CheckpointAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func CheckpointAddPos(builder *flatbuffers.Builder, pos uint64) {
	builder.PrependUint64Slot(1, pos, 0)
}
func CheckpointAddInstBefore(builder *flatbuffers.Builder, instBefore uint64) {
	builder.PrependUint64Slot(2, instBefore, 0)
}
func CheckpointEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
