package blockchain

// Block 区块定义
type Block struct {
	Index        int
	Proof        int
	TimeStamp    int64
	PreviousHash string
	Transactions []*Transaction
}
