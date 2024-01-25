package merkle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
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
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
	Data  string `json:"data"`
}

type Bundle struct {
	ChainId     int64        `json:"chain_id"`
	Calls       []BundleCall `json:"calls"`
	BlockNumber int64        `json:"block_number"`
}

type SimulationCallResult struct {
	Logs              []Log              `json:"logs"`
	GasUsed           *big.Int           `json:"gasUsed"`
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
	From   string   `json:"from"`
	To     string   `json:"to"`
	Amount *big.Int `json:"amount"`
}

type SimulationResult struct {
	ChainId     int                    `json:"chainId"`
	BlockNumber *big.Int               `json:"blockNumber"`
	ProcessTime int                    `json:"processTime"`
	Calls       []SimulationCallResult `json:"calls"`
}

func (s *SimulationAPI) SimulateBundle(bundle Bundle) (*SimulationResult, error) {
	var result SimulationResult

	// call the mbs-api.merkle.io/v1/simulate endpoint
	bodyBytes, err := json.Marshal(bundle)

	if err != nil {
		return nil, fmt.Errorf("error marshalling bundle: %v", err)
	}

	req, err := http.NewRequest("POST", "https://mbs-api.merkle.io/v1/simulate", bytes.NewBuffer(bodyBytes))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+s.sdk.GetApiKey())

	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	if res.StatusCode > 400 {
		return nil, fmt.Errorf("error sending request: code=%s", res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
