package merkle

import (
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
	From  string `json:"from,omitempty"`
	To    string `json:"to,omitempty"`
	Value string `json:"value,omitempty"`
	Nonce int64  `json:"nonce,omitempty"`
	Data  string `json:"data,omitempty"`
}

type SimulationBundle struct {
	ChainId     int64        `json:"chainId"`
	Calls       []BundleCall `json:"calls"`
	BlockNumber *int         `json:"blockNumber,omitempty"`
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

func (s *SimulationAPI) SimulateBundle(bundle *SimulationBundle) (*SimulationResult, error) {
	var result SimulationResult

	err := MakePost(
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