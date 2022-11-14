package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("Begin Extraction")

	//eth_address := "wss://eth-mainnet.g.alchemy.com/v2/" + os.Getenv("ETH_API_KEY")
	op_address := "https://opt-mainnet.g.alchemy.com/v2/" + os.Getenv("OP_API_KEY")
	poly_address := " wss://polygon-mainnet.g.alchemy.com/v2/" + os.Getenv("POLY_API_KEY")

	//runExtraction(eth_address, "0xdd4d117723C257CEe402285D3aCF218E9A8236E1", "eth_transactions.csv", 13420428, 100000)
	runExtraction(op_address, "0x62BB4fc73094c83B5e952C2180B23fA7054954c4", "op_transactions.csv", 14043021, 1000000)
	runExtraction(poly_address, "0x6a304dFdb9f808741244b6bfEe65ca7B3b3A6076", "poly_transactions.csv", 20226774, 100000)

}

func runExtraction(client_address string, address string, filePath string, minBlock int64, blockGap int64) {
	fmt.Println("Running ", filePath)
	client, err := ethclient.Dial(client_address)
	if err != nil {
		log.Fatal(err)
	}

	tokenAddress := common.HexToAddress(address)
	transferTopic := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

	csvFile, err := os.Create(filePath)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	_ = csvwriter.Write([]string{"block_time", "from", "to", "amt"})

	fromBlock := minBlock
	maxBlockNum, err := client.BlockNumber(context.Background())
	maxBlock := int64(maxBlockNum)

	for fromBlock <= maxBlock {

		fmt.Println("Querying from block ", fromBlock, " to ", fromBlock+blockGap)
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(fromBlock),
			ToBlock:   big.NewInt(fromBlock + blockGap),
			Addresses: []common.Address{tokenAddress},
			Topics:    [][]common.Hash{{transferTopic}},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Fatal(err)
		}

		for _, vLog := range logs {

			val := new(big.Int).SetBytes(vLog.Data)

			lblock, err := client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(vLog.BlockNumber))

			if err != nil {
				log.Fatal("Couldn;t find block")
			}

			ts := time.Unix(int64(lblock.Time), 0)
			//topics := vLog.Topics[2].Big()

			_ = csvwriter.Write([]string{
				ts.UTC().String(),
				common.HexToAddress(vLog.Topics[1].Hex()).String(),
				common.HexToAddress(vLog.Topics[2].Hex()).String(),
				val.String()})

		}

		fromBlock += blockGap
		csvwriter.Flush()

	}

	csvFile.Close()

	client.Close()

}
