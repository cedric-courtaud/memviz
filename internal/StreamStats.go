package internal

import (
	"errors"
	"fmt"
	"github.com/cedric-courtaud/memviz/internal/flatbuffers"
	"io"
	"math"
	"text/tabwriter"
)

type StatsConfig struct {
	AddrSlicing *AddrSlicing
}

func NewStatsConfig(sliceSpec string) (*StatsConfig, error) {
	s, err := ParseAddrSlicing(sliceSpec)

	if err != nil {
		return nil, err
	}

	return &StatsConfig{s}, nil
}

type phaseStats struct {
	Id             string
	AccessCount    uint64
	InversionCount uint64
	previousAccess *Access
	addrSliceSpec  *AddrSlicing
	addrDiffCount  map[uint64]map[uint64]uint64
}

func newPhaseStats(Id string, conf *StatsConfig) phaseStats {
	diff := make(map[uint64]map[uint64]uint64)

	for _, slice := range conf.AddrSlicing.Slices {
		diff[slice.mask] = make(map[uint64]uint64)
	}

	return phaseStats{
		Id:             Id,
		AccessCount:    0,
		InversionCount: 0,
		previousAccess: nil,
		addrSliceSpec:  conf.AddrSlicing,
		addrDiffCount:  diff,
	}
}

type StreamStats struct {
	Config     *StatsConfig
	phaseStats []phaseStats
	Writer     io.Writer
}

func (s *StreamStats) Finalize() {
	s.WriteSummary(s.Writer)
}

func (s *StreamStats) Start() {}

func (s *StreamStats) Stop() {
	/*
		if (s.Writer != nil) {
			s.WriteSummary(s.Writer)
		}*/
}

func (p *phaseStats) updateCount() {
	p.AccessCount += 1
}

func (p *phaseStats) updatePreviousAccess(access *Access) {
	p.previousAccess = access
}

func isRead(access Access) bool {
	return access.AccessType != flatbuffers.AccessTypeW
}

func (p *phaseStats) updateInversionCount(access *Access) {
	if p.previousAccess != nil {
		if isRead(*access) != isRead(*p.previousAccess) {
			p.InversionCount += 1
		}
	}
}

func ShannonEntropy(dist map[uint64]uint64, total uint64) float64 {
	var ret float64 = 0.0
	var n float64 = float64(total)

	for _, v := range dist {
		p := float64(v) / n
		ret += -p * math.Log2(p)
	}

	return ret
}

func (p *phaseStats) updateAddressDiffCount(access *Access) {
	if p.previousAccess == nil {
		return
	}

	for _, slice := range p.addrSliceSpec.Slices {
		m := p.addrDiffCount[slice.mask]
		a := (p.previousAccess.DestAddr & slice.mask) >> slice.begin
		b := (access.DestAddr) >> slice.begin
		r := (b - a) % ((slice.mask >> slice.begin) + 1)
		m[r] += 1
	}
}

func (p *phaseStats) handleAccess(access *Access) {
	p.updateCount()
	p.updateInversionCount(access)
	p.updateAddressDiffCount(access)

	p.updatePreviousAccess(access)
}

func NewStreamStats(config *StatsConfig) *StreamStats {
	return &StreamStats{Config: config}
}

func (s StreamStats) getCurrentPhase() (*phaseStats, error) {
	if len(s.phaseStats) == 0 {
		return nil, errors.New("Access does not belong to any phase")
	}

	return &s.phaseStats[len(s.phaseStats)-1], nil
}

func (s *StreamStats) HandleAccess(access *Access) error {
	p, err := s.getCurrentPhase()
	if err != nil {
		return err
	}

	p.handleAccess(access)
	return nil
}

func (s *StreamStats) HandleCheckpoint(checkpoint *Checkpoint) error {
	p := newPhaseStats(checkpoint.Id, s.Config)

	lastPhase, err := s.getCurrentPhase()

	if err == nil {
		p.updatePreviousAccess(lastPhase.previousAccess)
	}

	s.phaseStats = append(s.phaseStats, p)

	return nil
}

func (s *StreamStats) HandleForked(forked *Forked) error {
	return nil
}

func (s StreamStats) WriteSummary(writer io.Writer) error {
	w := new(tabwriter.Writer)
	w.Init(writer, 0, 0, 2, ' ', tabwriter.AlignRight)
	defer w.Flush()

	if len(s.phaseStats) > 0 {
		s.phaseStats[0].WriteSummaryHeader(w)
	}

	for _, phase := range s.phaseStats {
		phase.WriteSummary(w)
	}

	return nil
}

func (p phaseStats) WriteSummaryHeader(writer io.Writer) {
	fmt.Fprintf(writer, "%s\t%s\t%s", "Phase", "NAccess", "NInv")

	for _, s := range p.addrSliceSpec.Slices {
		s := fmt.Sprintf("HdAddr[%d:%d]", s.begin, s.end)
		fmt.Fprintf(writer, "\t%s", s)
	}

	fmt.Fprintf(writer, "\t\n%s\t%s\t%s", "---", "---", "---")

	for _, _ = range p.addrSliceSpec.Slices {
		fmt.Fprintf(writer, "\t%s", "---")
	}

	fmt.Fprintf(writer, "\t\n")
}

func (p phaseStats) WriteSummary(writer io.Writer) {
	fmt.Fprintf(writer, "%s\t%d\t%d", p.Id, p.AccessCount, p.InversionCount)
	for i, s := range p.addrSliceSpec.Slices {
		total := p.AccessCount
		if i == 0 {
			total -= 1
		}
		m := p.addrDiffCount[s.mask]
		fmt.Fprintf(writer, "\t%f", ShannonEntropy(m, total))
	}

	fmt.Fprintf(writer, "\t\n")
}
