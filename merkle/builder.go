package merkle

type BuilderSDK struct {
	sdk *MerkleSDK
}

func NewBuilderSDK(sdk *MerkleSDK) *BuilderSDK {
	return &BuilderSDK{
		sdk: sdk,
	}
}
