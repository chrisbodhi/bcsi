package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type UserRecord struct {
	Handle string
	Host   string
	Did    string
}

// {"did":"did:plc:..."}
type DidResponse struct {
	Did string
}

const (
	FIELD_SIZE_LENGTH = 2
	STORAGE_BASE      = "storage.json"
)

var mem = make(map[string]map[string][]byte)

func Decode(bytes []byte) UserRecord {
	user := UserRecord{}
	start, end := 0, FIELD_SIZE_LENGTH

	for _, key := range []string{"handle", "host", "did"} {
		snip := string(bytes[start:end])
		sansZeroes := []byte{}
		for _, b := range snip {
			if b != 0 {
				sansZeroes = append(sansZeroes, byte(b))
			}
		}
		// If the second byte is a 0, then Atoi errors. Hence, the need
		// to remove the 0s.
		fieldLen, err := strconv.Atoi(string(sansZeroes))
		if err != nil {
			fmt.Println("err in decoding field length:", err)
		}
		field := string(bytes[end : end+fieldLen])
		start = end + fieldLen
		end = start + FIELD_SIZE_LENGTH

		switch key {
		case "handle":
			user.Handle = field
		case "host":
			user.Host = field
		case "did":
			user.Did = field
		}
	}

	return user
}

func Encode(user UserRecord) []byte {
	var buffer []byte

	handle := user.Handle
	host := user.Host
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

func Get(key, table string) ([]byte, bool) {
	value, ok := mem[table][key]
	return value, ok
}

func Drop(table, port string) (string, error) {
	// Remove from mem
	if _, ok := mem[table]; ok {
		delete(mem, table)
	} else {
		msg := fmt.Sprintf("%s does not exist", table)
		return "", errors.New(msg)
	}

	// Rename backing datastore/file
	storage := fmt.Sprintf("%s_%s_%s", port, table, STORAGE_BASE)
	if err := os.Rename(storage, fmt.Sprintf("dropped_%s", table)); err != nil {
		msg := fmt.Sprintf("Failed to removing backing datastore for %s", table)
		return msg, nil
	}
	msg := fmt.Sprintf("Removed %s", table)
	return msg, nil
}

func Set(key string, value UserRecord, table, port string) {
	if _, ok := mem[table]; !ok {
		// Flush mem to {port}_{table}_storage.json
		mem[table] = make(map[string][]byte)
	}
	mem[table][key] = Encode(value)
	UpdateDatastore(table, port)
}

func InputToSetPieces(input []string) UserRecord {
	// Dragons! Hardcoded host until federation is implemented.
	host := "bsky.social"

	return UserRecord{
		Handle: input[0],
		Host:   host,
		Did:    FetchDid(input[0], host),
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

// Random returns a random number between x and y, inclusive.
func Random(x, y int) int {
	return x + rand.Intn(y-x)
}

// WriteAhead takes a command and writes it to a file.
func WriteAhead(line string) {
	filename := "ahead.txt"

	// Create write-ahead log if it doesn't exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			fmt.Println("err in creating file:", err)
		}
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("err in opening file:", err)
	}

	defer f.Close()

	if _, err := f.WriteString(line + "\n"); err != nil {
		fmt.Println("err in writing to file:", err)
	}
}

func LoadDatastore(port, table string) {
	storage := fmt.Sprintf("%s_%s_%s", port, table, STORAGE_BASE)
	// Create storage file if it doesn't exist
	if _, err := os.Stat(storage); os.IsNotExist(err) {
		_, err := os.Create(storage)
		if err != nil {
			fmt.Println(err)
		}
	}
	jsonFile, err := os.Open(storage)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	localTable := make(map[string][]byte)
	json.Unmarshal(byteValue, &localTable)
	mem[table] = localTable
}

func UpdateDatastore(table, port string) {
	storage := fmt.Sprintf("%s_%s_%s", port, table, STORAGE_BASE)

	// Create storage file if it doesn't exist
	if _, err := os.Stat(storage); os.IsNotExist(err) {
		_, err := os.Create(storage)
		if err != nil {
			fmt.Println(err)
		}
	}
	jsonFile, err := os.OpenFile(storage, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	jsonData, err := json.Marshal(mem[table])
	if err != nil {
		fmt.Println(err)
	}
	jsonFile.Write(jsonData)
}
