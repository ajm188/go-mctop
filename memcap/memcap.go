package memcap

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var MemcCommandRegexp = regexp.MustCompile(`^(?P<command>[a-z]+) (?P<key>\S+)`)

type Memcap struct {
	iface string
	port  int

	d time.Duration

	filter string

	stats *Stats
}

func NewMemcap(iface string, port int, d time.Duration) (*Memcap, error) {
	return &Memcap{
		iface: iface,
		port:  port,

		d: d,

		filter: fmt.Sprintf("tcp and port %d", port),

		stats: NewStats(),
	}, nil
}

func (mc *Memcap) GetStats() *Stats {
	return mc.stats
}

func (mc *Memcap) Run(done chan bool) error {
	handle, err := pcap.OpenLive(mc.iface, 3200, true, pcap.BlockForever)
	if err != nil {
		return err
	}

	if err := handle.SetBPFFilter(mc.filter); err != nil {
		return err
	}

	errCh := make(chan error)

	packets := 0

	go func() {
		var (
			eth     layers.Ethernet
			ip4     layers.IPv4
			tcp     layers.TCP
			payload gopacket.Payload
		)

		parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &tcp, &payload)
		decoded := []gopacket.LayerType{}

		src := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range src.Packets() {
			select {
			case <-done:
				return
			default:
			}
			//fmt.Printf("%s\n", packet.Dump())

			// for _, layer := range packet.Layers() {
			// 	fmt.Printf("- %s\n", layer.LayerType())
			// }

			if err := parser.DecodeLayers(packet.Data(), &decoded); err != nil {
				errCh <- err
				return
			}

			for _, layerType := range decoded {
				switch layerType {
				case layers.LayerTypeEthernet:
				case layers.LayerTypeIPv4:
				case layers.LayerTypeTCP:
				case gopacket.LayerTypePayload:
					packets++

					m := MemcCommandRegexp.FindStringSubmatch(string(payload.Payload()))
					if m == nil {
						continue
					}

					command := m[1]
					key := m[2]

					mc.stats.Add(key, command)
				}
			}
		}
	}()

	select {
	case errCh <- err:
		return err
	case <-time.After(mc.d):
		fmt.Printf("Finished after %s. Processed %d packets\n", mc.d.String(), packets)
	}

	done <- true
	// fmt.Printf("%+v\n", stats)

	return nil
}
