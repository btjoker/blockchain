package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// BlockChain 区块链定义
type BlockChain struct {
	Chain               []*Block
	CurrentTransactions []*Transaction
}

// NewTransaction 新的交易
func (bc *BlockChain) NewTransaction(sender, recipient string, amount int) int {
	bc.CurrentTransactions = append(bc.CurrentTransactions, &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	})
	return bc.LastBlock().Index + 1
}

// NewBlock 新的区块
func (bc *BlockChain) NewBlock(proof int, previousHash string) *Block {
	if previousHash == "" {
		previousHash = bc.LastBlock().Hash()
	}

	block := &Block{
		Index:        bc.Len() + 1,
		TimeStamp:    time.Now().Unix(),
		Transactions: bc.CurrentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}

	bc.CurrentTransactions = []*Transaction{}
	bc.AddBlock(block)
	return block
}

// ProofOfWord 工作证明
func (bc *BlockChain) ProofOfWord(lastProof int) int {
	var proof int
	for !ValidProof(lastProof, proof) {
		proof++
	}
	return proof
}

// AddBlock 添加新的区块
func (bc *BlockChain) AddBlock(block *Block) {
	bc.Chain = append(bc.Chain, block)
}

// LastBlock 最后一个区块
func (bc *BlockChain) LastBlock() *Block {
	if bc.Len() != 0 {
		return bc.Chain[bc.Len()-1]
	}
	return nil
}

// Len 区块链长度
func (bc *BlockChain) Len() int {
	return len(bc.Chain)
}

// ValidProof 验证证明
func ValidProof(lastProof, proof int) bool {
	guess := bytes.Join([][]byte{
		[]byte(strconv.Itoa(lastProof)),
		[]byte(strconv.Itoa(proof)),
	}, []byte{})

	guessHash := sha256.Sum256(guess)
	return bytes.HasSuffix(guessHash[:], []byte("00"))
}

// NewBlockChain 实例化区块链
func NewBlockChain() *BlockChain {
	return &BlockChain{
		Chain: []*Block{
			&Block{
				PreviousHash: "1",
				Proof:        100,
				TimeStamp:    time.Now().Unix(),
			},
		},
	}
}
