package main

import (
	"fmt"
	"os"

	"github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
	merkleSdk := merkle.New()

	merkleSdk.SetApiKey(os.Getenv("MERKLE_API_KEY"))

	txs, err := merkleSdk.Transactions().Stream(merkle.EthereumMainnet)

	for {
		select {
		case e := <-err:
			// error happened
			fmt.Printf("error: %v\n", e)
		case tx := <-txs:
			// process the transaction
			fmt.Printf("hash: %v\n", tx.Hash().String())
		}
	}
}
