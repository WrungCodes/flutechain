package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.Hex()), nil
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

type BlockHeader struct {
	Parent Hash
	Time   uint64
}

type Block struct {
	Header BlockHeader
	Tx     []Tx
}

type BlockFS struct {
	Key   Hash
	Value Block
}

func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}

func NewBlock(hash Hash, time uint64, tx []Tx) Block {
	return Block{BlockHeader{hash, time}, tx}
}
