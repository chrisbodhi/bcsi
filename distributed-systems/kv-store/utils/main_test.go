package utils

import (
	"testing"
)

func TestDecode(t *testing.T) {
	var u UserRecord

	encoded := []byte{57, 0, 108, 101, 97, 104, 112, 114, 105, 109, 101, 49, 49, 98, 115, 107, 121, 46, 115, 111, 99, 105, 97, 108, 51, 50, 100, 105, 100, 58, 112, 108, 99, 58, 104, 121, 98, 121, 108, 113, 104, 104, 106, 52, 102, 51, 104, 121, 119, 51, 122, 118, 112, 116, 120, 101, 106, 97}

	u = Decode(encoded)
	if u.Handle != "leahprime" {
		t.Errorf("expected %s, got %s", "leahprime", u.Handle)
	}
	if u.Host != "bsky.social" {
		t.Errorf("expected %s, got %s", "bsky.social", u.Host)
	}
	if u.Did != "did:plc:hybylqhhj4f3hyw3zvptxeja" {
		t.Errorf("expected %s, got %s", "did:plc:hybylqhhj4f3hyw3zvptxeja", u.Did)
	}
}

func TestEncode(t *testing.T) {
	u := UserRecord{
		Handle: "chrisbodhi",
		Host:   "bluesky.social",
		Did:    "did:plc:qwerty123456789",
	}

	decoded := Decode(Encode(u))

	if decoded.Handle != u.Handle {
		t.Errorf("expected %s, got %s", u.Handle, decoded.Handle)
	}

	if decoded.Host != u.Host {
		t.Errorf("expected %s, got %s", u.Host, decoded.Host)
	}

	if decoded.Did != u.Did {
		t.Errorf("expected %s, got %s", u.Did, decoded.Did)
	}
}
