![Logo](public/logo.png)

**Merkle is building crypto infrastructure**. [Join us on discord](https://discord.gg/Q9Dc7jVX6c).

# Merkle SDK

The Merkle SDK is a great way to access our products.

## Install

Install the Merkle broker client package:

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

    merkle.SetApiKey("sk_mbs_......") // get one at mbs.usemerkle.com
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

    merkle.SetApiKey("sk_mbs_......") // get one at mbs.usemerkle.com

    txs, err := merkle.Transactions().Stream(1) // pass a chain id

    for {
        select {
            case <-err:
            // error happened
            case tx <- txs:
            // process the transaction
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

    merkle.SetApiKey("sk_mbs_......") // get one at mbs.usemerkle.com

    auctions, err := merkle.Pool().Auctions()

    for {
        select {
            case <-err:
            // error happened
            case auction <- auctions:
            // process the auction, create a backrun

            // then send the bid
            auction.SendBid([]string{
                // hex encoded signed transactions
                "0x...."
            })
        }
    }
}
```