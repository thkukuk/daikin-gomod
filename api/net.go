package daikin

import (
	"fmt"
	"net"
	"time"

	log "github.com/thkukuk/mqtt-exporter/pkg/logger"
)

var wantFlags = net.FlagUp | net.FlagBroadcast | net.FlagMulticast

const (
	udpQueryPayload = "DAIKIN_UDP/common/basic_info"
)

// Option is an option type to pass to NewNetwork.
type Option func(*DaikinNetwork)

// InterfaceOption configures a specific interface to scan on.
func InterfaceOption(i string) func(*DaikinNetwork) {
	return func(d *DaikinNetwork) {
		d.Interface = i
	}
}

// AddressOption specifies a specific address to query.
func AddressOption(addr string) func(*DaikinNetwork) {
	return func(d *DaikinNetwork) {
		if addr != "" {
			d.Devices = map[string]*Daikin{
				addr: &Daikin{Address: addr},
			}
			d.PollCount = 0
		}
	}
}

// DebugOption configures debug logging
func DebugOption(i bool) func(*DaikinNetwork) {
	return func(d *DaikinNetwork) {
		d.verbose = i
	}
}

// NewNetwork returns a new DaikinNetwork, attached to the given interface.
func NewNetwork(o ...Option) (*DaikinNetwork, error) {
	dn := &DaikinNetwork{
		PollInterval: time.Second,
		PollCount:    1,
		Devices:      map[string]*Daikin{},
	}
	for _, opt := range o {
		opt(dn)
	}
	return dn, nil
}

// A DaikinNetwork represents a local network with Daikin device(s).
type DaikinNetwork struct {
	// Interface is the name of the local network interface.
	Interface string

	// PollInterval is the interval to poll for Daikin devices.
	PollInterval time.Duration
	// PollCount is the number of times to poll for Daikin devices.
	PollCount int

	// Devices are the Daikin devices found on the DaikinNetwork.
	Devices map[string]*Daikin

	broadcasts []net.IP

	verbose bool
}

// getBroadcastAddresses fetches and populates the interface broadcast addresses.
func (d *DaikinNetwork) getBroadcastAddresses() error {
	d.broadcasts = []net.IP{}
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range interfaces {
		if i.Flags != wantFlags || d.Interface != "" && i.Name != d.Interface {
			continue
		}
		// Fetch interface addresses.
		adr, err := i.Addrs()
		if err != nil {
			log.Warnf("%s: Can't get addresses, skipping.", i.Name)
			continue
		}
		for _, a := range adr {
			// Parse the address.
			ip, network, err := net.ParseCIDR(a.String())
			if err != nil {
				log.Infof("%s: Can't parse %s, skipping.", i.Name, a.String())
				continue
			}
			// Test if it is V4 (no daikin does ipv6).
			if four := ip.To4(); four == nil {
				if d.verbose {
					log.Debugf("%s: %s: Skipping non-v4 address", i.Name, ip)
				}
				continue
			}
			// Calculate and add the broadcast address.
			bCast := net.IP{0, 0, 0, 0}
			for i := 0; i < 4; i++ {
				bCast[i] = byte(network.IP[i]) | (0xff - network.Mask[i])
			}
			d.broadcasts = append(d.broadcasts, bCast)
		}
	}
	if len(d.broadcasts) == 0 && d.Interface != "" {
		return fmt.Errorf("no interface or no addresses: %s", d.Interface)
	}
	if d.verbose {
		log.Debugf("Broadcast addresses: %v", d.broadcasts)
	}
	return nil
}

// Discover runs a UDP polling cycle for Daikin devices.
// Sends UDP packet to broadcast address, dst port 30050 with payload:
// DAIKIN_UDP/common/basic_info
func (d *DaikinNetwork) Discover() error {
	if d.PollCount < 1 {
		return nil
	}
	if err := d.getBroadcastAddresses(); err != nil {
		return err
	}
	// Open a local listener.
	lAddr := net.UDPAddr{Port: 30000}
	conn, err := net.ListenUDP("udp", &lAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// A poller sends to broadcast and awaits replies.
	poller := func(bCast string, done chan bool) {
		if d.verbose {
			log.Debugf("Start polling to: %s", bCast)
		}
		for i := 0; i < d.PollCount; i++ {
			// Send broadcast packet.
			rAddr := &net.UDPAddr{IP: net.ParseIP(bCast), Port: 30050}
			if _, err := conn.WriteToUDP([]byte(udpQueryPayload), rAddr); err != nil {
				log.Errorf("write: err: %v\n", err)
				continue
			}
			// Read until the deadline.
			for {
				rBuf := make([]byte, 2048)
				conn.SetReadDeadline(time.Now().Add(d.PollInterval))
				n, rAddr, err := conn.ReadFromUDP(rBuf)
				if err != nil {
					if err, ok := err.(net.Error); ok && err.Timeout() {
						break
					}
					log.Errorf("read err: %v\n", err)
					continue
				}
				if d.verbose {
					log.Debugf("%d bytes from %v: %v\n", n, rAddr, string(rBuf))
				}

				ip := rAddr.IP.String()
				if _, ok := d.Devices[ip]; !ok {
					dev := &Daikin{Address: ip}
					d.Devices[ip] = dev
				}
			}
		}
		close(done)
	}

	// Start pollers per broadcast address, wait for them to complete.
	pollers := []chan bool{}
	for _, b := range d.broadcasts {
		ch := make(chan bool)
		go poller(b.String(), ch)
		pollers = append(pollers, ch)
	}
	for _, ch := range pollers {
		_, _ = <-ch
	}

	return nil
}
