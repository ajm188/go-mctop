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
		ticker := time.Tick(time.Second * 3)

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
