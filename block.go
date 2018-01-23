package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// Block 区块定义
type Block struct {
	Index        int
	Proof        int
	TimeStamp    int64
	PreviousHash string
	Transactions []*Transaction
}

// Hash 返回区块的 Hash
func (b *Block) Hash() string {
	headers, _ := json.Marshal(b)

	hash := sha256.Sum256(headers)
	return fmt.Sprintf("%x", hash)
}
