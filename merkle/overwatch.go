package merkle

type OverwatchAPI struct {
	sdk *MerkleSDK
}

func NewOverwatchAPI(sdk *MerkleSDK) *OverwatchAPI {
	return &OverwatchAPI{
		sdk: sdk,
	}
}

func (o *OverwatchAPI) WatchAddress(address string) error {
	type AddAddressRequest struct {
		Address string `json:"address"`
	}

	req := AddAddressRequest{
		Address: address,
	}

	err := MakePost("https://mbs-api.merkle.io/v1/overwatch/add_address", o.sdk.ApiKey, req, nil)

	if err != nil {
		return err
	}

	return nil
}

func (o *OverwatchAPI) UnwatchAddress(address string) error {
	err := MakeDel("https://mbs-api.merkle.io/v1/overwatch/addresses/"+address, o.sdk.ApiKey, nil, nil)

	if err != nil {
		return err
	}

	return nil
}

// declare hash
func (o *OverwatchAPI) Declare(hash string, chainId int64) error {
	err := MakePost("https://mbs-api.merkle.io/v1/overwatch/declare", o.sdk.ApiKey, map[string]interface{}{
		"hash":    hash,
		"chainId": chainId,
	}, nil)

	if err != nil {
		return err
	}

	return nil
}
