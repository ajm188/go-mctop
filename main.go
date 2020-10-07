package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"text/tabwriter"
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

	stats := mc.GetStats()

	tw := tabwriter.NewWriter(os.Stdout, 10, 4, 0, ' ', 0)
	fmt.Fprintln(tw, "call\tops")
	for call, count := range stats.Calls() {
		fmt.Fprintf(tw, "%s\t%d\n", call, count)
	}

	fmt.Fprintf(tw, "(total)\t%d\n", stats.TotalCalls())
	tw.Flush()

	// TODO: add flags for sorting by gets, by keysize, etc
	kss := memcap.KeyStatsList(stats.KeyStats())
	sort.Sort(sort.Reverse(kss)) // biggest first

	fmt.Fprintln(tw, "")

	fmt.Fprintln(tw, "key\tgets\tsets\tadds\tdeletes\tTOTAL")
	for _, ks := range kss {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\t%d\n", ks.Key(), ks.Gets(), ks.Sets(), ks.Adds(), ks.Deletes(), ks.TotalCalls())
	}

	tw.Flush()

	// fmt.Printf("%v %d\n", stats.Calls(), stats.TotalCalls())
}
