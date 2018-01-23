package blockchain

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/satori/go.uuid"
)

var (
	nodeidentifier = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	blockchain     = NewBlockChain()
	port           string
)

func init() {
	flag.StringVar(&port, "port", ":8080", "端口")
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	flag.Parse()
}

// Response .
type Response struct {
	Index        int            `json:"index,omitempty"`
	Length       int            `json:"length,omitempty"`
	Proof        int            `json:"proof,omitempty"`
	Message      string         `json:"message,omitempty"`
	PreviousHash string         `json:"previous_hash,omitempty"`
	Chain        []*Block       `json:"chain,omitempty"`
	Transactions []*Transaction `json:"transactions,omitempty"`
}

func mine(w http.ResponseWriter, r *http.Request) {
	lastBlock := blockchain.LastBlock()
	lastProof := lastBlock.Proof
	proof := blockchain.ProofOfWord(lastProof)

	blockchain.NewTransaction("0", nodeidentifier, 1)

	block := blockchain.NewBlock(proof, "")
	response := &Response{
		Message:      "New Block Forged",
		Index:        block.Index,
		Proof:        block.Proof,
		PreviousHash: block.PreviousHash,
		Transactions: block.Transactions,
	}
	doc, _ := json.Marshal(&response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(doc)
}

func newTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write([]byte("Missing values"))
		return
	}

	var values Transaction

	doc, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(doc, &values); err != nil {
		w.Write([]byte("Missing values"))
		return
	}
	index := blockchain.NewTransaction(values.Sender, values.Recipient, values.Amount)

	response, _ := json.Marshal(
		&Response{
			Message: fmt.Sprintf("Transaction will be added to Block %d", index),
		})
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(response)
}

func fullChain(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Chain:  blockchain.Chain,
		Length: blockchain.Len(),
	}
	doc, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(doc)
}

// Run .
func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/mine", mine)
	mux.HandleFunc("/transactions/new", newTransaction)
	mux.HandleFunc("/chain", fullChain)

	fmt.Printf("Listen: http://127.0.0.1%s\n", port)
	http.ListenAndServe(port, mux)
}
