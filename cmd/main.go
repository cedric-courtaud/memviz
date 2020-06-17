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
		return errors.New("Not yet implemented")
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

	flag.Parse()

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
