package wol

import (
	"fmt"
	"net"
	"strings"
)

const (
	defaultPort = "9"
	repeatCount = 16
)

// Send broadcasts a Wake-on-LAN magic packet for mac to the IPv4 broadcast address.
func Send(mac string) error {
	hw, err := parseMAC(mac)
	if err != nil {
		return err
	}
	packet := buildMagicPacket(hw)

	broadcast := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 9,
	}
	conn, err := net.DialUDP("udp4", nil, broadcast)
	if err != nil {
		return fmt.Errorf("udp dial: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(packet); err != nil {
		return fmt.Errorf("send packet: %w", err)
	}
	return nil
}

func parseMAC(s string) (net.HardwareAddr, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", ":")
	hw, err := net.ParseMAC(s)
	if err != nil {
		return nil, fmt.Errorf("invalid mac %q: %w", s, err)
	}
	if len(hw) != 6 {
		return nil, fmt.Errorf("mac must be 6 bytes, got %d", len(hw))
	}
	return hw, nil
}

// BuildMagicPacket returns the standard 102-byte WoL payload (for tests).
func BuildMagicPacket(hw net.HardwareAddr) []byte {
	return buildMagicPacket(hw)
}

func buildMagicPacket(hw net.HardwareAddr) []byte {
	out := make([]byte, 6+len(hw)*repeatCount)
	for i := range 6 {
		out[i] = 0xff
	}
	for i := 0; i < repeatCount; i++ {
		copy(out[6+i*len(hw):], hw)
	}
	return out
}

// DefaultPort is the usual UDP port used for WoL (informational).
func DefaultPort() string { return defaultPort }
