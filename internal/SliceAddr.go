package internal;

import (
	"strconv"
	"strings"
)

func setBits(val uint64 , n uint, offset uint) uint64 {
	var mask uint64 = 0
	var i uint
	for i = 0; i < n; i += 1 {
		mask |= 1 << i
	}

	mask <<= offset

	return val | mask
}

type AddrSlice struct {
	begin uint
	end uint
	mask uint64
}

type AddrSlicing struct {
	Slices []*AddrSlice
}

func NewAddrSlice (begin, end uint) *AddrSlice {
	return &AddrSlice {
		begin,
		end,
		setBits(0, end - begin, begin),
	}
}

func ParseAddrSlicing(specStr string) (*AddrSlicing, error) {
	tokens := strings.Split(specStr, ":")
	lengths := make([]uint, len(tokens))

	for i, token := range tokens {
		length, err := strconv.ParseUint(token, 10, 8)
		if err != nil {
			return nil, err
		}
		lengths[i] = uint(length)
	}

	var slices []*AddrSlice

	var begin uint
	begin = 0
	for i := len(lengths) - 1; i >= 0; i -= 1 {
		slices = append(slices, NewAddrSlice(begin, begin + lengths[i]))
		begin += lengths[i]
	}

	return &AddrSlicing{slices}, nil
}


