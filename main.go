package main

import (
	"context"
	"log"

	"github.com/btcs-longnp/b/contracts/btc_zombie"
	"github.com/btcs-longnp/b/event_handler"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("wss://bsc-testnet.nodereal.io/ws/v1/e9a36765eb8a40b9bd12e680a1fd2bc5")
	if err != nil {
		log.Fatal(err)
	}

	btcZombieJsonAbi, err := btc_zombie.ReadAbiJson()
	if err != nil {
		log.Fatal(err)
	}
	btcZombieEventHandler := event_handler.NewSmartContractEventHandler("0xE1aBa35771C24837F660430B0bf54c847bA18049", btcZombieJsonAbi)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{btcZombieEventHandler.Address},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening event from smart contract address " + btcZombieEventHandler.Address.String())

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			btcZombieEventHandler.HandleLog(vLog)
		}
	}
}
