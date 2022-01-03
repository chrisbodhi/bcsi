package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"
	// "encoding/binary"
)

const LITTLE_ENDIAN_MAGIC = "d4c3b2a1"

type CaptureHeader struct {
	magicNumber 	uint32
	majorVersion	uint16
	minorVersion	uint16
	timezoneOffset	uint32
	timestampAcc	uint32
	snapshotLen		uint32
	headerType		uint32
}

type PacketHeader struct {
	TimestampSec	uint32
	TimestampMicro	uint32 // could be nano for other cases
	CaptureLen		uint32
	UntruncLen		uint32
}

var ch CaptureHeader

func main() {
	path := "net.cap"

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Dang, ", err)
	}

	defer file.Close()
	magic, _ := readNextBytes(file, unsafe.Sizeof(ch.magicNumber), 0)
	checkEndian(magic)
	fmt.Println("opened!")
	startingPoint := int64(unsafe.Sizeof(ch))
	sum := countCaptured(file, startingPoint)
	fmt.Println("summed up, got", sum)
}

func readNextBytes(file *os.File, number uintptr, offset int64) ([]byte, error) {
	bytes := make([]byte, number)

	_, err := file.ReadAt(bytes, offset)
	if err != nil {
		if err == io.EOF {
			return nil, err
		} else {
			log.Fatal(err)
		}
	}

	return bytes, nil
}

func checkEndian(magic []byte) {
	magicHex := fmt.Sprintf("%x", magic)
	if magicHex != LITTLE_ENDIAN_MAGIC {
		log.Fatal("Wrong magic number: ", magicHex)
	}
}

func allCaptured(p PacketHeader) bool {
    return p.CaptureLen == p.UntruncLen
}

func countCaptured(file *os.File, from int64) int64 {
	var sum int64
	return countHelper(file, from, sum)
}

func countHelper(file *os.File, from int64, sum int64) int64 {
	var ph PacketHeader
	phSize := unsafe.Sizeof(ph)
	p := PacketHeader{}

	// get packet header bytes
	phBytes, err := readNextBytes(file, phSize, from)

	if err != nil {
		log.Println("Done with countHelper", err)
		return sum
	}
	
	buffer := bytes.NewBuffer(phBytes)
	err = binary.Read(buffer, binary.LittleEndian, &p)
	if err != nil {
		log.Fatal("cannot read bytes into struct")
	}
	fmt.Printf("Parsed data: %+v\n", p)
	
	if allCaptured(p) {
		capLen := int64(p.CaptureLen)
		return countHelper(file, from + int64(phSize) + capLen, sum + 1)
	} else {
		log.Println("Bytes truncated")
		return sum
	}

}
