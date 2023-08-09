package merkle

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/net/websocket"
)

type TransactionStream struct {
	sdk *MerkleSDK
}

func NewTransactionStream(sdk *MerkleSDK) *TransactionStream {
	return &TransactionStream{
		sdk: sdk,
	}
}

func (t *TransactionStream) Stream(chainId MerkleChainId) (chan *types.Transaction, chan error) {
	errStream := make(chan error)
	txs := make(chan *types.Transaction)

	var address = "txs.merkle.io"
	ws, err := websocket.Dial(fmt.Sprintf("wss://%s/ws/1/%s", address, t.sdk.ApiKey), "", fmt.Sprintf("http://%s/", address))

	if err != nil {
		go func() {
			errStream <- err
		}()
		return nil, errStream
	}
	incomingMessages := make(chan []uint8)

	// read incoming messages in a new goroutine
	go func(_ws *websocket.Conn, _mgs chan []uint8) {
		for {
			var message []uint8
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				fmt.Printf("Error::: %s\n", err.Error())
				return
			}
			_mgs <- message
		}
	}(ws, incomingMessages)

	go func() {
		for {
			select {
			case message := <-incomingMessages:
				tx := types.Transaction{}

				err := tx.UnmarshalBinary(message)

				if err != nil {
					errStream <- err
					return
				}

				txs <- &tx

			}
		}
	}()

	return txs, errStream
}
