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
	Connection   *websocket.Conn

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

type AuctionOptions struct {
	Transaction  *types.Transaction
	FeeRecipient common.Address
}

func (p *PrivatePool) CreateAuction(options *AuctionOptions) error {
	// send to the pool
	type PoolSubmission struct {
		// an array of transactions
		Transactions []string `json:"transactions"`

		// the fee recipient
		FeeRecipient string `json:"fee_recipient"`

		// optionally, a source
		Source string `json:"source"`
	}

	signer := types.LatestSignerForChainID(options.Transaction.ChainId())
	txFrom, err := signer.Sender(options.Transaction)

	if err != nil {
		return fmt.Errorf("failed to get transaction sender: %s", err)
	}

	feeReceipient := txFrom.String()

	if options.FeeRecipient.String() != (common.Address{}).String() {
		feeReceipient = options.FeeRecipient.String()
	}

	txBytes, err := options.Transaction.MarshalBinary()

	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %s", err)
	}

	submission := &PoolSubmission{
		Transactions: []string{common.Bytes2Hex(txBytes)},
		FeeRecipient: feeReceipient,
		Source:       "public",
	}

	submissionBody, err := json.Marshal(submission)

	if err != nil {
		return fmt.Errorf("failed to marshal submission: %s", err)
	}

	client := &http.Client{}

	// send to the pool
	req, err := http.NewRequest("POST", "https://pool.usemerkle.com/transactions", bytes.NewBuffer(submissionBody))

	if err != nil {
		return fmt.Errorf("failed to create request to pool: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if p.sdk.GetApiKey() != "" {
		req.Header.Set("X-MBS-Key", p.sdk.GetApiKey())
	}

	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("failed to send request to pool: %s", err)
	}

	if res.StatusCode > 400 {
		return fmt.Errorf("failed to send request to pool: code=%s", res.Status)
	}

	return nil
}

func (p *PrivatePool) Auctions() (chan *Auction, chan error) {
	auctions := make(chan *Auction)
	errStream := make(chan error)

	connect := func(auctionChannel chan *Auction, errChannel chan error) {
		conn, err := websocket.Dial("wss://pool.usemerkle.com/stream/auctions?apiKey="+p.sdk.GetApiKey(), "", "http://localhost/")

		if err != nil {
			go func() {
				errStream <- err
			}()
			return
		}

		go func() {
			for {
				var rawAuction RawAuction
				var rawJSON string
				var rawJSONTotal = ""

				// sometimes, the auctions are too big and split
				// into multiple frames, we need to combine them
				for {
					err := websocket.Message.Receive(conn, &rawJSON)

					if err != nil {
						errStream <- fmt.Errorf("failed to receive message: %s", err)
						return
					}

					rawJSONTotal += rawJSON

					err = json.Unmarshal([]byte(rawJSONTotal), &rawAuction)

					if err == nil {
						break
					}
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
					// keep track of the connection for bids
					Connection: conn,
				}

				auctions <- &auction
			}
		}()
	}

	connect(auctions, errStream)

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

	_, err = a.Connection.Write(body)

	if err != nil {
		return fmt.Errorf("failed to send request: %s", err)
	}

	return nil
}
