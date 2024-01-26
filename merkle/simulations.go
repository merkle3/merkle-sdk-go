package merkle

import (
	"context"
	"fmt"
)

type SimulationAPI struct {
	sdk *MerkleSDK
}

func NewSimulationAPI(sdk *MerkleSDK) *SimulationAPI {
	return &SimulationAPI{
		sdk: sdk,
	}
}

type BundleCall struct {
	From      string                   `json:"from,omitempty"`
	To        string                   `json:"to,omitempty"`
	Value     string                   `json:"value,omitempty"`
	Nonce     int64                    `json:"nonce,omitempty"`
	Data      string                   `json:"data,omitempty"`
	GasLimit  int64                    `json:"gasLimit,omitempty"`
	Overrides *StateOverrideParameters `json:"overrides,omitempty"`
}

type SimulationBundle struct {
	ChainId     MerkleChainId            `json:"chainId"`
	Calls       []BundleCall             `json:"calls"`
	BlockNumber *int                     `json:"blockNumber,omitempty"`
	Overrides   *StateOverrideParameters `json:"overrides,omitempty"`
}

type StateOverrideParameters struct {
	Accounts      map[string]*AccountParameters `json:"accounts,omitempty"`
	ContractCodes map[string]string             `json:"contractCodes,omitempty"`
	Storage       map[string]map[string]string  `json:"storage,omitempty"`
}

type AccountParameters struct {
	Nonce   *int `json:"nonce,omitempty"`
	Balance *int `json:"balance,omitempty"`
}

type SimulationCallResult struct {
	Logs              []Log              `json:"logs"`
	GasUsed           *BigInt            `json:"gasUsed"`
	Result            string             `json:"result"`
	AddressCreated    *string            `json:"addressCreated,omitempty"`
	Status            int                `json:"status"`
	Error             *ErrorDetails      `json:"error,omitempty"`
	InternalTransfers []InternalTransfer `json:"internalTransfers,omitempty"`
}

type Log struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

type ErrorDetails struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type InternalTransfer struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount *BigInt `json:"amount"`
}

type SimulationResult struct {
	ChainId     int                    `json:"chainId"`
	BlockNumber *BigInt                `json:"blockNumber"`
	ProcessTime int                    `json:"processTime"`
	Calls       []SimulationCallResult `json:"calls"`
}

func (s *SimulationAPI) SimulateBundle(ctx context.Context, bundle *SimulationBundle) (*SimulationResult, error) {
	var result SimulationResult

	err := MakePost(
		ctx,
		"https://mbs-api.merkle.io/v1/simulate",
		s.sdk.GetApiKey(),
		bundle,
		&result,
	)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return &result, nil
}
