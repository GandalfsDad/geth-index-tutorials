package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	rpc := "tbd"

	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal("Failed to connect to the websocket of the RPC ", err)
	} else {
		fmt.Println("Successfully connected to the RPC Endpoint")
	}
}
