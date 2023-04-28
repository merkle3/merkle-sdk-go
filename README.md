![Logo](public/logo.png)

**Merkle is building crypto infrastructure**. [Join us on discord](https://discord.gg/Q9Dc7jVX6c).

# Merkle SDK

The Merkle SDK is a great way to access our products.

## Install

Install the Merkle SDK package:

```
go get github.com/merkle3/merkle-sdk-go
```

## Authentication

Get an API key from [Merkle Blockchain Services (MBS)](https://mbs.usemerkle.com). It is free.

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at mbs.usemerkle.com
}
```

## Features

### Stream transactions

Access Merkle's private stream of transactions on Ethereum.

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at mbs.usemerkle.com

    txs, err := merkleSdk.Transactions().Stream(merkle.EthereumMainnet) // pass a chain id

    for {
        select {
            case <-err:
			    // error happened
			    fmt.Printf("error: %v\n", err)
            case tx := <-txs:
                // process the transaction
                fmt.Printf("hash: %v\n", tx.Hash().String())
        }
    }
}
```

### Stream auctions

Stream auctions from the Merkle Private Pool. [Learn more in the docs](https://docs.usemerkle.com/private-pool/what-is-merkle-private-pool).

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at mbs.usemerkle.com

    auctions, err := merkleSdk.Pool().Auctions()

    for {
        select {
            case <-err:
            // error happened
            case auction <- auctions:
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

<!-- ### Send bundles

Send bundles to Merkle's high performance low latency builder.

```golang
package main

import (
    "github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
    merkleSdk := merkle.New()

    merkleSdk.SetApiKey("sk_mbs_......") // get one at mbs.usemerkle.com

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