package internal

import (
	"github.com/realistschuckle/testify/assert"
	"reflect"
	"testing"
)
/*
func TestAddrSliceSpec_MakeMask(t *testing.T) {

	tests := []struct {
		name   string
		spec  []uint
		want   []uint64
	}{
		{"1 bit", []uint{1}, []uint64{0x1}},
		{"3 fields", []uint{3, 2, 1}, []uint64{0x1, 0x6, 0x38}},
		{"IMX6Spec", []uint{8, 8, 11, 5}, []uint64{ 0x1f, 0xffe0, 0xff0000, 0xff000000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AddrSlicing{
				Slices: tt.spec,
			}
			if got := a.MakeMask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeMask() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
func Test_parseAddrSliceSpec(t *testing.T) {
	example := "1"
	s, err := ParseAddrSlicing(example)
	assert.Nil(t, err)
	assert.Equal(t, *s.Slices[0], AddrSlice{0,1, 0x1})

	example = "3:2:1"
	s, err = ParseAddrSlicing(example)
	assert.Nil(t, err)
	assert.Equal(t, *s.Slices[0], AddrSlice{0,1, 0x1})
	assert.Equal(t, *s.Slices[1], AddrSlice{1,3, 0x6})
	assert.Equal(t, *s.Slices[2], AddrSlice{3,6, 0x38})

	example = "8:8:11:5"
	s, err = ParseAddrSlicing(example)
	assert.Nil(t, err)
	assert.Equal(t, *s.Slices[0], AddrSlice{0,5, 0x1f})
	assert.Equal(t, *s.Slices[1], AddrSlice{5,16, 0xffe0})
	assert.Equal(t, *s.Slices[2], AddrSlice{16,24, 0xff0000})
	assert.Equal(t, *s.Slices[3], AddrSlice{24,32, 0xff000000})
}

func TestNewAddrSlice(t *testing.T) {
	type args struct {
		begin uint
		end   uint
	}
	tests := []struct {
		name string
		args args
		want AddrSlice
	}{
		{"", args{0,1}, AddrSlice{0,1, 0x1}},
		{"", args{0,0}, AddrSlice{0,0, 0x0}},
		{"", args{0,5}, AddrSlice{0,5, 0x1f}},
		{"", args{5,16}, AddrSlice{5,16, 0xffe0}},
		{"", args{16,24}, AddrSlice{16,24, 0xff0000}},
		{"", args{24,32}, AddrSlice{24,32, 0xff000000}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAddrSlice(tt.args.begin, tt.args.end); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewAddrSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}