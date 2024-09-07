package main

import (
	"context"
	"log"
	"math/big"
	"strings"
	"time"

	txp "github.com/brankomiric/transactions-parser/parser"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	rpcNode := "https://mainnet.infura.io/v3/c9657d3c5621495c9f6b60c3913df958"

	// Initialize go-eth client
	client, err := ethclient.Dial(rpcNode)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// Get current block number using go-eth
	bn, err := client.BlockNumber(ctx)
	if err != nil {
		log.Println(err)
	}

	// Create transactions-parser instance
	parser := txp.New()

	// Get current block number using tx parser
	parserBn := parser.GetCurrentBlock()

	// Check block numbers equal
	blockNumberAssert := AssertEqual(int(bn), parserBn)
	log.Printf("Block numbers eq: %t", blockNumberAssert)

	bnBig := big.NewInt(int64(bn))

	block, err := client.BlockByNumber(ctx, bnBig)
	if err != nil {
		log.Println(err)
	}

	// Get random addresses from block txs
	addr1 := strings.ToLower(block.Transactions()[0].To().Hex())
	addr2 := strings.ToLower(block.Transactions()[10].To().Hex())

	// Subscribe to parser
	parser.Subscribe(addr1)
	parser.Subscribe(addr2)

	// Poll for txs
	for {
		time.Sleep(10 * time.Second)
		txs1 := parser.GetTransactions(addr1)
		txs2 := parser.GetTransactions(addr2)
		log.Printf("Addr %s tx count %d\n", addr1, len(txs1))
		log.Printf("Addr %s tx count %d\n", addr2, len(txs2))

	}
}

func AssertEqual[T comparable](a, b T) bool {
	return a == b
}
