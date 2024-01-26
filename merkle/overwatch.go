package merkle

import "context"

type OverwatchAPI struct {
	sdk *MerkleSDK
}

func NewOverwatchAPI(sdk *MerkleSDK) *OverwatchAPI {
	return &OverwatchAPI{
		sdk: sdk,
	}
}

func (o *OverwatchAPI) WatchAddress(ctx context.Context, address string) error {
	type AddAddressRequest struct {
		Address string `json:"address"`
	}

	req := AddAddressRequest{
		Address: address,
	}

	err := MakePost(ctx, "https://mbs-api.merkle.io/v1/overwatch/addresses", o.sdk.ApiKey, req, nil)

	if err != nil {
		return err
	}

	return nil
}

func (o *OverwatchAPI) UnwatchAddress(ctx context.Context, address string) error {
	err := MakeDel(ctx, "https://mbs-api.merkle.io/v1/overwatch/addresses/"+address, o.sdk.ApiKey, nil, nil)

	if err != nil {
		return err
	}

	return nil
}

// declare hash
func (o *OverwatchAPI) Declare(ctx context.Context, chainId MerkleChainId, hash string) error {
	err := MakePost(ctx, "https://mbs-api.merkle.io/v1/overwatch/declare", o.sdk.ApiKey, map[string]interface{}{
		"hash":    hash,
		"chainId": chainId,
	}, nil)

	if err != nil {
		return err
	}

	return nil
}
