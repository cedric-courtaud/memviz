package internal

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_parseAddrSliceSpec(t *testing.T) {
	example := "1"
	s, err := ParseAddrSlicing(example)
	assert.Nil(t, err)
	assert.Equal(t, s.Total, 1)
	assert.Equal(t, *s.Slices[0], AddrSlice{0,1, 0x1, 2})

	example = "3:2:1"
	s, err = ParseAddrSlicing(example)
	assert.Nil(t, err)
	assert.Equal(t, s.Total, 6)
	assert.Equal(t, *s.Slices[0], AddrSlice{0,1, 0x1, 2})
	assert.Equal(t, *s.Slices[1], AddrSlice{1,3, 0x6, 4})
	assert.Equal(t, *s.Slices[2], AddrSlice{3,6, 0x38, 8})

	example = "8:8:11:5"
	s, err = ParseAddrSlicing(example)
	assert.Nil(t, err)
	assert.Equal(t, s.Total, 32)
	assert.Equal(t, *s.Slices[0], AddrSlice{0,5, 0x1f, 32})
	assert.Equal(t, *s.Slices[1], AddrSlice{5,16, 0xffe0, 2048})
	assert.Equal(t, *s.Slices[2], AddrSlice{16,24, 0xff0000, 256})
	assert.Equal(t, *s.Slices[3], AddrSlice{24,32, 0xff000000, 256})
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
		{"", args{0,1}, AddrSlice{0,1, 0x1, 2}},
		{"", args{0,0}, AddrSlice{0,0, 0x0, 1}},
		{"", args{0,5}, AddrSlice{0,5, 0x1f, 32}},
		{"", args{5,16}, AddrSlice{5,16, 0xffe0, 2048}},
		{"", args{16,24}, AddrSlice{16,24, 0xff0000, 256}},
		{"", args{24,32}, AddrSlice{24,32, 0xff000000, 256}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAddrSlice(tt.args.begin, tt.args.end); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("NewAddrSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}