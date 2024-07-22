package main

import (
	"fmt"
	"os"

	"github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
	merkleSdk := merkle.New()

	merkleSdk.SetApiKey(os.Getenv("MERKLE_API_KEY"))

	auctions, err := merkleSdk.Pool().Auctions()

	for {
		select {
		case e := <-err:
			// error happened
			fmt.Printf("error: %v\n", e)
		case auction := <-auctions:
			// process the transaction
			fmt.Printf("auction tx: %+v\n", auction.Transaction)
		}
	}
}
