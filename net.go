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
	nodeidentifier = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	blockchain     = NewBlockChain()
	port           string
)

// nodeTmp 临时存放节点
type nodeTmp struct {
	Nodes []string `json:"nodes"`
}

func init() {
	flag.StringVar(&port, "port", ":8080", "端口")
	flag.Parse()
}

// Response josn 响应
type Response struct {
	Index        int            `json:"index,omitempty"`
	Length       int            `json:"length,omitempty"`
	Proof        int            `json:"proof,omitempty"`
	Message      string         `json:"message,omitempty"`
	PreviousHash string         `json:"previous_hash,omitempty"`
	Chain        []*Block       `json:"chain,omitempty"`
	NewChain     []*Block       `json:"new_chain,omitempty"`
	Transactions []*Transaction `json:"transactions,omitempty"`
	TotalNodes   []string       `json:"total_nodes,omitempty"`
}

// mine 挖矿处理
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

// newTransaction 新交易处理
func newTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
}

// fullChain 返回全部区块链
func fullChain(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Chain:  blockchain.Chain,
		Length: blockchain.Len(),
	}
	doc, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(doc)
}

// registerNodes 注册节点
func registerNodes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var values nodeTmp
		doc, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(doc, &values); err != nil {
			return
		}

		if len(values.Nodes) == 0 {
			w.WriteHeader(400)
			w.Write([]byte("Error: Please supply a valid list of nodes"))
			return
		}

		for _, node := range values.Nodes {
			blockchain.RegisterNode(node)
		}

		data, _ := json.Marshal(Response{
			Message:    "New nodes have been added",
			TotalNodes: blockchain.Nodes.List(),
		})
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(data)
	}
}

// consensus 处理同步
func consensus(w http.ResponseWriter, r *http.Request) {
	var data []byte
	if blockchain.ResolveConflicts() {
		data, _ = json.Marshal(Response{
			Message:  "Our chain was replaced",
			NewChain: blockchain.Chain,
		})
	} else {
		data, _ = json.Marshal(Response{
			Message: "Our chain is authoritative",
			Chain:   blockchain.Chain,
		})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

// Run .
func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/mine", mine)
	mux.HandleFunc("/transactions/new", newTransaction)
	mux.HandleFunc("/chain", fullChain)
	mux.HandleFunc("/nodes/resolve", consensus)
	mux.HandleFunc("/nodes/register", registerNodes)

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	fmt.Printf("Listen: http://127.0.0.1%s\n", port)
	http.ListenAndServe(port, mux)
}
