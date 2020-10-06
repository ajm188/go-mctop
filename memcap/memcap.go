package memcap

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var ValueRegexp = regexp.MustCompile(`VALUE (\S+) \S+ (\S+)`)

type Memcap struct {
	iface string
	port  int

	d time.Duration

	filter string
}

func NewMemcap(iface string, port int, d time.Duration) (*Memcap, error) {
	return &Memcap{
		iface: iface,
		port:  port,

		d: d,

		filter: fmt.Sprintf("port %d", port),
	}, nil
}

func (mc *Memcap) Run() error {
	handle, err := pcap.OpenLive(mc.iface, 1600, true, pcap.BlockForever)
	if err != nil {
		return err
	}

	if err := handle.SetBPFFilter(mc.filter); err != nil {
		return err
	}

	var (
		// eth layers.Ethernet
		// ip4 layers.IPv4
		tcp layers.TCP
	)

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeTCP, &tcp)
	decoded := []gopacket.LayerType{}

	errCh := make(chan error)

	go func() {
		src := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range src.Packets() {
			fmt.Printf("%s\n", packet.Dump())

			if err := parser.DecodeLayers(packet.Data(), &decoded); err != nil {
				errCh <- err
				return
			}

			for _, layerType := range decoded {
				switch layerType {
				case layers.LayerTypeTCP:
					fmt.Printf("%s\n", tcp.Payload)
				}
			}
		}
	}()

	select {
	case errCh <- err:
		return err
	case <-time.After(mc.d):
		fmt.Printf("Finished after %s\n", mc.d.String())
	}

	return nil
}
