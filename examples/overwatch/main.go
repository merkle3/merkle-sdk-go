package main

import (
	"os"

	"github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
	merkleSdk := merkle.New()

	merkleSdk.SetApiKey(os.Getenv("MERKLE_API_KEY")) // get one at mbs.https://mbs.merkle.io

	err := merkleSdk.Overwatch().WatchAddress("0x3b42a0ed9050A79d8F35B07021272B3ef073266A")

	if err != nil {
		panic(err)
	}

	// declare hashes on Ethereum mainnet
	err = merkleSdk.Overwatch().Declare(merkle.EthereumMainnet, "0x....")

	if err != nil {
		panic(err)
	}

	// unwatch address on Ethereum mainnet
	err = merkleSdk.Overwatch().UnwatchAddress("0x3b42a0ed9050A79d8F35B07021272B3ef073266A")

	if err != nil {
		panic(err)
	}
}
