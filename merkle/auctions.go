package merkle

type PrivatePool struct {
	sdk *MerkleSDK
}

func NewPrivatePool(sdk *MerkleSDK) *PrivatePool {
	return &PrivatePool{
		sdk: sdk,
	}
}

type Auction struct {
}

func Auctions() (chan *Auction, chan error) {
	return nil, nil
}
