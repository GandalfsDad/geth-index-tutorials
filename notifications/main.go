package main

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gen2brain/beeep"
)

func main() {
	api_key := os.Getenv("API_KEY")
	client, err := ethclient.Dial("wss://eth-mainnet.g.alchemy.com/v2/" + api_key)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x5A98FcBEA516Cf06857215779Fd812CA3beF1B32")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")}},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:

			id := new(big.Float).SetInt(new(big.Int).SetBytes(vLog.Data))
			dec := new(big.Float).SetFloat64(1e-18)
			amt := id.Mul(id, dec)

			msg := "From: " + common.HexToAddress(vLog.Topics[1].Hex()).String() + " To: " + common.HexToAddress(vLog.Topics[2].Hex()).String() + " Amt: " + amt.String()

			err := beeep.Notify("Transfer", msg, "")
			if err != nil {
				panic(err)
			}
		}
	}

}
