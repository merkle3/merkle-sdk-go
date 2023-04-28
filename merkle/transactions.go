package merkle

import "github.com/ethereum/go-ethereum/core/types"

type TransactionStream struct {
	sdk *MerkleSDK
}

func NewTransactionStream(sdk *MerkleSDK) *TransactionStream {
	return &TransactionStream{
		sdk: sdk,
	}
}

func Stream() (chan *types.Transaction, chan error) {
	return nil, nil
}
