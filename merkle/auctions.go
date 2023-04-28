package merkle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/net/websocket"
)

type PrivatePool struct {
	sdk *MerkleSDK
}

func NewPrivatePool(sdk *MerkleSDK) *PrivatePool {
	return &PrivatePool{
		sdk: sdk,
	}
}

type AuctionTransaction struct {
	Hash  common.Hash
	From  common.Address
	To    common.Address
	Value *big.Int
	Data  []byte
	Gas   uint64
}

type Auction struct {
	Id           string
	FeeRecipient string
	ClosesAt     time.Time
	CreatedAt    time.Time

	Transaction *AuctionTransaction
}

type RawAuction struct {
	Id           string `json:"id"`
	FeeRecipient string `json:"fee_recipient"`
	ClosesAtUnix int64  `json:"closes_at_unix"`
	CreatedAt    int64  `json:"created_at_unix"`
	Transaction  struct {
		Data  string `json:"data"`
		From  string `json:"from"`
		Gas   int64  `json:"gas"`
		Hash  string `json:"hash"`
		To    string `json:"to"`
		Value string `json:"value"`
	} `json:"transaction"`
}

func Auctions() (chan *Auction, chan error) {
	auctions := make(chan *Auction)
	errStream := make(chan error)

	conn, err := websocket.Dial("wss://pool.usemerkle.com/stream/auctions", "", "")

	if err != nil {
		go func() {
			errStream <- err
		}()
		return nil, errStream
	}

	go func() {
		for {
			var rawAuction RawAuction

			err := websocket.JSON.Receive(conn, &rawAuction)

			if err != nil {
				errStream <- err
				return
			}

			auction := Auction{
				Id:           rawAuction.Id,
				FeeRecipient: rawAuction.FeeRecipient,
				ClosesAt:     time.Unix(rawAuction.ClosesAtUnix, 0),
				CreatedAt:    time.Unix(rawAuction.CreatedAt, 0),
				Transaction: &AuctionTransaction{
					Hash:  common.HexToHash(rawAuction.Transaction.Hash),
					From:  common.HexToAddress(rawAuction.Transaction.From),
					To:    common.HexToAddress(rawAuction.Transaction.To),
					Value: new(big.Int),
					Data:  []byte(rawAuction.Transaction.Data),
					Gas:   uint64(rawAuction.Transaction.Gas),
				},
			}

			auctions <- &auction
		}
	}()

	return auctions, errStream
}

func (a *Auction) SendBid(tx types.Transaction) error {
	bin, err := tx.MarshalBinary()

	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %s", err)
	}

	hex := common.Bytes2Hex(bin)

	return a.SendRawBid([]string{hex})
}

type RelaySubmitRequest struct {
	Method  string         `json:"method"`
	Params  []BundleParams `json:"params"`
	Jsonrpc string         `json:"jsonrpc"`
}

type BundleParams struct {
	Txs         []string `json:"txs"`
	BlockNumber string   `json:"blockNumber"`
}

func (a *Auction) SendRawBid(txs []string) error {
	// send a request to https://pool.usemerkle.com/relay
	payload := &RelaySubmitRequest{
		Method: "eth_sendBundle",
		Params: []BundleParams{
			{
				Txs:         txs,
				BlockNumber: "0",
			},
		},
		Jsonrpc: "2.0",
	}

	body, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal payload: %s", err)
	}

	res, err := http.Post("https://pool.usemerkle.com/relay", "application/json", bytes.NewBuffer(body))

	if err != nil {
		return fmt.Errorf("failed to send request: %s", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("failed to send request: %s", res.Status)
	}

	return nil
}
