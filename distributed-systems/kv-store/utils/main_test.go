package utils

import (
	"testing"
)

func TestEncode(t *testing.T) {
	u := UserRecord{
		HH: HH{
			Handle: "chrisbodhi",
			Host:   "github.com",
		},
		Did: "did:plc:qwerty123456789",
	}

	decoded := Decode(Encode(u))

	if decoded.HH.Handle != u.HH.Handle {
		t.Errorf("expected %s, got %s", u.HH.Handle, decoded.HH.Handle)
	}

	if decoded.HH.Host != u.HH.Host {
		t.Errorf("expected %s, got %s", u.HH.Host, decoded.HH.Host)
	}

	if decoded.Did != u.Did {
		t.Errorf("expected %s, got %s", u.Did, decoded.Did)
	}
}
