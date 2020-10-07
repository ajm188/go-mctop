package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ajm188/go-mctop/memcap"
	"github.com/ajm188/go-mctop/ui"
)

var (
	flags = &flag.FlagSet{}

	iface = flags.String("interface", "", "interface to sniff")
	port  = flags.Int("port", 11211, "port to sniff")
	t     = flags.Duration("duration", time.Minute, "how long to sniff")
)

func main() {
	flag.CommandLine = flags
	flag.Parse()

	mc, err := memcap.NewMemcap(*iface, *port, *t)
	if err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	doneSniffing := make(chan bool, 1)
	doneDrawing := make(chan bool, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		<-quit
		doneSniffing <- true
		doneDrawing <- true
	}()

	// TODO: more complex, but allow mctop to operate in two modes:
	// 1. report mode, which never renders a termui, but produces a detailed text report
	// 2. ui mode, which does all this drawing stuff. it's way more expensive to do,
	// 	  mostly due to the constant sorting of an ever-growing list of key stats, so
	//	  mctop should let you get stats without forcing you to pay that cost.
	//
	// Thought: if you include a -duration flag, it runs in report mode for the specified
	// duration, but if you omit that flag it runs in ui mode indefinitely.

	if err := ui.Init(doneDrawing); err != nil {
		panic(err)
	}
	defer ui.Close()

	go func() {
		if err := mc.Run(doneSniffing); err != nil {
			fmt.Println(err)
		}

		done <- true
	}()

	go func() {
		// TODO: this interval should be configurable
		ticker := time.Tick(time.Second * 3)

		// TODO: make the initial wait time configurable
		// Wait for some data to arrive, then do an initial draw.
		time.Sleep(time.Millisecond * 100)
		draw(mc.GetStats())

		for {
			select {
			case <-ticker:
				stats := mc.GetStats()
				draw(stats)
			case <-doneDrawing:
				done <- true
				return
			}
		}
	}()

	<-done
}

func draw(stats *memcap.Stats) {
	ui.UpdateCalls(stats)
	ui.UpdateKeys(stats)
	ui.Render()
}
