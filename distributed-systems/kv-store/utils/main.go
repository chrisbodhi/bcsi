package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

const FIELD_SIZE_LENGTH = 2

func Decode(bytes []byte) UserRecord {
	u := UserRecord{}

	handleLen, err := strconv.Atoi(string(bytes[:FIELD_SIZE_LENGTH]))
	fmt.Println("handleLen", handleLen)
	if err != nil {
		fmt.Println("err in decoding handle length:", err)
	}
	u.HH.Handle = string(bytes[FIELD_SIZE_LENGTH : handleLen+FIELD_SIZE_LENGTH])

	hostLen, err := strconv.Atoi(string(bytes[handleLen+FIELD_SIZE_LENGTH : handleLen+FIELD_SIZE_LENGTH+FIELD_SIZE_LENGTH]))

	if err != nil {
		fmt.Println("err in decoding host length:", err)
	}

	u.HH.Host = string(bytes[handleLen+FIELD_SIZE_LENGTH+FIELD_SIZE_LENGTH : handleLen+FIELD_SIZE_LENGTH+FIELD_SIZE_LENGTH+hostLen])

	didLen, err := strconv.Atoi(string(bytes[handleLen+FIELD_SIZE_LENGTH+FIELD_SIZE_LENGTH+hostLen : handleLen+FIELD_SIZE_LENGTH+FIELD_SIZE_LENGTH+hostLen+FIELD_SIZE_LENGTH]))

	if err != nil {
		fmt.Println("err in decoding did length:", err)
	}

	u.Did = string(bytes[handleLen+FIELD_SIZE_LENGTH+FIELD_SIZE_LENGTH+hostLen+FIELD_SIZE_LENGTH : handleLen+FIELD_SIZE_LENGTH+FIELD_SIZE_LENGTH+hostLen+FIELD_SIZE_LENGTH+didLen])

	return u
}

func Encode(user UserRecord) []byte {
	var buffer []byte

	handle := user.HH.Handle
	host := user.HH.Host
	did := user.Did

	for _, field := range []string{handle, host, did} {
		fieldLen, fieldBytes := getBytesForEncoding(user, field)
		buffer = append(buffer, fieldLen...)
		buffer = append(buffer, fieldBytes...)
	}

	return buffer
}

func getBytesForEncoding(user UserRecord, field string) ([]byte, []byte) {
	fieldValueLen := len(field)

	fieldValueLenBytes := make([]byte, FIELD_SIZE_LENGTH)
	copy(fieldValueLenBytes, []byte(strconv.Itoa(fieldValueLen)))

	return fieldValueLenBytes, []byte(field)
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

	if resp.StatusCode != http.StatusOK {
		fmt.Println("status code error:", resp.StatusCode, resp.Status)
	}

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
