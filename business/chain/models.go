package chain

// Block - represents a block in the blockchain.
type Block struct {
	BlockNumber       int64         `json:"blockNumber"`
	Timestamp         uint64        `json:"timestamp"`
	Difficulty        uint64        `json:"difficulty"`
	Hash              string        `json:"hash"`
	TransactionsCount int           `json:"transactionsCount"`
	Transactions      []Transaction `json:"transactions"`
}

// Transaction - represents single transaction in the blockchain.
type Transaction struct {
	Hash     string `json:"hash"`
	Value    string `json:"value"`
	Gas      uint64 `json:"gas"`
	GasPrice uint64 `json:"gasPrice"`
	Nonce    uint64 `json:"nonce"`
	To       string `json:"to"`
	Pending  bool   `json:"pending"`
}

// TransferEthRequest - represents request details for transferring ETH.
type TransferEthRequest struct {
	PrivateKey string `json:"privateKey"`
	To         string `json:"to"`
	Amount     int64  `json:"amount"`
}

// HashResponse - represents response for hash.
type HashResponse struct {
	Hash string `json:"hash"`
}

// BalanceResponse - represents response for balance.
type BalanceResponse struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
	Symbol  string `json:"symbol"`
	Units   string `json:"units"`
}

// Error - represents error response.
type Error struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}
