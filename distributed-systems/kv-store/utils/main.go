package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type HH struct {
	Handle string
	Host   string
}

type UserRecord struct {
	HH
	Did string
}

// {"did":"did:plc:..."}
type DidResponse struct {
	Did string
}

// FetchDid fetches the DID for a given handle and host.
// Currently, the host is hardcoded to "bsky.social" where
// this function is called.
func FetchDid(handle, host string) string {
	url := fmt.Sprintf("https://%s/xrpc/com.atproto.identity.resolveHandle?handle=%s.%s", host, handle, host)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err in request to", host, "with", handle, ":", err)
	}
	defer resp.Body.Close()
	// Check status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("status code error:", resp.StatusCode, resp.Status)
	}
	// Decode response body
	var didResponse DidResponse
	err = json.NewDecoder(resp.Body).Decode(&didResponse)
	if err != nil {
		fmt.Println("err in decoding response body:", err)
	}
	return didResponse.Did
}

func InputToSetPieces(input []string) UserRecord {
	// Dragons! Hardcoded host until federation is implemented.
	host := "bsky.social"

	return UserRecord{
		HH:  HH{input[0], host},
		Did: FetchDid(input[0], host),
	}
}

// ValidateSet ensures we have a Twitter handle and an AT proto handle. Or at least, strings for each of those.
// Dragons! We do not check for a host, but we will need to do so when federation is implemented.
func ValidateSet(parts []string) error {
	if len(parts) != 2 {
		return errors.New("usage: `set tw-handle handle`")
	}

	twHandle, handle := parts[0], parts[1]

	if twHandle == "" || handle == "" {
		return errors.New("usage: `set tw-handle handle`")
	}

	return nil
}

func WithSpace(inputs []string) string {
	return strings.Join(inputs, " ")
}
