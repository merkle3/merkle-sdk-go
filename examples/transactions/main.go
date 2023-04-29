package main

import (
	"fmt"

	"github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
	merkleSdk := merkle.New()

	txs, err := merkleSdk.Transactions().Stream(1) // pass a chain id

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
