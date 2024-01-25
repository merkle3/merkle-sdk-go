package merkle

type OverwatchAPI struct {
	sdk *MerkleSDK
}

func NewOverwatchAPI(sdk *MerkleSDK) *OverwatchAPI {
	return &OverwatchAPI{
		sdk: sdk,
	}
}

func (o *OverwatchAPI) WatchAddress(address string) {
	// TODO
}

func (o *OverwatchAPI) UnwatchAddress(address string) {
	// TODO
}

// declare hash
func (o *OverwatchAPI) Declare(hash string, chainId int64) {
	// TODO
}
