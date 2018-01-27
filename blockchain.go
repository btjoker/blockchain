package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// BlockChain 区块链定义
type BlockChain struct {
	Chain               []*Block       `json:"chain"`
	CurrentTransactions []*Transaction `json:"current_transactions"`
	Nodes               *Set           `json:"nodes"`
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
		previousHash = Hash(bc.LastBlock())
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

// RegisterNode 注册节点
func (bc *BlockChain) RegisterNode(address string) {
	u, err := url.Parse(address)
	if err != nil {
		return
	}
	bc.Nodes.Add(u.Host)
}

// ValidChain 验证区块链
func (bc *BlockChain) ValidChain(chain []*Block) bool {
	lastBlock := chain[0]
	currentIndex := 1
	chainLength := len(chain)

	for currentIndex < chainLength {
		block := chain[currentIndex]
		fmt.Println(lastBlock)
		fmt.Println(block)
		fmt.Println("\n----------\n")

		if block.PreviousHash != Hash(lastBlock) {
			return false
		}

		if !ValidProof(lastBlock.Proof, block.Proof) {
			return false
		}

		lastBlock = block
		currentIndex++
	}
	return true
}

// ResolveConflicts 解决区块链冲突
func (bc *BlockChain) ResolveConflicts() bool {
	var newChain []*Block
	neighbours := bc.Nodes
	maxLength := bc.Len()

	for node := range neighbours.items {
		resp, err := http.Get(fmt.Sprintf("http://%s/chain", node))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var newBlockChain BlockChain
		doc, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(doc, &newBlockChain); err != nil {
			continue
		}

		length := newBlockChain.Len()
		chain := newBlockChain.Chain

		if length > maxLength && bc.ValidChain(chain) {
			maxLength = length
			newChain = chain
		}
	}

	if len(newChain) != 0 {
		bc.Chain = newChain
		return true
	}
	return false
}

// Hash 返回区块的 Hash
func Hash(b *Block) string {
	headers, _ := json.Marshal(b)

	value := sha256.Sum256(headers)
	return fmt.Sprintf("%x", value)
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
		Nodes: NewSet(),
	}
}
