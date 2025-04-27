package main

import "time"

//1. Interface composing

/*

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
)

func hashAndBroadcast(r HashReader) error {
	hash := r.hash()
	fmt.Println("Hash: ", hash)
	return broadcast(r)
}

func broadcast(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	fmt.Println("string of bytes: ", string(b))
	return nil
}

type HashReader interface {
	io.Reader
	hash() string
}

type hashReader struct {
	*bytes.Reader
	buf *bytes.Buffer
}

func NewHashReader(b []byte) *hashReader {
	return &hashReader{
		Reader: bytes.NewReader(b),
		buf:    bytes.NewBuffer(b),
	}
}

func (h *hashReader) hash() string {
	hash := sha1.Sum(h.buf.Bytes())
	return hex.EncodeToString(hash[:])
}

func main() {
	payload := []byte("Hello, World!")
	hashAndBroadcast(NewHashReader(payload))
}
*/

//2. Aggregate Data

func main() {

}

func fetch() string {
	time.Sleep(100 * time.Millisecond)
	return "data"
}
