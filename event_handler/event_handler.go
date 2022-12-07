package event_handler

import (
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type SmartContractEventHandler struct {
	Address                   common.Address
	JsonAbi                   string
	GoAbi                     *abi.ABI
	MapEventSigHexToEventName map[string]string
}

func NewSmartContractEventHandler(address, jsonAbi string) SmartContractEventHandler {
	contractABI, err := abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		log.Fatal(err)
	}

	return SmartContractEventHandler{
		Address:                   common.HexToAddress("0xE1aBa35771C24837F660430B0bf54c847bA18049"),
		JsonAbi:                   jsonAbi,
		GoAbi:                     &contractABI,
		MapEventSigHexToEventName: CreateMapEventSigHexToName(&contractABI),
	}
}

func (self SmartContractEventHandler) HandleLog(vLog types.Log) error {
	fmt.Println("----------------------------------")
	// event signature hash is always vLog.Topics[0]
	eventSigHex := vLog.Topics[0].Hex()
	eventName, ok := self.MapEventSigHexToEventName[eventSigHex]

	if !ok {
		fmt.Println("event " + eventName + " not defined")
	}

	fmt.Println("event name: " + eventName)
	event, err := self.GoAbi.Unpack(eventName, vLog.Data)
	if err != nil {
		fmt.Printf("err when unpack event %s: %v\n", eventName, err)
	}
	fmt.Println("data:")
	fmt.Println(event)

	return nil
}

func GetEventSigHex(eventSig string) string {
	sigByte := []byte(eventSig)
	return crypto.Keccak256Hash(sigByte).Hex()
}

func CreateMapEventSigHexToName(abiInstance *abi.ABI) map[string]string {
	mapEventSigHexToName := make(map[string]string)

	for k, v := range abiInstance.Events {
		sigHex := GetEventSigHex(v.Sig)
		mapEventSigHexToName[sigHex] = k
	}

	return mapEventSigHexToName
}
