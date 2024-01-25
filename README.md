<img src="public/merkle-large.png" width="80" height="80" style="border-radius: 4px"/>

**merkle is building crypto infrastructure**. [Join us on discord](https://discord.gg/Q9Dc7jVX6c).

# merkle SDK

The merkle SDK is a great way to access our products.

## Install

Install the merkle SDK package:

```
go get github.com/merkle3/merkle-sdk-go
```

## Authentication

Get an API key from [merkle Blockchain Services (MBS)](https://mbs.merkle.io).

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at mbs.https://mbs.merkle.io
}
```

# Features

## Transaction Network

### Stream transactions

Access merkle's private stream of transactions on Ethereum & Polygon. [Learn more](https://docs.merkle.io/transaction-network/what-is-transaction-network)

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at https://mbs.merkle.io

    txs, err := merkleSdk.Transactions().Stream(merkle.EthereumMainnet) // pass a chain id, e.g. merkle.EthereumMainnet, merkle.PolygonMainnet or merkle.BnbMainnet

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
```

### Transaction tracing

Know exactly when and where a transaction was broadcasted. [Learn more](https://docs.merkle.io/transaction-network/tracing)

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at https://mbs.merkle.io

    trace, err := merkleSdk.Transactions().Trace("0x....") // a transaction hash

    // check for error
    if err != nil {
        fmt.Printf("error: %v\n", err)
        return
    }

    fmt.Printf("first seen at: %v\n", trace.FirstSeenAt.String())
}
```

## Private Mempool

### Stream auctions

Stream auctions from the Merkle Private Pool. [Learn more](https://docs.merkle.io/private-pool/what-is-private-mempool).

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at https://mbs.merkle.io

    auctions, err := merkleSdk.Pool().Auctions()

    for {
        select {
            case e := <-err:
            // error happened
            case auction := <-auctions:
            // process the auction, create a backrun

            // then send the bid
            err := auction.SendBid(tx) // a signed types.Transaction

            // or send a raw bid
            err := auction.SendRawBid([]string{
                // hex encoded bid
                "0x...."
            })

            // check for error in case the auction is already closed
        }
    }
}
```

### Send transaction to the private mempool

Send Ethereum, BSC and Polygon transactions to the private mempool to get MEV protection and recovery. [Learn more](https://docs.merkle.io/private-pool/what-is-private-mempool)

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at https://mbs.merkle.io

    err := merkleSdk.Pool().Send(&merkle.NewTransactionOptions{
        tx: nil, // a types.Transaction from go-ethereum
    })

    if err != nil {
        fmt.Printf("error: %v\n", err)
    }
}
```

## Simulations

### Simulate a bundle of transactions

The simulation API allows you to simulate a bundle of transactions on Ethereum, BSC and Polygon. [Learn more](https://docs.merkle.io/simulations/what-are-simulations)

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
   merkleSdk := merkle.New()

	merkleSdk.SetApiKey("sk_mbs_........") // get one at https://mbs.merkle.io

	block := 19078685

	simulationResult, err := merkleSdk.Simulation().SimulateBundle(&merkle.SimulationBundle{
		ChainId:     1,      // Ethereum Mainnet
		BlockNumber: &block, // nil for latest block, or a block number
		Calls: []merkle.BundleCall{
			{
				From: "0x3b42a0ed9050A79d8F35B07021272B3ef073266A",
				To:   "0x881D40237659C251811CEC9c364ef91dC08D300C",
				Data: "0x5f5755290000000000000000000000000000000000000000000000000000000000000080000000000000000000000000b2617246d0c6c0087f18703d576831899ca94f0100000000000000000000000000000000000000000000152d02c7e14af680000000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000136f6e65496e6368563546656544796e616d6963000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000b2617246d0c6c0087f18703d576831899ca94f01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000152d02c7e14af68000000000000000000000000000000000000000000000000000001464a6568a138eed0000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000002f063f978aab00000000000000000000000000f326e4de8f66a0bdc0970b79e0924e33c79f1915000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c80502b1c5000000000000000000000000b2617246d0c6c0087f18703d576831899ca94f0100000000000000000000000000000000000000000000152d02c7e14af68000000000000000000000000000000000000000000000000000001492bbd24caad01b0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000140000000000000003b6d0340b36ec83d844c0579ec2493f10b2087e96bb65460ab4991fe00000000000000000000000000000000000000000000000000a0",
			},
		},
	})

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Printf("simulation result: %+v\n", simulationResult)
}
```
