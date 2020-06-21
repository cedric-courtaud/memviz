package internal

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"memrec/internal/flatbuffers"
)

const DEFAULT_IMAGE_GENERATOR_CAPACITY = 2048

type AccessCount struct {
	NI int
	NR int
	NW int
}

type CheckpointPos struct {
	Checkpoint *Checkpoint
	XPos		int
}

type ImageGenerator struct {
	Width                int
	Height               int
	InstructionPerColumn int
	AddrSlicing          *AddrSlicing
	buff                 [][]AccessCount
	lastAccess           *Access
	instBefore           uint64
	yPositioners         []*YPositioner
	yPos                 []int
	checkpoints			 []CheckpointPos
	Beta			     float32
	Writer io.Writer
}

func (i * ImageGenerator) pushNewRow() {
	row := make([]AccessCount, i.Height)
	i.buff = append(i.buff, row)
	i.Width += 1
}

func (i * ImageGenerator) getXPos(progress uint64) int {

	return int(progress) / i.InstructionPerColumn
}

func (i * ImageGenerator) getYPos(a * Access, buff []int) {
	addr := a.DestAddr

	for j, positioner := range i.yPositioners {
		slice := positioner.slice
		a1 := (addr & slice.mask) >> slice.begin
		a2 := (i.lastAccess.DestAddr & slice.mask) >> slice.begin
		diff := (a1 - a2) % (slice.nVal)
		buff[j] = positioner.GetYPos(diff)
	}
}

type YPositioner struct {
	slice * AddrSlice
	LayerHeight int
	Offset int
}

func NewYPositioner(slice * AddrSlice, totalHeight, totalAddrBits int) *YPositioner {
	hLayer := int((float32(slice.end - slice.begin) / float32(totalAddrBits)) * float32(totalHeight))
	offset := int((float32(slice.begin) / float32(totalAddrBits)) * float32(totalHeight))

	return &YPositioner{
		slice: slice,
		LayerHeight: hLayer,
		Offset: offset,
	}
}

func (p * YPositioner) GetYPos(value uint64) int {
	return int((float64(value) / float64(p.slice.nVal)) * float64(p.LayerHeight)) + p.Offset
}

func NewImageGenerator(Height, InstructionPerColumn int, slicing *AddrSlicing) *ImageGenerator {
	buff := make([][]AccessCount, 0, DEFAULT_IMAGE_GENERATOR_CAPACITY)
	positioners := make([]*YPositioner, 0, len(slicing.Slices))

	for _, slice := range slicing.Slices {
		positioners = append(positioners, NewYPositioner(slice, Height, slicing.Total))
	}

	ret := ImageGenerator{
		Width:                0,
		Height:               Height,
		InstructionPerColumn: InstructionPerColumn,
		AddrSlicing:          slicing,
		buff:                 buff,
		lastAccess:           nil,
		yPositioners:         positioners,
		yPos:                 make([]int, len(positioners)),
		Beta: 50.0,
	}

	ret.pushNewRow()

	return &ret
}

func (i *ImageGenerator) HandleAccess(access *Access) error {
	i.instBefore += access.InstBefore

	if i.lastAccess == nil {
		i.lastAccess = access
		return nil
	}

	x := i.getXPos(i.instBefore)
	if x >= i.Width {
		for i.Width <= x {
			i.pushNewRow()
		}
	}

	i.getYPos(access, i.yPos)

	for _, y := range i.yPos {
		switch access.AccessType {
		case flatbuffers.AccessTypeI:
			i.buff[x][y].NI += 1
			break
		case flatbuffers.AccessTypeR:
			i.buff[x][y].NR += 1
			break
		case flatbuffers.AccessTypeW:
			i.buff[x][y].NW += 1
			break
		}
	}

	i.lastAccess = access

	return nil
}

func (i *ImageGenerator) HandleCheckpoint(checkpoint *Checkpoint) error {
	progress := i.instBefore + checkpoint.InstBefore

	c := CheckpointPos{checkpoint, i.getXPos(progress)}
	i.checkpoints = append(i.checkpoints, c)

	return nil
}

func (i * ImageGenerator) GenerateImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, i.Width, i.Height))
	ipc := float32(i.InstructionPerColumn)

	// Draw accesses
	for x := 0; x < i.Width; x++ {
		for y := 0; y < i.Height; y++ {
			acc := i.buff[x][y]
			t := float32(acc.NI + acc.NR + acc.NW)

			r := uint8((float32(acc.NW)/ t) * 255.0)
			b := uint8((float32(acc.NR)/ t) * 255.0)
			g := uint8((float32(acc.NI)/ t) * 255.0)
			a := uint8(i.Beta * (t / ipc) * 255.0)

			px := color.RGBA{r,g,b,a}
			img.Set(x, i.Height - y, px)
		}
	}

	// Draw checkpoints
	for _, c := range i.checkpoints {
		x := c.XPos
		for y := 0; y < i.Height; y++ {
			img.Set(x, y, color.RGBA{0,0,0, 127})
		}
	}

	return img
}

func (i *ImageGenerator) Start() {
}

func (i *ImageGenerator) Stop() {
}

func (i *ImageGenerator) Finalize() {
	img := i.GenerateImage()
	encoder := png.Encoder{
		CompressionLevel: -1,
	}

	encoder.Encode(i.Writer, img)
}
