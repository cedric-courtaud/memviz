// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flatbuffers

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type BaseProfile struct {
	_tab flatbuffers.Table
}

func GetRootAsBaseProfile(buf []byte, offset flatbuffers.UOffsetT) *BaseProfile {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &BaseProfile{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *BaseProfile) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *BaseProfile) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *BaseProfile) Checkpoints(obj *Checkpoint, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *BaseProfile) CheckpointsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *BaseProfile) AccessType(j int) AccessType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return AccessType(rcv._tab.GetInt8(a + flatbuffers.UOffsetT(j*1)))
	}
	return 0
}

func (rcv *BaseProfile) AccessTypeLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *BaseProfile) MutateAccessType(j int, n AccessType) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateInt8(a+flatbuffers.UOffsetT(j*1), int8(n))
	}
	return false
}

func (rcv *BaseProfile) InstAddr(j int) uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint64(a + flatbuffers.UOffsetT(j*8))
	}
	return 0
}

func (rcv *BaseProfile) InstAddrLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *BaseProfile) MutateInstAddr(j int, n uint64) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint64(a+flatbuffers.UOffsetT(j*8), n)
	}
	return false
}

func (rcv *BaseProfile) DestAddr(j int) uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint64(a + flatbuffers.UOffsetT(j*8))
	}
	return 0
}

func (rcv *BaseProfile) DestAddrLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *BaseProfile) MutateDestAddr(j int, n uint64) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint64(a+flatbuffers.UOffsetT(j*8), n)
	}
	return false
}

func (rcv *BaseProfile) InstBefore(j int) uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint64(a + flatbuffers.UOffsetT(j*8))
	}
	return 0
}

func (rcv *BaseProfile) InstBeforeLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *BaseProfile) MutateInstBefore(j int, n uint64) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint64(a+flatbuffers.UOffsetT(j*8), n)
	}
	return false
}

func BaseProfileStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func BaseProfileAddCheckpoints(builder *flatbuffers.Builder, checkpoints flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(checkpoints), 0)
}
func BaseProfileStartCheckpointsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func BaseProfileAddAccessType(builder *flatbuffers.Builder, AccessType flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(AccessType), 0)
}
func BaseProfileStartAccessTypeVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func BaseProfileAddInstAddr(builder *flatbuffers.Builder, InstAddr flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(InstAddr), 0)
}
func BaseProfileStartInstAddrVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func BaseProfileAddDestAddr(builder *flatbuffers.Builder, DestAddr flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(DestAddr), 0)
}
func BaseProfileStartDestAddrVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func BaseProfileAddInstBefore(builder *flatbuffers.Builder, InstBefore flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(InstBefore), 0)
}
func BaseProfileStartInstBeforeVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func BaseProfileEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
