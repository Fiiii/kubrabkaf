package chain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

// ETHClient - Ethereum client wrapper.
type ETHClient struct {
	Client *ethclient.Client
}

// GetLatestBlock - get latest block from the chain. It returns a Block struct with all the transactions.
func (e *ETHClient) GetLatestBlock() (*Block, error) {
	// Query last block
	header, _ := e.Client.HeaderByNumber(context.Background(), nil)

	blockNum := big.NewInt(header.Number.Int64())
	b, err := e.Client.BlockByNumber(context.Background(), blockNum)
	if err != nil {
		return nil, err
	}

	// Create block
	block := &Block{
		BlockNumber:       b.Number().Int64(),
		Timestamp:         b.Time(),
		Difficulty:        b.Difficulty().Uint64(),
		Hash:              b.Hash().String(),
		TransactionsCount: len(b.Transactions()),
		Transactions:      make([]Transaction, 0),
	}

	// Fill transactions
	for _, tx := range b.Transactions() {
		block.Transactions = append(block.Transactions, Transaction{
			Hash:     tx.Hash().String(),
			Value:    tx.Value().String(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice().Uint64(),
			Nonce:    tx.Nonce(),
			To:       tx.To().String(),
		})
	}

	return block, nil
}

// GetTxByHash - get transaction by hash. Hash = 32 byte Keccak256 hash.
func (e *ETHClient) GetTxByHash(hash common.Hash) (*Transaction, error) {
	tx, isPending, err := e.Client.TransactionByHash(context.Background(), hash)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Hash:     tx.Hash().String(),
		Value:    tx.Value().String(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice().Uint64(),
		Nonce:    tx.Nonce(),
		Pending:  isPending,
		To:       tx.To().String(),
	}, nil
}

// TransferEth - transfer ETH from one address to another.
func (e *ETHClient) TransferEth(privateKey string, toAddress string, amount int64) (*HashResponse, error) {
	var (
		sk       = crypto.ToECDSAUnsafe(common.FromHex(privateKey))
		to       = common.HexToAddress(toAddress)
		value    = new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether))
		sender   = common.HexToAddress(privateKey)
		gasLimit = uint64(21000)
	)
	chainid, err := e.Client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chainID: %w", err)
	}

	nonce, err := e.Client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}

	tipCap, _ := e.Client.SuggestGasTipCap(context.Background())
	feeCap, _ := e.Client.SuggestGasPrice(context.Background())

	// Create a new transaction
	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:   chainid,
			Nonce:     nonce,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       gasLimit,
			To:        &to,
			Value:     value,
			Data:      nil,
		})
	// Sign the transaction using our keys
	signedTx, _ := types.SignTx(tx, types.NewLondonSigner(chainid), sk)

	// Send the transaction to our node
	err = e.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("send tx: %w", err)
	}

	return &HashResponse{
		Hash: signedTx.Hash().String(),
	}, err
}

// GetAddressBalance - get address balance.
func (e *ETHClient) GetAddressBalance(address string) (*BalanceResponse, error) {
	balance, err := e.Client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return nil, err
	}

	return &BalanceResponse{
		Balance: balance.String(),
	}, nil
}
