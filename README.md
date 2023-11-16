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

## Features

### Stream transactions

Access merkle's private stream of transactions on Ethereum & Polygon. [Learn more](https://docs.merkle.io/transaction-network/what-is-transaction-network)

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at mbs.https://mbs.merkle.io

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

### Stream auctions

Stream auctions from the Merkle Private Pool. [Learn more](https://docs.merkle.io/private-pool/what-is-private-mempool).

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at mbs.merkle.io

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

Send Ethereum, BSC and Polygon transactions to the private mempool to get MEV protection and recovery. [Learn more](https://docs.merkle.io/private-mempool/what-is-private-mempool)

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at mbs.https://mbs.merkle.io

    err := merkleSdk.Pool().Send(&merkle.NewTransactionOptions{
        tx: nil, // a types.Transaction from go-ethereum
    })

    if err != nil {
        fmt.Printf("error: %v\n", err)
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

    merkleSdk.SetApiKey("sk_mbs_......") // get one at mbs.merkle.io

    trace, err := merkleSdk.Transactions().Trace("0x....") // a transaction hash

    // check for error
    if err != nil {
        fmt.Printf("error: %v\n", err)
        return
    }

    fmt.Printf("first seen at: %v\n", trace.FirstSeenAt.String())
}
```

<!-- ### Send bundles

Send bundles to Merkle's high performance low latency builder.

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at mbs.merkle.io

    builder := merkleSdk.Builder()

    err := builder.SendBundle(&merkle.Bundle{
        Transactions: []merkle.BundleTx{
            merkle.Tx(tx).CanRevert(),
            merkle.RawTx("0x.....")
        },
        TargetBlock: 300000,
    })

    // check for error
}
``` -->
