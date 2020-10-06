package memcap

import (
	"time"

	"github.com/google/gopacket/pcapgo"
)

type Memcap struct {
	iface string
	port  int

	d time.Duration
}

func NewMemcap(iface string) *Memcap {
	return &Memcap{
		iface: iface,
	}
}

func (mc *Memcap) Run() error {
	_, err := pcapgo.NewEthernetHandle(mc.iface)
	if err != nil {
		return err
	}

	return nil
}
