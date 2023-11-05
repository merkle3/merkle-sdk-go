package merkle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
			ws, err := websocket.Dial(fmt.Sprintf("wss://%s/ws/%s/%d", address, t.sdk.ApiKey, int64(chainId)), "", fmt.Sprintf("http://%s/", address))

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

type MerkleTrace struct {
	Hash        string        `json:"hash"`
	FirstSeenAt time.Time     `json:"firstSeenAt"`
	ChainId     MerkleChainId `json:"chainId"`
	Trace       []Observation `json:"trace"`
	TxData      string        `json:"txData"`
}

type Observation struct {
	Time   time.Time
	Origin string
}

// trace a transaction
func (t *TransactionStream) Trace(hash string) (*MerkleTrace, error) {
	// url is https://txs.merkle.io/trace/<hash>
	res, err := http.Get(fmt.Sprintf("https://txs.merkle.io/trace/%s", hash))

	if err != nil {
		return nil, fmt.Errorf("error fetching trace: %s", err)
	}

	// check if we got a 404
	if res.StatusCode == 404 {
		return nil, nil
	}

	// check if we got a 200
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching trace: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading trace: %s", err)
	}

	// decode the response
	var trace MerkleTrace

	err = json.Unmarshal(body, &trace)

	if err != nil {
		return nil, fmt.Errorf("error decoding trace: %s", err)
	}

	return &trace, nil
}
