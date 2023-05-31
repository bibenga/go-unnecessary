package main

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"log"
	"math/big"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	i := int64(123)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	log.Printf("hex   : %v", hex.EncodeToString(b))
	log.Printf("base32: %v", base32.StdEncoding.EncodeToString(b))
	log.Printf("base64: %v", base64.StdEncoding.EncodeToString(b))

	log.Printf("big: %v - %v", big.MaxBase, 10+('Z'-'A'+1))
	bi := big.NewInt(9_123_456_789_012_456_789)
	log.Printf("10: %v", bi.Text(10))
	log.Printf("16: %v", bi.Text(16))
	log.Printf("32: %v", bi.Text(32))
	log.Printf("36: %v", bi.Text(36))
	log.Printf("62: %v", bi.Text(62))
}
