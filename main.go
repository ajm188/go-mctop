package main

import (
	"flag"
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

	if err := mc.Run(); err != nil {
		panic(err)
	}
}
