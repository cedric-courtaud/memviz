package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"memrec/internal"
	"os"
)

type MemRecConfig struct {
	SliceSpecString       string
	ComputeStats          bool
	StatsOutputWriter 	  io.Writer
	RecordInFlatBuffers   bool
	FlatBuffersOutputPath string
	InputPath 			  string
	InputReader 		  io.Reader
	Dispatcher			  *internal.EventDispatcher
	GenerateImage		  bool
	ImageOutputPath		  string
	Beta				  float64
	Height				  int
	Width				  int
	InstructionPerColumn  int
}

func afterCLIParseInit(config * MemRecConfig) error {
	config.Dispatcher = internal.NewEventDispatcher()

	if (config.ComputeStats) {
		sc, err := internal.NewStatsConfig(config.SliceSpecString)

		if err != nil {
			panic(err)
		}

		stats := internal.NewStreamStats(sc)
		config.StatsOutputWriter = os.Stdout
		stats.Writer = config.StatsOutputWriter

		config.Dispatcher.AddHandler(stats)
	}

	if (config.RecordInFlatBuffers) {
		f, err := os.Create(config.FlatBuffersOutputPath)
		if err != nil {
			return err
		}

		config.Dispatcher.AddHandler(internal.NewFlatBuffersRecorder(f))
	}

	if (config.GenerateImage) {
		// This line will be removed once fixed width image generation is implemented
		config.Width = 0

		if config.Width == 0 && config.InstructionPerColumn == 0 {
			return errors.New("You must set one of those: --ipc and --width")
		}

		if config.Width != 0 && config.InstructionPerColumn != 0 {
			return errors.New("You can set --ipc and --width at the same time")
		}

		slicing, err := internal.ParseAddrSlicing(config.SliceSpecString)
		if err != nil {
			return err
		}

		g := internal.NewImageGenerator(config.Height, config.InstructionPerColumn, slicing)
		if config.Beta == 0 {
			g.Beta = float32(config.InstructionPerColumn) / 5.0
		} else {
			g.Beta = float32(config.Beta)
		}

		g.Writer, err = os.Create(config.ImageOutputPath)
		if err != nil {
			return err
		}

		config.Dispatcher.AddHandler(g)
	}

	if len(config.Dispatcher.Handlers) == 0 {
		return errors.New("Nothing to do")
	}

	if config.InputPath == "" {
		config.InputReader = os.Stdin
	} else {
		f, err := os.Open(config.InputPath)
		if err != nil {
			return err
		}

		config.InputReader = f
	}



	return nil
}

func ParseCLIArgs(config * MemRecConfig) error {
	flag.StringVar(&config.FlatBuffersOutputPath, "record-to", "", "")
	flag.BoolVar(&config.ComputeStats, "stats", false, "")
	flag.StringVar(&config.SliceSpecString, "addr-slicing", "8:8:11:5", "")

	flag.StringVar(&config.ImageOutputPath, "generate-image", "", "")
	flag.Float64Var(&config.Beta, "beta", 0, "")
	// flag.IntVar(&config.Width, "width", 0, "")
	flag.IntVar(&config.InstructionPerColumn, "ipc", 0, "")
	flag.IntVar(&config.Height, "height", 1024, "")

	flag.Parse()

	if config.FlatBuffersOutputPath != "" {
		config.RecordInFlatBuffers = true
	}

	if config.ImageOutputPath != "" {
		config.GenerateImage = true
	}

	if len(flag.Args()) == 0 {
		config.InputPath = ""
	} else {
		config.InputPath = flag.Args()[0]
	}

	return afterCLIParseInit(config)
}

func Run(config * MemRecConfig) error {
	reader := bufio.NewReader(config.InputReader)
	parser := internal.NewEventParser(config.Dispatcher)
	config.Dispatcher.Start()

	err := parser.Parse(reader)
	if err != nil {
		return err
	}

	config.Dispatcher.Stop()
	config.Dispatcher.Finalize()

	return nil
}

func main() {
	config := MemRecConfig{
		SliceSpecString:       "",
		ComputeStats:          false,
		StatsOutputWriter:     nil,
		RecordInFlatBuffers:   false,
		FlatBuffersOutputPath: "",
		InputPath:             "",
		InputReader:           nil,
		Dispatcher:            nil,
	}

	err := ParseCLIArgs(&config)

	if err != nil {
		panic(err)
	}

	err = Run(&config)

	if err != nil {
		panic(err)
	}
}
