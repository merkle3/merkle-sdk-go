package main

import (
	"fmt"

	"github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
	merkleSdk := merkle.New()

	auctions, err := merkleSdk.Pool().Auctions() // pass a chain id

	for {
		select {
		case e := <-err:
			// error happened
			fmt.Printf("error: %v\n", e)
		case auction := <-auctions:
			// process the transaction
			fmt.Printf("hash: %v\n", auction.Id)
		}
	}
}
