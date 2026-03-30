package wol

import (
	"net"
	"testing"
)

func TestBuildMagicPacket(t *testing.T) {
	hw := net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}
	p := BuildMagicPacket(hw)
	if len(p) != 102 {
		t.Fatalf("len got %d want 102", len(p))
	}
	for i := 0; i < 6; i++ {
		if p[i] != 0xff {
			t.Fatalf("sync byte %d: got %02x want ff", i, p[i])
		}
	}
	for r := 0; r < 16; r++ {
		off := 6 + r*6
		for j := 0; j < 6; j++ {
			if p[off+j] != hw[j] {
				t.Fatalf("repeat %d byte %d: got %02x want %02x", r, j, p[off+j], hw[j])
			}
		}
	}
}

func TestParseMAC(t *testing.T) {
	for _, s := range []string{
		"00:11:22:33:44:55",
		"00-11-22-33-44-55",
		"  aa:bb:cc:dd:ee:ff ",
	} {
		_, err := parseMAC(s)
		if err != nil {
			t.Errorf("parseMAC(%q): %v", s, err)
		}
	}
}
