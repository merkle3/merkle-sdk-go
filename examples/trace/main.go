package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/merkle3/merkle-sdk-go/merkle"
)

func main() {
	godotenv.Load()

	merkleSdk := merkle.New()

	merkleSdk.SetApiKey(os.Getenv("MERKLE_API_KEY"))

	// trace 0xfdccc024d726b2c3e7131cb75949d4ac8616c36c1ef2bfb18fcd14cb0d0f1a61 on Polygon
	trace, _ := merkleSdk.Transactions().Trace("0xfdccc024d726b2c3e7131cb75949d4ac8616c36c1ef2bfb18fcd14cb0d0f1a61")

	fmt.Printf("trace: %+v\n", trace)
}
