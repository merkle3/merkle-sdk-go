package proto

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

func ToProto(tx *types.Transaction) *Transaction {
	txBytes, err := tx.MarshalBinary()

	if err != nil {
		log.WithError(err).Error("failed to marshal transaction")
		return nil
	}

	hex := common.Bytes2Hex(txBytes)

	return &Transaction{
		// transaction data
		TxHash:  tx.Hash().String(),
		TxBytes: hex,
	}
}

func ToTransaction(protoTx *Transaction) (*types.Transaction, error) {
	txBytes := common.Hex2Bytes(protoTx.TxBytes)

	tx := new(types.Transaction)

	err := tx.UnmarshalBinary(txBytes)

	if err != nil {
		log.WithError(err).Error("failed to unmarshal transaction")
		return nil, err
	}

	return tx, nil
}

// transforms an array into a protobuf array
func ToProtoArray(array []*types.Transaction) []*Transaction {
	protoArray := make([]*Transaction, 0)

	for _, tx := range array {
		protoArray = append(protoArray, ToProto(tx))
	}

	return protoArray
}
