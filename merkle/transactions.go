package merkle

import (
	"fmt"
	"time"

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
	txStream := make(chan *types.Transaction)

	incomingMessages := make(chan []uint8)

	if t.sdk.ApiKey == "" {
		go func() {
			errStream <- fmt.Errorf("API key is not set")
		}()
		return txStream, errStream
	}

	go func() {
		for {
			var address = "txs.merkle.io"
			ws, err := websocket.Dial(fmt.Sprintf("wss://%s/ws/%s", address, t.sdk.ApiKey), "", fmt.Sprintf("http://%s/", address))

			if err != nil {
				go func() {
					errStream <- err
				}()
				return
			}

			for {
				var message []uint8

				// set a deadline of 5 seconds to receive the next tx,
				// should be plenty of time
				ws.SetReadDeadline(time.Now().Add(5 * time.Second))
				err := websocket.Message.Receive(ws, &message)
				if err != nil {
					// if we couldn't read the message, try to reconnect
					break
				}
				incomingMessages <- message
			}

			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case message := <-incomingMessages:
				tx := types.Transaction{}

				err := tx.UnmarshalBinary(message)

				if err != nil {
					// if we couldn't parse the transaction, skip it
					continue
				}

				txStream <- &tx
			}
		}
	}()

	return txStream, errStream
}
