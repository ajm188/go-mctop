package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ajm188/go-mctop/memcap"
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

	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		<-quit
		done <- true
	}()

	if err := mc.Run(done); err != nil {
		panic(err)
	}

	fmt.Printf("%v %d\n", mc.GetStats().Calls(), mc.GetStats().TotalCalls())
}
