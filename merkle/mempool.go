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
	ChainId      int64
	CreatedAt    time.Time
	Connection   *websocket.Conn

	Transaction *AuctionTransaction
}

type RawRpcResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
}

type RawAuction struct {
	Id           string `json:"id"`
	FeeRecipient string `json:"fee_recipient"`
	ClosesAtUnix int64  `json:"closes_at_unix"`
	ChainId      int64  `json:"chain_id"`
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

type NewTransactionOptions struct {
	Transaction  *types.Transaction
	FeeRecipient common.Address

	// optionally, a source
	Source string

	// prevent reverts
	PreventRevert bool

	// hints
	Hints []string

	// privacy profile
	PrivacyProfile string
}

func (p *PrivatePool) Send(options *NewTransactionOptions) error {
	// send to the pool
	type PoolSubmission struct {
		// An array of transactions
		Transactions []string `json:"transactions"`

		// The fee recipient
		FeeRecipient string `json:"fee_recipient"`

		// Optional, a source tag
		Source string `json:"source"`

		// Optional, a privacy profile
		Privacy string `json:"privacy"`

		// Optional, a list of hints, overrides the privacy profile
		Hints []string `json:"hints"`

		// Optional, a list of allowed bundles for this transaction
		BundleTypes []string `json:"bundle_types"`

		// Optional, a list of release targets
		ReleaseTargets []string `json:"release_targets"`

		// Optional, prevent reverts
		PreventReverts bool `json:"prevent_reverts"`
	}

	signer := types.LatestSignerForChainID(options.Transaction.ChainId())
	txFrom, err := signer.Sender(options.Transaction)

	if err != nil {
		return fmt.Errorf("failed to get transaction sender: %s", err)
	}

	feeRecipient := txFrom.String()

	if options.FeeRecipient.String() != (common.Address{}).String() {
		feeRecipient = options.FeeRecipient.String()
	}

	txBytes, err := options.Transaction.MarshalBinary()

	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %s", err)
	}

	submission := &PoolSubmission{
		Transactions:   []string{common.Bytes2Hex(txBytes)},
		FeeRecipient:   feeRecipient,
		Source:         options.Source,
		Privacy:        options.PrivacyProfile,
		Hints:          options.Hints,
		PreventReverts: options.PreventRevert,
	}

	submissionBody, err := json.Marshal(submission)

	if err != nil {
		return fmt.Errorf("failed to marshal submission: %s", err)
	}

	client := &http.Client{}

	// send to the pool
	req, err := http.NewRequest("POST", "https://mempool.merkle.io/transactions", bytes.NewBuffer(submissionBody))

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
		conn, err := websocket.Dial("wss://mempool.merkle.io/stream/auctions?apiKey="+p.sdk.GetApiKey(), "", "http://localhost/")

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
					ChainId:      rawAuction.ChainId,
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

func (a *Auction) SendBid(tx types.Transaction) (string, error) {
	bin, err := tx.MarshalBinary()

	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction: %s", err)
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

func (a *Auction) SendRawBid(txs []string) (string, error) {
	// send a request to https://pool.merkle.io/relay
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
		return "", fmt.Errorf("failed to marshal payload: %s", err)
	}

	res, err := http.Post("https://pool.merkle.io/relay", "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", fmt.Errorf("failed to send request: %s", err)
	}

	if res.StatusCode > 400 {
		return "", fmt.Errorf("failed to send request: code=%s", res.Status)
	}

	defer res.Body.Close()

	// decode the response
	var resBody map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&resBody)

	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %s", err)
	}

	bidId, ok := resBody["result"].(string)

	if !ok {
		return "", fmt.Errorf("failed to get bid id")
	}

	return bidId, nil
}
