package merkle

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/merkle3/merkle-sdk-go/proto"
	brokerProto "github.com/merkle3/merkle-sdk-go/proto"
	"google.golang.org/grpc"
)

type TransactionStream struct {
	sdk *MerkleSDK
}

func NewTransactionStream(sdk *MerkleSDK) *TransactionStream {
	return &TransactionStream{
		sdk: sdk,
	}
}

func (t *TransactionStream) Stream(chainId int32) (chan *types.Transaction, chan error) {
	errStream := make(chan error)
	txs := make(chan *types.Transaction)

	conn, err := grpc.Dial("txs.usemerkle.com:80", grpc.WithInsecure())

	if err != nil {
		go func() {
			errStream <- err
		}()
		return nil, errStream
	}

	broker := brokerProto.NewBrokerApiClient(conn)

	stream, err := broker.StreamReceivedTransactions(context.Background(), &brokerProto.TxStreamRequest{
		ApiKey:  t.sdk.GetApiKey(),
		ChainId: chainId,
	})

	go func() {
		for {
			txPacket, err := stream.Recv()

			if err != nil {
				errStream <- err
				return
			}

			tx, err := proto.ToTransaction(txPacket)

			if err != nil {
				errStream <- err
				return
			}

			txs <- tx
		}
	}()

	return txs, errStream
}
