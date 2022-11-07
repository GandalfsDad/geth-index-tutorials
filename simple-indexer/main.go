package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	api_key := os.Getenv("API_KEY")
	client, err := ethclient.Dial("wss://eth-mainnet.g.alchemy.com/v2/" + api_key)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xb9a179DcA5a7bf5f8B9E088437B3A85ebB495eFe")

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(13419942),
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{common.HexToHash("0xda18d31fbb73ed04b84307ef1bc6602e02c855af9f65b53ed10ba43e8d35b7dd")}},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	reader, err := os.Open("./contracts/pt.abi")
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(reader)
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {
		// fmt.Println(vLog.BlockHash.Hex())
		// fmt.Println(vLog.BlockNumber)
		// fmt.Println(vLog.TxHash.Hex())

		val, err := contractAbi.Unpack("ClaimedDraw", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		topics := vLog.Topics[2].Big()

		fmt.Println("Winner:", common.HexToAddress(vLog.Topics[1].Hex()), "Draw:", topics, " Prize:", val[0])
	}

	fmt.Println(len(logs))

}
